package products

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

// ErrorCode identifies a stable machine-readable product API error.
type ErrorCode string

const (
	ErrorCodeInvalidProductSection ErrorCode = "invalid_product_section"
	ErrorCodeInvalidProductSlug    ErrorCode = "invalid_product_slug"
	ErrorCodeInvalidActiveFilter   ErrorCode = "invalid_active_filter"
	ErrorCodeInvalidProductDraft   ErrorCode = "invalid_product_draft"
	ErrorCodeDuplicateProductSlug  ErrorCode = "duplicate_product_slug"
	ErrorCodeProductNotFound       ErrorCode = "product_not_found"
	ErrorCodeSearchQueryTooLong    ErrorCode = "search_query_too_long"
	ErrorCodeInvalidBody           ErrorCode = "invalid_body"
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

// AdminProductManager defines the owner-facing product operations exposed over HTTP.
type AdminProductManager interface {
	List(context.Context, AdminListInput) (AdminProductPage, error)
	GetBySlug(context.Context, string) (AdminProduct, error)
	Create(context.Context, AdminProductDraft) (AdminProduct, error)
	Update(context.Context, string, AdminProductDraft) (AdminProduct, error)
	UpdateActive(context.Context, string, AdminActiveUpdate) error
}

// AdminHandler translates authenticated admin product requests.
type AdminHandler struct {
	productManager AdminProductManager
}

type listResponse struct {
	Products []Product `json:"products"`
	Count    int       `json:"count"`
}

type detailResponse struct {
	Product Product `json:"product"`
}

type adminDetailResponse struct {
	Product AdminProduct `json:"product"`
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

// NewAdminHandler creates an admin product HTTP handler.
func NewAdminHandler(productManager AdminProductManager) *AdminHandler {
	return &AdminHandler{productManager: productManager}
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

// List returns owner-facing products, including inactive products.
func (handler *AdminHandler) List(responseWriter http.ResponseWriter, request *http.Request) {
	queryParameters := request.URL.Query()

	page, err := handler.productManager.List(request.Context(), AdminListInput{
		Section: queryParameters.Get("section"),
		Query:   queryParameters.Get("q"),
		Active:  queryParameters.Get("active"),
		Limit:   positiveInt(queryParameters.Get("limit")),
		Offset:  positiveInt(queryParameters.Get("offset")),
	})
	if err != nil {
		handler.writeServiceError(responseWriter, err)
		return
	}
	if page.Products == nil {
		page.Products = make([]AdminProduct, 0)
	}

	writeJSON(responseWriter, http.StatusOK, page)
}

// GetBySlug returns one owner-facing product for editing.
func (handler *AdminHandler) GetBySlug(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	product, err := handler.productManager.GetBySlug(
		request.Context(),
		request.PathValue("slug"),
	)
	if err != nil {
		handler.writeServiceError(responseWriter, err)
		return
	}

	writeJSON(responseWriter, http.StatusOK, adminDetailResponse{Product: product})
}

// Create inserts one owner-managed product.
func (handler *AdminHandler) Create(responseWriter http.ResponseWriter, request *http.Request) {
	draft, ok := decodeProductDraft(responseWriter, request)
	if !ok {
		return
	}

	product, err := handler.productManager.Create(request.Context(), draft)
	if err != nil {
		handler.writeServiceError(responseWriter, err)
		return
	}

	writeJSON(responseWriter, http.StatusCreated, adminDetailResponse{Product: product})
}

// Update changes owner-editable product metadata.
func (handler *AdminHandler) Update(responseWriter http.ResponseWriter, request *http.Request) {
	draft, ok := decodeProductDraft(responseWriter, request)
	if !ok {
		return
	}

	product, err := handler.productManager.Update(
		request.Context(),
		request.PathValue("slug"),
		draft,
	)
	if err != nil {
		handler.writeServiceError(responseWriter, err)
		return
	}

	writeJSON(responseWriter, http.StatusOK, adminDetailResponse{Product: product})
}

// UpdateActive archives or restores one product from the public catalogue.
func (handler *AdminHandler) UpdateActive(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	var update AdminActiveUpdate

	decoder := json.NewDecoder(http.MaxBytesReader(responseWriter, request.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&update); err != nil {
		writeError(responseWriter, http.StatusBadRequest, ErrorCodeInvalidBody,
			"the product update could not be read")
		return
	}

	err := handler.productManager.UpdateActive(
		request.Context(),
		request.PathValue("slug"),
		update,
	)
	if err != nil {
		handler.writeServiceError(responseWriter, err)
		return
	}

	responseWriter.WriteHeader(http.StatusNoContent)
}

// decodeProductDraft reads a strict owner-editable product payload.
func decodeProductDraft(
	responseWriter http.ResponseWriter,
	request *http.Request,
) (AdminProductDraft, bool) {
	var draft AdminProductDraft

	decoder := json.NewDecoder(http.MaxBytesReader(responseWriter, request.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&draft); err != nil {
		writeError(responseWriter, http.StatusBadRequest, ErrorCodeInvalidBody,
			"the product could not be read")
		return AdminProductDraft{}, false
	}

	return draft, true
}

// writeServiceError maps domain errors to stable public status codes and messages.
func (handler *Handler) writeServiceError(
	responseWriter http.ResponseWriter,
	serviceError error,
) {
	writeServiceError(responseWriter, serviceError)
}

// writeServiceError maps domain errors to stable admin status codes and messages.
func (handler *AdminHandler) writeServiceError(
	responseWriter http.ResponseWriter,
	serviceError error,
) {
	writeServiceError(responseWriter, serviceError)
}

// writeServiceError maps domain errors to stable HTTP status codes and messages.
func writeServiceError(
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
	case errors.Is(serviceError, ErrInvalidActiveFilter):
		writeError(
			responseWriter,
			http.StatusBadRequest,
			ErrorCodeInvalidActiveFilter,
			"active must be true or false",
		)
	case errors.Is(serviceError, ErrInvalidProductDraft):
		writeError(
			responseWriter,
			http.StatusBadRequest,
			ErrorCodeInvalidProductDraft,
			"product fields are incomplete or invalid",
		)
	case errors.Is(serviceError, ErrDuplicateProductSlug):
		writeError(
			responseWriter,
			http.StatusConflict,
			ErrorCodeDuplicateProductSlug,
			"product slug already exists",
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

// positiveInt reads a paging parameter and leaves invalid values for service defaults.
func positiveInt(value string) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return 0
	}

	return parsed
}
