package products

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalidProductSection = errors.New("invalid product section")
	ErrInvalidProductSlug    = errors.New("invalid product slug")
	ErrInvalidActiveFilter   = errors.New("invalid product active filter")
	ErrInvalidProductDraft   = errors.New("invalid product draft")
	ErrDuplicateProductSlug  = errors.New("duplicate product slug")
	ErrProductNotFound       = errors.New("product not found")
	ErrSearchQueryTooLong    = errors.New("search query is too long")
)

var productSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

var validProductSections = map[Section]struct{}{
	SectionSuits:        {},
	SectionTrad:         {},
	SectionCasualsMen:   {},
	SectionCasualsWomen: {},
}

// Service validates storefront inputs and prepares repository products for delivery.
type Service struct {
	productRepository Repository
	r2PublicBase      string
}

// AdminRepository is the persistence contract for owner-facing product management.
type AdminRepository interface {
	ListAdmin(context.Context, AdminListFilter) (AdminProductPage, error)
	GetAdminBySlug(context.Context, string) (AdminProduct, error)
	CreateAdmin(context.Context, AdminProductDraft) (AdminProduct, error)
	UpdateAdmin(context.Context, string, AdminProductDraft) (AdminProduct, error)
	UpdateActive(context.Context, string, bool) error
}

// AdminService validates owner-facing catalogue management inputs.
type AdminService struct {
	productRepository AdminRepository
	r2PublicBase      string
}

// NewService creates a product service with a normalized public image base URL.
func NewService(productRepository Repository, r2PublicBase string) *Service {
	return &Service{
		productRepository: productRepository,
		r2PublicBase:      strings.TrimRight(strings.TrimSpace(r2PublicBase), "/"),
	}
}

// NewAdminService creates the owner-facing catalogue service.
func NewAdminService(productRepository AdminRepository, r2PublicBase string) *AdminService {
	return &AdminService{
		productRepository: productRepository,
		r2PublicBase:      strings.TrimRight(strings.TrimSpace(r2PublicBase), "/"),
	}
}

// List validates optional filters and returns storefront-ready active products.
func (service *Service) List(
	requestContext context.Context,
	input ListInput,
) ([]Product, error) {
	// Normalize and validate all user-provided filters before database access.
	productSection, err := parseProductSection(input.Section)
	if err != nil {
		return nil, err
	}
	searchQuery := strings.TrimSpace(input.Query)
	if utf8.RuneCountInString(searchQuery) > MaximumSearchLength {
		return nil, ErrSearchQueryTooLong
	}

	// Query active products with the validated repository filter.
	products, err := service.productRepository.List(requestContext, ListFilter{
		Section: productSection,
		Query:   searchQuery,
	})
	if err != nil {
		return nil, err
	}

	// Convert internal R2 keys into public URLs before returning the response model.
	service.attachPublicImageURLs(products)
	return products, nil
}

// GetBySlug validates a canonical slug and returns one storefront-ready product.
func (service *Service) GetBySlug(
	requestContext context.Context,
	productSlug string,
) (Product, error) {
	if !productSlugPattern.MatchString(productSlug) {
		return Product{}, ErrInvalidProductSlug
	}

	product, err := service.productRepository.GetBySlug(requestContext, productSlug)
	if err != nil {
		return Product{}, err
	}

	products := []Product{product}
	service.attachPublicImageURLs(products)

	return products[0], nil
}

// List validates admin catalogue filters and returns products including archived rows.
func (service *AdminService) List(
	requestContext context.Context,
	input AdminListInput,
) (AdminProductPage, error) {
	productSection, err := parseProductSection(input.Section)
	if err != nil {
		return AdminProductPage{}, err
	}

	activeFilter, err := parseActiveFilter(input.Active)
	if err != nil {
		return AdminProductPage{}, err
	}

	searchQuery := strings.TrimSpace(input.Query)
	if utf8.RuneCountInString(searchQuery) > MaximumSearchLength {
		return AdminProductPage{}, ErrSearchQueryTooLong
	}

	page, err := service.productRepository.ListAdmin(requestContext, AdminListFilter{
		Section: productSection,
		Query:   searchQuery,
		Active:  activeFilter,
		Limit:   clampAdminProductLimit(input.Limit),
		Offset:  max(input.Offset, 0),
	})
	if err != nil {
		return AdminProductPage{}, err
	}

	service.attachAdminPublicImageURLs(page.Products)
	return page, nil
}

// GetBySlug returns one owner-facing product for editing.
func (service *AdminService) GetBySlug(
	requestContext context.Context,
	productSlug string,
) (AdminProduct, error) {
	normalizedSlug := strings.TrimSpace(productSlug)
	if !productSlugPattern.MatchString(normalizedSlug) {
		return AdminProduct{}, ErrInvalidProductSlug
	}

	product, err := service.productRepository.GetAdminBySlug(requestContext, normalizedSlug)
	if err != nil {
		return AdminProduct{}, err
	}

	products := []AdminProduct{product}
	service.attachAdminPublicImageURLs(products)

	return products[0], nil
}

