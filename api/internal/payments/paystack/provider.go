package paystack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/uzor-co/oyecouture/api/internal/payments"
)

const defaultBaseURL = "https://api.paystack.co"

type Provider struct {
	secretKey  string
	baseURL    string
	httpClient *http.Client
}

func NewProvider(secretKey string) *Provider {
	return &Provider{
		secretKey:  secretKey,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: 20 * time.Second},
	}
}

func (provider *Provider) Name() string { return payments.ProviderPaystack }

type initializeRequest struct {
	Email       string         `json:"email"`
	Amount      string         `json:"amount"`
	Currency    string         `json:"currency"`
	Reference   string         `json:"reference"`
	CallbackURL string         `json:"callback_url"`
	Metadata    map[string]any `json:"metadata"`
}

type initializeResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

type verifyResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Reference       string     `json:"reference"`
		Status          string     `json:"status"`
		Amount          int64      `json:"amount"`
		Currency        string     `json:"currency"`
		PaidAt          *time.Time `json:"paid_at"`
		GatewayResponse string     `json:"gateway_response"`
	} `json:"data"`
}

func (provider *Provider) Initialize(
	requestContext context.Context,
	request payments.InitializeRequest,
) (payments.InitializeResult, []byte, error) {
	if provider.secretKey == "" {
		return payments.InitializeResult{}, nil, payments.ErrPaymentUnavailable
	}

	payload, err := json.Marshal(initializeRequest{
		Email:       request.Email,
		Amount:      strconv.FormatInt(request.AmountKobo, 10),
		Currency:    request.Currency,
		Reference:   request.Reference,
		CallbackURL: request.CallbackURL,
		Metadata: map[string]any{
			"order_no": request.OrderNo,
			"customer": request.Customer,
		},
	})
	if err != nil {
		return payments.InitializeResult{}, nil, fmt.Errorf("encode paystack initialize: %w", err)
	}

	responseBody, err := provider.do(requestContext, http.MethodPost, "/transaction/initialize", payload)
	if err != nil {
		return payments.InitializeResult{}, responseBody, err
	}

	var response initializeResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return payments.InitializeResult{}, responseBody, fmt.Errorf("decode paystack initialize: %w", err)
	}
	if !response.Status || response.Data.AuthorizationURL == "" || response.Data.Reference == "" {
		return payments.InitializeResult{}, responseBody, fmt.Errorf("paystack initialize rejected: %s", response.Message)
	}

	return payments.InitializeResult{
		AuthorizationURL: response.Data.AuthorizationURL,
		AccessCode:       response.Data.AccessCode,
		Reference:        response.Data.Reference,
	}, responseBody, nil
}

func (provider *Provider) Verify(
	requestContext context.Context,
	reference string,
) (payments.Verification, error) {
	if provider.secretKey == "" {
		return payments.Verification{}, payments.ErrPaymentUnavailable
	}

	responseBody, err := provider.do(
		requestContext, http.MethodGet, "/transaction/verify/"+url.PathEscape(reference), nil,
	)
	if err != nil {
		return payments.Verification{Reference: reference, ProviderPayload: responseBody}, err
	}

	var response verifyResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return payments.Verification{Reference: reference, ProviderPayload: responseBody},
			fmt.Errorf("decode paystack verify: %w", err)
	}
	if !response.Status {
		return payments.Verification{Reference: reference, ProviderPayload: responseBody},
			fmt.Errorf("paystack verify rejected: %s", response.Message)
	}

	return payments.Verification{
		Reference:       response.Data.Reference,
		Status:          response.Data.Status,
		AmountKobo:      response.Data.Amount,
		Currency:        response.Data.Currency,
		PaidAt:          response.Data.PaidAt,
		GatewayResponse: response.Data.GatewayResponse,
		ProviderPayload: responseBody,
	}, nil
}

func (provider *Provider) do(
	requestContext context.Context,
	method string,
	path string,
	payload []byte,
) ([]byte, error) {
	var body io.Reader = http.NoBody
	if payload != nil {
		body = bytes.NewReader(payload)
	}

	request, err := http.NewRequestWithContext(
		requestContext, method, provider.baseURL+path, body,
	)
	if err != nil {
		return nil, fmt.Errorf("create paystack request: %w", err)
	}

	request.Header.Set("Authorization", "Bearer "+provider.secretKey)
	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := provider.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("call paystack: %w", err)
	}
	defer response.Body.Close()

	responseBody := new(bytes.Buffer)
	if _, err := responseBody.ReadFrom(response.Body); err != nil {
		return nil, fmt.Errorf("read paystack response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return responseBody.Bytes(), fmt.Errorf("paystack returned %s", response.Status)
	}

	return responseBody.Bytes(), nil
}
