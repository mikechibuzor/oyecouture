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
	ErrInvalidProductSection  = errors.New("invalid product section")
	ErrInvalidProductSlug     = errors.New("invalid product slug")
	ErrInvalidActiveFilter    = errors.New("invalid product active filter")
	ErrInvalidProductDraft    = errors.New("invalid product draft")
	ErrInvalidImageUpload     = errors.New("invalid product image upload")
	ErrInvalidImageReorder    = errors.New("invalid product image reorder")
	ErrImageLimitReached      = errors.New("product image limit reached")
	ErrImageNotFound          = errors.New("product image not found")
	ErrImageUploadUnavailable = errors.New("product image upload unavailable")
	ErrDuplicateProductSlug   = errors.New("duplicate product slug")
	ErrProductNotFound        = errors.New("product not found")
	ErrSearchQueryTooLong     = errors.New("search query is too long")
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
	AddImage(context.Context, string, string) (AdminProduct, error)
	CreateAdmin(context.Context, AdminProductDraft) (AdminProduct, error)
	DeleteImage(context.Context, string, string) (AdminProduct, error)
	ReorderImages(context.Context, string, []string) (AdminProduct, error)
	UpdateAdmin(context.Context, string, AdminProductDraft) (AdminProduct, error)
	UpdateActive(context.Context, string, bool) error
	DeleteAdmin(context.Context, string) error
}

// ImageUploader persists uploaded image bytes and returns a public storage key.
type ImageUploader interface {
	Upload(context.Context, string, AdminImageUpload) (string, error)
}

