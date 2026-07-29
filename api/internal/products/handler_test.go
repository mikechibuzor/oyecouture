package products

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// productReaderStub records handler inputs and returns controlled service results.
type productReaderStub struct {
	listInput     ListInput
	listProducts  []Product
	listError     error
	detailSlug    string
	detailProduct Product
	detailError   error
}

// List records collection input and returns the configured product result.
func (reader *productReaderStub) List(
	_ context.Context,
	input ListInput,
) ([]Product, error) {
	reader.listInput = input
	return reader.listProducts, reader.listError
}

// GetBySlug records the detail slug and returns the configured product result.
func (reader *productReaderStub) GetBySlug(
	_ context.Context,
	productSlug string,
) (Product, error) {
	reader.detailSlug = productSlug
	return reader.detailProduct, reader.detailError
}

// TestHandlerListReturnsTypedEnvelopeAndFilters verifies collection request translation.
func TestHandlerListReturnsTypedEnvelopeAndFilters(t *testing.T) {
	t.Parallel()

	reader := &productReaderStub{
		listProducts: []Product{{Slug: "linen", PriceKobo: 7_800_000}},
	}
	handler := NewHandler(reader)
	request := httptest.NewRequest(
		http.MethodGet,
		"/products?section=casuals-men&q=linen",
		nil,
	)
	responseRecorder := httptest.NewRecorder()

	handler.List(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}
	if responseRecorder.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("expected JSON content type, got %q", responseRecorder.Header().Get("Content-Type"))
	}
	if reader.listInput.Section != "casuals-men" || reader.listInput.Query != "linen" {
		t.Fatalf("unexpected list input: %+v", reader.listInput)
	}

	var response listResponse
	if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if response.Count != 1 || len(response.Products) != 1 {
		t.Fatalf("unexpected list response: %+v", response)
	}
	if response.Products[0].PriceKobo != 7_800_000 {
		t.Fatalf("expected integer-kobo price, got %d", response.Products[0].PriceKobo)
	}
}

// TestHandlerListReturnsEmptyArrayForNoProducts verifies stable empty-list JSON.
func TestHandlerListReturnsEmptyArrayForNoProducts(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&productReaderStub{})
	request := httptest.NewRequest(http.MethodGet, "/products?q=missing", nil)
	responseRecorder := httptest.NewRecorder()

	handler.List(responseRecorder, request)

	if responseRecorder.Body.String() != "{\"products\":[],\"count\":0}\n" {
		t.Fatalf("unexpected empty response %q", responseRecorder.Body.String())
	}
}

// TestHandlerGetBySlugReturnsProduct verifies path-value translation and detail response.
func TestHandlerGetBySlugReturnsProduct(t *testing.T) {
	t.Parallel()

	reader := &productReaderStub{
		detailProduct: Product{Slug: "agbada", HasSizes: false},
	}
	handler := NewHandler(reader)
	request := httptest.NewRequest(http.MethodGet, "/products/agbada", nil)
	request.SetPathValue("slug", "agbada")
	responseRecorder := httptest.NewRecorder()

	handler.GetBySlug(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}
	if reader.detailSlug != "agbada" {
		t.Fatalf("expected agbada slug, got %q", reader.detailSlug)
	}

	var response detailResponse
	if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode detail response: %v", err)
	}
	if response.Product.Slug != "agbada" || response.Product.HasSizes {
		t.Fatalf("unexpected detail response: %+v", response.Product)
	}
}

// TestHandlerMapsProductNotFoundToJSON verifies the public not-found contract.
func TestHandlerMapsProductNotFoundToJSON(t *testing.T) {
	t.Parallel()

	reader := &productReaderStub{detailError: ErrProductNotFound}
	handler := NewHandler(reader)
	request := httptest.NewRequest(http.MethodGet, "/products/missing", nil)
	request.SetPathValue("slug", "missing")
	responseRecorder := httptest.NewRecorder()

	handler.GetBySlug(responseRecorder, request)

	if responseRecorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, responseRecorder.Code)
	}
	var response errorResponse
	if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if response.Error.Code != ErrorCodeProductNotFound {
		t.Fatalf("expected not-found code, got %q", response.Error.Code)
	}
}

// TestHandlerMapsInvalidSectionToBadRequest verifies client validation errors.
func TestHandlerMapsInvalidSectionToBadRequest(t *testing.T) {
	t.Parallel()

	reader := &productReaderStub{listError: ErrInvalidProductSection}
	handler := NewHandler(reader)
	request := httptest.NewRequest(http.MethodGet, "/products?section=invalid", nil)
	responseRecorder := httptest.NewRecorder()

	handler.List(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, responseRecorder.Code)
	}
	var response errorResponse
	if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if response.Error.Code != ErrorCodeInvalidProductSection {
		t.Fatalf("expected invalid-section code, got %q", response.Error.Code)
	}
}

// TestHandlerMapsUnexpectedErrorToInternalFailure verifies persistence errors stay private.
func TestHandlerMapsUnexpectedErrorToInternalFailure(t *testing.T) {
	t.Parallel()

	reader := &productReaderStub{listError: errors.New("database password leaked")}
	handler := NewHandler(reader)
	request := httptest.NewRequest(http.MethodGet, "/products", nil)
	responseRecorder := httptest.NewRecorder()

	handler.List(responseRecorder, request)

	if responseRecorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "{\"error\":{\"code\":\"internal_error\",\"message\":\"the product request could not be completed\"}}\n" {
		t.Fatalf("unexpected internal error response %q", responseRecorder.Body.String())
	}
}
