package products

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// ErrorCode identifies a stable machine-readable product API error.
type ErrorCode string

const (
	ErrorCodeInvalidProductSection ErrorCode = "invalid_product_section"
	ErrorCodeInvalidProductSlug    ErrorCode = "invalid_product_slug"
	ErrorCodeProductNotFound       ErrorCode = "product_not_found"
	ErrorCodeSearchQueryTooLong    ErrorCode = "search_query_too_long"
	ErrorCodeInternal              ErrorCode = "internal_error"
)

// ProductReader defines the product service operations exposed over HTTP.
type ProductReader interface {
	List(context.Context, ListInput) ([]Product, error)
	GetBySlug(context.Context, string) (Product, error)
}

// Handler translates HTTP requests into product service calls.
type Handler struct {
	productReader ProductReader
}

type listResponse struct {
	Products []Product `json:"products"`
	Count    int       `json:"count"`
}

type detailResponse struct {
	Product Product `json:"product"`
}

type errorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

// NewHandler creates a product HTTP handler around the service contract.
func NewHandler(productReader ProductReader) *Handler {
	return &Handler{productReader: productReader}
}

// List handles product collection requests with optional section and search filters.
func (handler *Handler) List(responseWriter http.ResponseWriter, request *http.Request) {
	// Translate query parameters into the untrusted service input model.
	queryParameters := request.URL.Query()
	products, err := handler.productReader.List(request.Context(), ListInput{
		Section: queryParameters.Get("section"),
		Query:   queryParameters.Get("q"),
	})
	if err != nil {
		handler.writeServiceError(responseWriter, err)
		return
	}

	// Preserve an empty JSON array so frontend iteration never receives null.
	if products == nil {
		products = make([]Product, 0)
	}

	writeJSON(responseWriter, http.StatusOK, listResponse{
		Products: products,
		Count:    len(products),
	})
}

// GetBySlug handles one product detail request using the canonical path slug.
func (handler *Handler) GetBySlug(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	product, err := handler.productReader.GetBySlug(
		request.Context(),
		request.PathValue("slug"),
	)
	if err != nil {
		handler.writeServiceError(responseWriter, err)
		return
	}

	writeJSON(responseWriter, http.StatusOK, detailResponse{Product: product})
}

// writeServiceError maps domain errors to stable public status codes and messages.
func (handler *Handler) writeServiceError(
	responseWriter http.ResponseWriter,
	serviceError error,
) {
	// Return specific client-safe failures without exposing persistence details.
	switch {
	case errors.Is(serviceError, ErrInvalidProductSection):
		writeError(
			responseWriter,
			http.StatusBadRequest,
			ErrorCodeInvalidProductSection,
			"section must be suits, trad, casuals-men, or casuals-women",
		)
	case errors.Is(serviceError, ErrInvalidProductSlug):
		writeError(
			responseWriter,
			http.StatusBadRequest,
			ErrorCodeInvalidProductSlug,
			"product slug is invalid",
		)
	case errors.Is(serviceError, ErrSearchQueryTooLong):
		writeError(
			responseWriter,
			http.StatusBadRequest,
			ErrorCodeSearchQueryTooLong,
			"search query must be 100 characters or fewer",
		)
	case errors.Is(serviceError, ErrProductNotFound):
		writeError(
			responseWriter,
			http.StatusNotFound,
			ErrorCodeProductNotFound,
			"product was not found",
		)
	default:
		log.Printf("product request failed: %v", serviceError)
		writeError(
			responseWriter,
			http.StatusInternalServerError,
			ErrorCodeInternal,
			"the product request could not be completed",
		)
	}
}

// writeError writes one consistent JSON error envelope.
func writeError(
	responseWriter http.ResponseWriter,
	statusCode int,
	errorCode ErrorCode,
	errorMessage string,
) {
	writeJSON(responseWriter, statusCode, errorResponse{
		Error: errorDetail{
			Code:    errorCode,
			Message: errorMessage,
		},
	})
}

// writeJSON writes a JSON response with the requested HTTP status.
func writeJSON(responseWriter http.ResponseWriter, statusCode int, payload any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)

	if err := json.NewEncoder(responseWriter).Encode(payload); err != nil {
		log.Printf("encode product response: %v", err)
	}
}
