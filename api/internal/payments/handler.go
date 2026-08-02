package payments

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const maximumRequestBytes = 32 << 10

type ErrorCode string

const (
	ErrorCodeInvalidBody   ErrorCode = "invalid_body"
	ErrorCodeOrderNotFound ErrorCode = "order_not_found"
	ErrorCodeInvalidStatus ErrorCode = "invalid_status"
	ErrorCodeInternal      ErrorCode = "internal_error"
)

type PaymentService interface {
	Summary(context.Context, string) (Summary, error)
	Initialize(context.Context, string) (InitializeResult, error)
	Verify(context.Context, string, string) (VerifyResult, error)
}

type Handler struct {
	service PaymentService
}

func NewHandler(service PaymentService) *Handler {
	return &Handler{service: service}
}

type initializeResponse struct {
	Payment InitializeResult `json:"payment"`
}

type verifyRequest struct {
	Reference string `json:"reference"`
}

type verifyResponse struct {
	Result VerifyResult `json:"result"`
}

func (handler *Handler) Summary(responseWriter http.ResponseWriter, request *http.Request) {
	summary, err := handler.service.Summary(request.Context(), request.PathValue("payToken"))
	if err != nil {
		handler.writePaymentError(responseWriter, err)
		return
	}

	writeJSON(responseWriter, http.StatusOK, summary)
}

func (handler *Handler) Initialize(responseWriter http.ResponseWriter, request *http.Request) {
	payment, err := handler.service.Initialize(request.Context(), request.PathValue("payToken"))
	if err != nil {
		handler.writePaymentError(responseWriter, err)
		return
	}

	writeJSON(responseWriter, http.StatusOK, initializeResponse{Payment: payment})
}

func (handler *Handler) Verify(responseWriter http.ResponseWriter, request *http.Request) {
	var input verifyRequest

	decoder := json.NewDecoder(http.MaxBytesReader(responseWriter, request.Body, maximumRequestBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil || input.Reference == "" {
		writeError(responseWriter, http.StatusBadRequest, ErrorCodeInvalidBody,
			"the payment reference could not be read")
		return
	}

	result, err := handler.service.Verify(
		request.Context(), request.PathValue("payToken"), input.Reference,
	)
	if err != nil {
		handler.writePaymentError(responseWriter, err)
		return
	}

	writeJSON(responseWriter, http.StatusOK, verifyResponse{Result: result})
}

func (handler *Handler) writePaymentError(responseWriter http.ResponseWriter, paymentError error) {
	switch {
	case errors.Is(paymentError, ErrOrderNotFound), errors.Is(paymentError, ErrReferenceNotFound):
		writeError(responseWriter, http.StatusNotFound, ErrorCodeOrderNotFound,
			"the payment link could not be found")
	case errors.Is(paymentError, ErrOrderNotPayable), errors.Is(paymentError, ErrPaymentAlreadyPaid):
		writeError(responseWriter, http.StatusConflict, ErrorCodeInvalidStatus,
			"this order cannot be paid right now")
	case errors.Is(paymentError, ErrPaymentInProgress):
		writeError(responseWriter, http.StatusConflict, ErrorCodeInvalidStatus,
			"a payment is already being prepared for this order")
	case errors.Is(paymentError, ErrPaymentUnavailable):
		writeError(responseWriter, http.StatusServiceUnavailable, ErrorCodeInternal,
			"payments are not available right now")
	case errors.Is(paymentError, ErrReferenceMismatch),
		errors.Is(paymentError, ErrPaymentAmountMismatch),
		errors.Is(paymentError, ErrPaymentCurrency):
		writeError(responseWriter, http.StatusConflict, ErrorCodeInvalidStatus,
			"the payment could not be verified")
	default:
		log.Printf("payment request failed: %v", paymentError)
		writeError(responseWriter, http.StatusInternalServerError, ErrorCodeInternal,
			"the payment request could not be completed")
	}
}

type errorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

func writeError(
	responseWriter http.ResponseWriter,
	statusCode int,
	errorCode ErrorCode,
	errorMessage string,
) {
	writeJSON(responseWriter, statusCode, errorResponse{
		Error: errorDetail{Code: errorCode, Message: errorMessage},
	})
}

func writeJSON(responseWriter http.ResponseWriter, statusCode int, payload any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)

	if err := json.NewEncoder(responseWriter).Encode(payload); err != nil {
		log.Printf("encode payment response: %v", err)
	}
}
