package paystack

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const maximumWebhookBytes = 128 << 10

type WebhookService interface {
	ConfirmWebhookReference(context.Context, string) error
}

type WebhookHandler struct {
	service   WebhookService
	secretKey string
}

func NewWebhookHandler(service WebhookService, secretKey string) *WebhookHandler {
	return &WebhookHandler{service: service, secretKey: secretKey}
}

type webhookEvent struct {
	Event string `json:"event"`
	Data  struct {
		Reference string `json:"reference"`
	} `json:"data"`
}

func (handler *WebhookHandler) Handle(responseWriter http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(http.MaxBytesReader(responseWriter, request.Body, maximumWebhookBytes))
	if err != nil {
		writeWebhookError(responseWriter, http.StatusBadRequest, "the webhook could not be read")
		return
	}
	if !validSignature(body, request.Header.Get("x-paystack-signature"), handler.secretKey) {
		writeWebhookError(responseWriter, http.StatusUnauthorized, "the webhook signature is invalid")
		return
	}

	var event webhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		writeWebhookError(responseWriter, http.StatusBadRequest, "the webhook could not be parsed")
		return
	}

	if event.Event == "charge.success" && event.Data.Reference != "" {
		if err := handler.service.ConfirmWebhookReference(request.Context(), event.Data.Reference); err != nil {
			log.Printf("confirm paystack webhook failed: %v", err)
			writeWebhookError(responseWriter, http.StatusInternalServerError,
				"the webhook could not be processed")
			return
		}
	}

	responseWriter.WriteHeader(http.StatusOK)
}

func validSignature(body []byte, signature string, secretKey string) bool {
	if signature == "" || secretKey == "" {
		return false
	}

	mac := hmac.New(sha512.New, []byte(secretKey))
	mac.Write(body)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

func writeWebhookError(responseWriter http.ResponseWriter, statusCode int, message string) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)

	if err := json.NewEncoder(responseWriter).Encode(map[string]map[string]string{
		"error": {"message": message},
	}); err != nil {
		log.Printf("encode paystack webhook response: %v", err)
	}
}
