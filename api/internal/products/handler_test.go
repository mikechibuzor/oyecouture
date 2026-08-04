package products

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

// adminProductManagerStub records owner-facing handler inputs.
type adminProductManagerStub struct {
	listInput          AdminListInput
	listPage           AdminProductPage
	listError          error
	detailSlug         string
	detailProduct      AdminProduct
	detailError        error
	createDraft        AdminProductDraft
	createProduct      AdminProduct
	createError        error
	updateSlug         string
	updateDraft        AdminProductDraft
	updateProduct      AdminProduct
	updateError        error
	updateActiveSlug   string
	updateActive       bool
	updateActiveError  error
	updateActiveCalled bool
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

// List records admin list input and returns the configured page.
func (manager *adminProductManagerStub) List(
	_ context.Context,
	input AdminListInput,
) (AdminProductPage, error) {
	manager.listInput = input
	return manager.listPage, manager.listError
}

// GetBySlug records admin detail input.
func (manager *adminProductManagerStub) GetBySlug(
	_ context.Context,
	productSlug string,
) (AdminProduct, error) {
	manager.detailSlug = productSlug
	return manager.detailProduct, manager.detailError
}

// Create records admin create input.
func (manager *adminProductManagerStub) Create(
	_ context.Context,
	draft AdminProductDraft,
) (AdminProduct, error) {
	manager.createDraft = draft
	return manager.createProduct, manager.createError
}

// Update records admin update input.
func (manager *adminProductManagerStub) Update(
	_ context.Context,
	productSlug string,
	draft AdminProductDraft,
) (AdminProduct, error) {
	manager.updateSlug = productSlug
	manager.updateDraft = draft
	return manager.updateProduct, manager.updateError
}

// UpdateActive records one admin availability mutation.
func (manager *adminProductManagerStub) UpdateActive(
	_ context.Context,
	productSlug string,
	update AdminActiveUpdate,
) error {
	manager.updateActiveCalled = true
	manager.updateActiveSlug = productSlug
	manager.updateActive = update.Active
	return manager.updateActiveError
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

// TestAdminHandlerListReturnsProductsAndFilters verifies owner list request translation.
func TestAdminHandlerListReturnsProductsAndFilters(t *testing.T) {
	t.Parallel()

	manager := &adminProductManagerStub{
		listPage: AdminProductPage{
			Products: []AdminProduct{{
				Slug:      "gold-agbada",
				Active:    false,
				PriceKobo: 78_000_000,
			}},
			Total:  1,
			Limit:  25,
			Offset: 0,
		},
	}
	handler := NewAdminHandler(manager)
	request := httptest.NewRequest(
		http.MethodGet,
		"/admin/products?section=trad&q=gold&active=false&limit=25&offset=0",
		nil,
	)
	responseRecorder := httptest.NewRecorder()

	handler.List(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}
	if manager.listInput.Section != "trad" || manager.listInput.Query != "gold" ||
		manager.listInput.Active != "false" || manager.listInput.Limit != 25 {
		t.Fatalf("unexpected list input: %+v", manager.listInput)
	}

	var response AdminProductPage
	if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode admin product page: %v", err)
	}
	if response.Total != 1 || len(response.Products) != 1 || response.Products[0].Active {
		t.Fatalf("unexpected admin product response: %+v", response)
	}
}

// TestAdminHandlerMapsInvalidActiveFilter verifies admin validation errors.
func TestAdminHandlerMapsInvalidActiveFilter(t *testing.T) {
	t.Parallel()

	handler := NewAdminHandler(&adminProductManagerStub{listError: ErrInvalidActiveFilter})
	request := httptest.NewRequest(http.MethodGet, "/admin/products?active=archived", nil)
	responseRecorder := httptest.NewRecorder()

	handler.List(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, responseRecorder.Code)
	}
	var response errorResponse
	if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if response.Error.Code != ErrorCodeInvalidActiveFilter {
		t.Fatalf("expected invalid-active code, got %q", response.Error.Code)
	}
}

// TestAdminHandlerGetBySlugReturnsProduct verifies edit-detail responses.
func TestAdminHandlerGetBySlugReturnsProduct(t *testing.T) {
	t.Parallel()

	manager := &adminProductManagerStub{
		detailProduct: AdminProduct{Slug: "gold-agbada", Active: true},
	}
	handler := NewAdminHandler(manager)
	request := httptest.NewRequest(http.MethodGet, "/admin/products/gold-agbada", nil)
	request.SetPathValue("slug", "gold-agbada")
	responseRecorder := httptest.NewRecorder()

	handler.GetBySlug(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}
	if manager.detailSlug != "gold-agbada" {
		t.Fatalf("expected detail slug, got %q", manager.detailSlug)
	}

	var response adminDetailResponse
	if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode admin detail response: %v", err)
	}
	if response.Product.Slug != "gold-agbada" || !response.Product.Active {
		t.Fatalf("unexpected admin detail response: %+v", response.Product)
	}
}

