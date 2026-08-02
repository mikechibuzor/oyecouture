package orders

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// ErrorCode identifies a stable machine-readable order API error.
type ErrorCode string

const (
	ErrorCodeInvalidBody     ErrorCode = "invalid_body"
	ErrorCodeInvalidCustomer ErrorCode = "invalid_customer"
	ErrorCodeInvalidItems    ErrorCode = "invalid_items"
	ErrorCodeInvalidDelivery ErrorCode = "invalid_delivery"
	ErrorCodeUnavailable     ErrorCode = "product_unavailable"
	ErrorCodeOrderNotFound   ErrorCode = "order_not_found"
	ErrorCodeInvalidStatus   ErrorCode = "invalid_status"
	ErrorCodeInternal        ErrorCode = "internal_error"
)

// Body larger than this is rejected before parsing.
const maximumRequestBytes = 32 << 10

// OrderCreator is the service contract exposed over HTTP.
type OrderCreator interface {
	Create(context.Context, CreateOrderInput) (Order, error)
}

type Handler struct {
	orderCreator OrderCreator
}

func NewHandler(orderCreator OrderCreator) *Handler {
	return &Handler{orderCreator: orderCreator}
}

type createResponse struct {
	Order Order `json:"order"`
}

type errorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

// Create handles order submission from the storefront.
func (handler *Handler) Create(responseWriter http.ResponseWriter, request *http.Request) {
	var input CreateOrderInput

	decoder := json.NewDecoder(http.MaxBytesReader(responseWriter, request.Body, maximumRequestBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		writeError(responseWriter, http.StatusBadRequest, ErrorCodeInvalidBody,
			"the order could not be read")
		return
	}

	order, err := handler.orderCreator.Create(request.Context(), input)
	if err != nil {
		handler.writeServiceError(responseWriter, err)
		return
	}

	writeJSON(responseWriter, http.StatusCreated, createResponse{Order: order})
}

// writeServiceError maps domain errors to stable public status codes, without
// leaking persistence detail.
func (handler *Handler) writeServiceError(responseWriter http.ResponseWriter, serviceError error) {
	switch {
	case errors.Is(serviceError, ErrInvalidCustomer):
		writeError(responseWriter, http.StatusBadRequest, ErrorCodeInvalidCustomer,
			"please check the name, email and phone number")
	case errors.Is(serviceError, ErrNoItems),
		errors.Is(serviceError, ErrTooManyItems),
		errors.Is(serviceError, ErrInvalidItem),
		errors.Is(serviceError, ErrSizeRequired),
		errors.Is(serviceError, ErrSizeNotApplicable):
		writeError(responseWriter, http.StatusBadRequest, ErrorCodeInvalidItems,
			"the bag contains an item we cannot order")
	case errors.Is(serviceError, ErrDeliveryRequired),
		errors.Is(serviceError, ErrUnknownZone):
		writeError(responseWriter, http.StatusBadRequest, ErrorCodeInvalidDelivery,
			"please add a delivery address and choose a delivery area")
	case errors.Is(serviceError, ErrUnknownProduct):
		writeError(responseWriter, http.StatusConflict, ErrorCodeUnavailable,
			"one of the pieces is no longer available")
	default:
		log.Printf("create order failed: %v", serviceError)
		writeError(responseWriter, http.StatusInternalServerError, ErrorCodeInternal,
			"the order could not be placed")
	}
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
		log.Printf("encode order response: %v", err)
	}
}