// Create validates and inserts one product.
func (service *AdminService) Create(
	requestContext context.Context,
	draft AdminProductDraft,
) (AdminProduct, error) {
	normalizedDraft, err := normalizeProductDraft(draft)
	if err != nil {
		return AdminProduct{}, err
	}

	product, err := service.productRepository.CreateAdmin(requestContext, normalizedDraft)
	if err != nil {
		return AdminProduct{}, err
	}

	products := []AdminProduct{product}
	service.attachAdminPublicImageURLs(products)

	return products[0], nil
}

// Update validates and applies owner-editable product metadata.
func (service *AdminService) Update(
	requestContext context.Context,
	productSlug string,
	draft AdminProductDraft,
) (AdminProduct, error) {
	normalizedSlug := strings.TrimSpace(productSlug)
	if !productSlugPattern.MatchString(normalizedSlug) {
		return AdminProduct{}, ErrInvalidProductSlug
	}

	normalizedDraft, err := normalizeProductDraft(draft)
	if err != nil {
		return AdminProduct{}, err
	}

	product, err := service.productRepository.UpdateAdmin(
		requestContext,
		normalizedSlug,
		normalizedDraft,
	)
	if err != nil {
		return AdminProduct{}, err
	}

	products := []AdminProduct{product}
	service.attachAdminPublicImageURLs(products)

	return products[0], nil
}

// UpdateActive changes one product's public availability.
func (service *AdminService) UpdateActive(
	requestContext context.Context,
	productSlug string,
	update AdminActiveUpdate,
) error {
	normalizedSlug := strings.TrimSpace(productSlug)
	if !productSlugPattern.MatchString(normalizedSlug) {
		return ErrInvalidProductSlug
	}

	return service.productRepository.UpdateActive(requestContext, normalizedSlug, update.Active)
}

// normalizeProductDraft trims text and enforces admin-editable product invariants.
func normalizeProductDraft(draft AdminProductDraft) (AdminProductDraft, error) {
	draft.Slug = strings.TrimSpace(draft.Slug)
	draft.Name = strings.TrimSpace(draft.Name)
	draft.Section = strings.TrimSpace(draft.Section)
	draft.Category = strings.TrimSpace(draft.Category)
	draft.Description = strings.TrimSpace(draft.Description)
	draft.Tag = strings.TrimSpace(draft.Tag)

	if !productSlugPattern.MatchString(draft.Slug) {
		return AdminProductDraft{}, ErrInvalidProductSlug
	}

	section, err := parseProductSection(draft.Section)
	if err != nil {
		return AdminProductDraft{}, err
	}
	if section == "" ||
		draft.Name == "" ||
		draft.Category == "" ||
		draft.Description == "" ||
		draft.PriceKobo < 0 ||
		draft.SortOrder < 0 {
		return AdminProductDraft{}, ErrInvalidProductDraft
	}

	draft.Section = string(section)
	return draft, nil
}

// parseProductSection converts an optional section value into a validated enum.
func parseProductSection(sectionValue string) (Section, error) {
	normalizedSection := Section(strings.TrimSpace(sectionValue))
	if normalizedSection == "" {
		return "", nil
	}
	if _, isValidSection := validProductSections[normalizedSection]; !isValidSection {
		return "", ErrInvalidProductSection
	}

	return normalizedSection, nil
}

// parseActiveFilter reads optional active=true|false admin filters.
func parseActiveFilter(activeValue string) (*bool, error) {
	normalizedActive := strings.TrimSpace(activeValue)
	if normalizedActive == "" {
		return nil, nil
	}

	active, err := strconv.ParseBool(normalizedActive)
	if err != nil {
		return nil, ErrInvalidActiveFilter
	}

	return &active, nil
}

// clampAdminProductLimit bounds owner catalogue reads.
func clampAdminProductLimit(limit int) int {
	if limit <= 0 {
		return defaultAdminProductLimit
	}
	if limit > maximumAdminProductLimit {
		return maximumAdminProductLimit
	}

	return limit
}

// attachPublicImageURLs enriches every product image without exposing its storage key.
func (service *Service) attachPublicImageURLs(products []Product) {
	for productIndex := range products {
		for imageIndex := range products[productIndex].Images {
			products[productIndex].Images[imageIndex].URL = buildPublicImageURL(
				service.r2PublicBase,
				products[productIndex].Images[imageIndex].R2Key,
			)
		}
	}
}

// attachAdminPublicImageURLs enriches admin products without exposing storage keys.
func (service *AdminService) attachAdminPublicImageURLs(products []AdminProduct) {
	for productIndex := range products {
		for imageIndex := range products[productIndex].Images {
			products[productIndex].Images[imageIndex].URL = buildPublicImageURL(
				service.r2PublicBase,
				products[productIndex].Images[imageIndex].R2Key,
			)
		}
	}
}

// buildPublicImageURL joins a configured R2 public base with a normalized object key.
func buildPublicImageURL(r2PublicBase string, r2Key string) string {
	if r2PublicBase == "" || strings.TrimSpace(r2Key) == "" {
		return ""
	}

	return r2PublicBase + "/" + strings.TrimLeft(r2Key, "/")
}