// AdminService validates owner-facing catalogue management inputs.
type AdminService struct {
	productRepository AdminRepository
	imageUploader     ImageUploader
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

// NewAdminServiceWithUploader creates the owner-facing service with image uploads enabled.
func NewAdminServiceWithUploader(
	productRepository AdminRepository,
	imageUploader ImageUploader,
	r2PublicBase string,
) *AdminService {
	return &AdminService{
		productRepository: productRepository,
		imageUploader:     imageUploader,
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

// AddImage uploads one image and appends it to the product gallery.
// Delete removes a product outright.
//
// Order history survives this: order_items snapshots the name and price it was
// sold at, and its product_id is ON DELETE SET NULL, so a past order still
// reads correctly. product_images cascades, though the uploaded files stay in
// Cloudinary; nothing here reaches out to remove them.
//
// Archiving via UpdateActive remains the gentler option and is what the admin
// offers first, since a deleted product cannot be brought back.
func (service *AdminService) Delete(
	requestContext context.Context,
	productSlug string,
) error {
	productSlug = strings.TrimSpace(productSlug)
	if !productSlugPattern.MatchString(productSlug) {
		return ErrInvalidProductSlug
	}

	return service.productRepository.DeleteAdmin(requestContext, productSlug)
}

func (service *AdminService) AddImage(
	requestContext context.Context,
	productSlug string,
	upload AdminImageUpload,
) (AdminProduct, error) {
	normalizedSlug := strings.TrimSpace(productSlug)
	if !productSlugPattern.MatchString(normalizedSlug) {
		return AdminProduct{}, ErrInvalidProductSlug
	}
	if err := validateImageUpload(upload); err != nil {
		return AdminProduct{}, err
	}
	if service.imageUploader == nil {
		return AdminProduct{}, ErrImageUploadUnavailable
	}

	product, err := service.GetBySlug(requestContext, normalizedSlug)
	if err != nil {
		return AdminProduct{}, err
	}
	if len(product.Images) >= maximumProductImages {
		return AdminProduct{}, ErrImageLimitReached
	}

	imageKey, err := service.imageUploader.Upload(requestContext, normalizedSlug, upload)
	if err != nil {
		return AdminProduct{}, err
	}

	updatedProduct, err := service.productRepository.AddImage(
		requestContext,
		normalizedSlug,
		imageKey,
	)
	if err != nil {
		return AdminProduct{}, err
	}

	products := []AdminProduct{updatedProduct}
	service.attachAdminPublicImageURLs(products)

	return products[0], nil
}

// DeleteImage removes one image from the product and compacts positions.
func (service *AdminService) DeleteImage(
	requestContext context.Context,
	productSlug string,
	imageID string,
) (AdminProduct, error) {
	normalizedSlug := strings.TrimSpace(productSlug)
	normalizedImageID := strings.TrimSpace(imageID)
	if !productSlugPattern.MatchString(normalizedSlug) || normalizedImageID == "" {
		return AdminProduct{}, ErrInvalidProductSlug
	}

	product, err := service.productRepository.DeleteImage(
		requestContext,
		normalizedSlug,
		normalizedImageID,
	)
	if err != nil {
		return AdminProduct{}, err
	}

	products := []AdminProduct{product}
	service.attachAdminPublicImageURLs(products)

	return products[0], nil
}

// ReorderImages persists the owner-selected gallery order.
func (service *AdminService) ReorderImages(
	requestContext context.Context,
	productSlug string,
	reorder AdminImageReorder,
) (AdminProduct, error) {
	normalizedSlug := strings.TrimSpace(productSlug)
	if !productSlugPattern.MatchString(normalizedSlug) {
		return AdminProduct{}, ErrInvalidProductSlug
	}

	imageIDs, err := normalizeImageIDs(reorder.ImageIDs)
	if err != nil {
		return AdminProduct{}, err
	}

	product, err := service.productRepository.ReorderImages(
		requestContext,
		normalizedSlug,
		imageIDs,
	)
	if err != nil {
		return AdminProduct{}, err
	}

	products := []AdminProduct{product}
	service.attachAdminPublicImageURLs(products)

	return products[0], nil
}

// normalizeProductDraft trims text and enforces admin-editable product invariants.
func normalizeProductDraft(draft AdminProductDraft) (AdminProductDraft, error) {
	draft.Slug = strings.TrimSpace(draft.Slug)
	draft.Name = strings.TrimSpace(draft.Name)
	draft.Section = strings.TrimSpace(draft.Section)
	draft.Category = strings.TrimSpace(draft.Category)
	draft.Description = normalizeProductDescription(draft.Description)
	draft.Tag = strings.TrimSpace(draft.Tag)

	// The slug is optional: leaving it blank derives one from the name. Asking
	// the owner to invent a URL-safe string for every piece is a technical
	// detail leaking into a catalogue task, and a hand-typed one is usually
	// worse than the generated version anyway.
	if draft.Slug == "" {
		draft.Slug = slugFromName(draft.Name)
	}
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

// slugFromName builds a URL-safe slug from a product name: lowercase, runs of
// anything non-alphanumeric collapsed to single hyphens, no leading or trailing
// hyphen. Returns empty for a name with nothing usable in it, which the caller
// then rejects as an invalid slug rather than storing something meaningless.
func slugFromName(name string) string {
	var builder strings.Builder
	previousWasHyphen := true

	for _, character := range strings.ToLower(name) {
		switch {
		case (character >= 'a' && character <= 'z') || (character >= '0' && character <= '9'):
			builder.WriteRune(character)
			previousWasHyphen = false
		case !previousWasHyphen:
			builder.WriteByte('-')
			previousWasHyphen = true
		}
	}

	return strings.Trim(builder.String(), "-")
}

// normalizeProductDescription keeps catalogue prose aligned with the approved
// copy style while preserving ordinary hyphens inside compound words.
func normalizeProductDescription(description string) string {
	normalizedDescription := strings.TrimSpace(description)
	normalizedDescription = strings.ReplaceAll(normalizedDescription, " \u2014 ", ", ")
	return strings.ReplaceAll(normalizedDescription, "\u2014", "-")
}

// validateImageUpload keeps invalid or unsafe files away from the image provider.
func validateImageUpload(upload AdminImageUpload) error {
	if strings.TrimSpace(upload.FileName) == "" || len(upload.Content) == 0 {
		return ErrInvalidImageUpload
	}

	switch strings.ToLower(strings.TrimSpace(upload.ContentType)) {
	case "image/jpeg", "image/png", "image/webp":
		return nil
	default:
		return ErrInvalidImageUpload
	}
}

// normalizeImageIDs trims and de-duplicates a requested image order.
func normalizeImageIDs(imageIDs []string) ([]string, error) {
	if len(imageIDs) == 0 || len(imageIDs) > maximumProductImages {
		return nil, ErrInvalidImageReorder
	}

	seenImageIDs := make(map[string]struct{}, len(imageIDs))
	normalizedImageIDs := make([]string, 0, len(imageIDs))
	for _, imageID := range imageIDs {
		normalizedImageID := strings.TrimSpace(imageID)
		if normalizedImageID == "" {
			return nil, ErrInvalidImageReorder
		}
		if _, seen := seenImageIDs[normalizedImageID]; seen {
			return nil, ErrInvalidImageReorder
		}

		seenImageIDs[normalizedImageID] = struct{}{}
		normalizedImageIDs = append(normalizedImageIDs, normalizedImageID)
	}

	return normalizedImageIDs, nil
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
