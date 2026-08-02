package sendbyte

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/uzor-co/oyecouture/api/internal/email"
)

const defaultBaseURL = "https://api.sendbyte.africa"

type SendByte struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewSendByte configures the SendByte REST email provider.
func NewSendByte(apiKey string) *SendByte {
	return &SendByte{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// Name returns the provider identifier used in email audit rows.
func (sendByte *SendByte) Name() string { return email.ProviderSendByte }

type sendRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	HTML    string `json:"html"`
	Text    string `json:"text,omitempty"`
}

// Send posts a transactional email to SendByte and returns its response payload.
func (sendByte *SendByte) Send(
	requestContext context.Context,
	message email.Message,
) (email.SendResult, []byte, error) {
	if sendByte.apiKey == "" {
		return email.SendResult{}, nil, fmt.Errorf("sendbyte api key is not configured")
	}

	payload, err := json.Marshal(sendRequest{
		From:    message.From,
		To:      message.To,
		Subject: message.Subject,
		HTML:    message.HTML,
		Text:    message.Text,
	})
	if err != nil {
		return email.SendResult{}, nil, fmt.Errorf("encode sendbyte email: %w", err)
	}

	responseBody, err := sendByte.do(requestContext, payload, message.IdempotencyKey)
	if err != nil {
		return email.SendResult{}, responseBody, err
	}

	return email.SendResult{MessageID: messageID(responseBody)}, responseBody, nil
}

// do execute the SendByte HTTP request with an idempotency key when available.
func (sendByte *SendByte) do(
	requestContext context.Context,
	payload []byte,
	idempotencyKey string,
) ([]byte, error) {
	request, err := http.NewRequestWithContext(
		requestContext, http.MethodPost, sendByte.baseURL+"/v1/emails", bytes.NewReader(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("create sendbyte request: %w", err)
	}

	request.Header.Set("Authorization", "Bearer "+sendByte.apiKey)
	request.Header.Set("Content-Type", "application/json")
	if idempotencyKey != "" {
		request.Header.Set("Idempotency-Key", idempotencyKey)
	}

	response, err := sendByte.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("call sendbyte: %w", err)
	}
	defer response.Body.Close()

	responseBody := new(bytes.Buffer)
	if _, err := io.Copy(responseBody, response.Body); err != nil {
		return nil, fmt.Errorf("read sendbyte response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return responseBody.Bytes(), fmt.Errorf("sendbyte returned %s", response.Status)
	}

	return responseBody.Bytes(), nil
}

// messageID extracts common message id shapes from SendByte responses.
func messageID(responseBody []byte) string {
	var response map[string]any
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return ""
	}

	for _, key := range []string{"id", "message_id", "messageId"} {
		if value, ok := response[key].(string); ok {
			return value
		}
	}

	data, ok := response["data"].(map[string]any)
	if !ok {
		return ""
	}
	for _, key := range []string{"id", "message_id", "messageId"} {
		if value, ok := data[key].(string); ok {
			return value
		}
	}

	return ""
}
