package paystack

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/uzor-co/oyecouture/api/internal/payments"
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
			// Paystack retries any non-2xx. Redelivering an event we can never
			// accept would retry forever, so those are acknowledged and logged
			// instead. An amount or currency mismatch is a signal for a human,
			// not something a retry can resolve.
			if isTerminal(err) {
				log.Printf("paystack webhook not retryable, acknowledged: reference=%s: %v",
					event.Data.Reference, err)
				responseWriter.WriteHeader(http.StatusOK)
				return
			}

			log.Printf("confirm paystack webhook failed, will be retried: reference=%s: %v",
				event.Data.Reference, err)
			writeWebhookError(responseWriter, http.StatusInternalServerError,
				"the webhook could not be processed")
			return
		}
	}

	responseWriter.WriteHeader(http.StatusOK)
}

// terminalErrors cannot be resolved by delivering the event again: the
// reference is unknown to us, or the payment does not match the order it
// claims to settle.
var terminalErrors = []error{
	payments.ErrReferenceNotFound,
	payments.ErrReferenceMismatch,
	payments.ErrPaymentAmountMismatch,
	payments.ErrPaymentCurrency,
	payments.ErrOrderNotPayable,
	payments.ErrOrderNotFound,
}

// isTerminal reports whether retrying is pointless. Anything unrecognised is
// deliberately treated as transient: a database blip or an unreachable Paystack
// must be retried, and wrongly acknowledging one would drop a real payment.
func isTerminal(err error) bool {
	for _, terminal := range terminalErrors {
		if errors.Is(err, terminal) {
			return true
		}
	}

	// Paystack answering "this transaction is not valid" is authoritative;
	// asking again returns the same answer.
	if errors.Is(err, ErrVerifyRejected) {
		return true
	}

	// A 4xx from Paystack means the request itself is wrong, so redelivering it
	// changes nothing. A 5xx is their side struggling and is worth retrying.
	var apiError *APIError
	if errors.As(err, &apiError) {
		return !apiError.Retryable()
	}

	return false
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