// TestAdminHandlerCreateReturnsCreatedProduct verifies product create decoding.
func TestAdminHandlerCreateReturnsCreatedProduct(t *testing.T) {
	t.Parallel()

	manager := &adminProductManagerStub{
		createProduct: AdminProduct{Slug: "linen-shirt", Active: true},
	}
	handler := NewAdminHandler(manager)
	request := httptest.NewRequest(http.MethodPost, "/admin/products",
		strings.NewReader(`{
			"slug":"linen-shirt",
			"name":"Linen Shirt",
			"section":"casuals-men",
			"category":"Shirts",
			"priceKobo":12000000,
			"description":"A linen shirt.",
			"tag":"",
			"hasSizes":true,
			"active":true,
			"sortOrder":3
		}`))
	responseRecorder := httptest.NewRecorder()

	handler.Create(responseRecorder, request)

	if responseRecorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, responseRecorder.Code)
	}
	if manager.createDraft.Slug != "linen-shirt" || manager.createDraft.PriceKobo != 12_000_000 {
		t.Fatalf("unexpected create draft: %+v", manager.createDraft)
	}
}

// TestAdminHandlerUpdateReturnsProduct verifies product update decoding.
func TestAdminHandlerUpdateReturnsProduct(t *testing.T) {
	t.Parallel()

	manager := &adminProductManagerStub{
		updateProduct: AdminProduct{Slug: "edited-shirt", Active: false},
	}
	handler := NewAdminHandler(manager)
	request := httptest.NewRequest(http.MethodPut, "/admin/products/linen-shirt",
		strings.NewReader(`{
			"slug":"edited-shirt",
			"name":"Edited Shirt",
			"section":"casuals-women",
			"category":"Sets",
			"priceKobo":10500000,
			"description":"Edited.",
			"tag":"New",
			"hasSizes":true,
			"active":false,
			"sortOrder":8
		}`))
	request.SetPathValue("slug", "linen-shirt")
	responseRecorder := httptest.NewRecorder()

	handler.Update(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}
	if manager.updateSlug != "linen-shirt" || manager.updateDraft.Slug != "edited-shirt" {
		t.Fatalf("unexpected update call: %+v", manager)
	}
}

// TestAdminHandlerCreateMapsDuplicateSlug verifies conflict mapping.
func TestAdminHandlerCreateMapsDuplicateSlug(t *testing.T) {
	t.Parallel()

	handler := NewAdminHandler(&adminProductManagerStub{createError: ErrDuplicateProductSlug})
	request := httptest.NewRequest(http.MethodPost, "/admin/products",
		strings.NewReader(`{
			"slug":"linen-shirt",
			"name":"Linen Shirt",
			"section":"casuals-men",
			"category":"Shirts",
			"priceKobo":12000000,
			"description":"A linen shirt.",
			"tag":"",
			"hasSizes":true,
			"active":true,
			"sortOrder":3
		}`))
	responseRecorder := httptest.NewRecorder()

	handler.Create(responseRecorder, request)

	if responseRecorder.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, responseRecorder.Code)
	}
	var response errorResponse
	if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if response.Error.Code != ErrorCodeDuplicateProductSlug {
		t.Fatalf("expected duplicate slug code, got %q", response.Error.Code)
	}
}

// TestAdminHandlerUpdateActiveReturnsNoContent verifies active-state mutation.
func TestAdminHandlerUpdateActiveReturnsNoContent(t *testing.T) {
	t.Parallel()

	manager := &adminProductManagerStub{}
	handler := NewAdminHandler(manager)
	request := httptest.NewRequest(http.MethodPatch, "/admin/products/gold-agbada/active",
		strings.NewReader(`{"active":false}`))
	request.SetPathValue("slug", "gold-agbada")
	responseRecorder := httptest.NewRecorder()

	handler.UpdateActive(responseRecorder, request)

	if responseRecorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, responseRecorder.Code)
	}
	if !manager.updateActiveCalled || manager.updateActiveSlug != "gold-agbada" || manager.updateActive {
		t.Fatalf("unexpected update call: %+v", manager)
	}
}

// TestAdminHandlerUpdateActiveRejectsInvalidBody verifies strict JSON decoding.
func TestAdminHandlerUpdateActiveRejectsInvalidBody(t *testing.T) {
	t.Parallel()

	manager := &adminProductManagerStub{}
	handler := NewAdminHandler(manager)
	request := httptest.NewRequest(http.MethodPatch, "/admin/products/gold-agbada/active",
		strings.NewReader(`{"active":false,"extra":true}`))
	request.SetPathValue("slug", "gold-agbada")
	responseRecorder := httptest.NewRecorder()

	handler.UpdateActive(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, responseRecorder.Code)
	}
	if manager.updateActiveCalled {
		t.Fatalf("expected invalid body to skip service call")
	}
}
