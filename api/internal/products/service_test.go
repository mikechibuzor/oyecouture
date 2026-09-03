package products

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// repositoryStub records service calls and returns controlled product results.
type repositoryStub struct {
	listFilter         ListFilter
	listProducts       []Product
	listError          error
	getProductSlug     string
	getProduct         Product
	getError           error
	listCallCount      int
	getBySlugCallCount int
}

// adminRepositoryStub records admin service calls and returns controlled results.
type adminRepositoryStub struct {
	deleteAdminSlug       string
	deleteAdminCallCount  int
	deleteAdminError      error
	listFilter            AdminListFilter
	listPage              AdminProductPage
	listError             error
	getSlug               string
	getProduct            AdminProduct
	getError              error
	addImageSlug          string
	addImageKey           string
	addImageProduct       AdminProduct
	addImageError         error
	createDraft           AdminProductDraft
	createProduct         AdminProduct
	createError           error
	deleteImageSlug       string
	deleteImageID         string
	deleteImageProduct    AdminProduct
	deleteImageError      error
	reorderSlug           string
	reorderImageIDs       []string
	reorderProduct        AdminProduct
	reorderError          error
	updateSlug            string
	updateDraft           AdminProductDraft
	updateProduct         AdminProduct
	updateError           error
	updateActiveSlug      string
	updateActiveValue     bool
	updateActiveError     error
	listCallCount         int
	getCallCount          int
	addImageCallCount     int
	createCallCount       int
	deleteImageCallCount  int
	reorderCallCount      int
	updateCallCount       int
	updateActiveCallCount int
}

type imageUploaderStub struct {
	productSlug string
	upload      AdminImageUpload
	imageKey    string
	err         error
	callCount   int
}

// List records the validated filter and returns the configured products.
func (repository *repositoryStub) List(
	_ context.Context,
	filter ListFilter,
) ([]Product, error) {
	repository.listCallCount++
	repository.listFilter = filter
	return repository.listProducts, repository.listError
}

// GetBySlug records the validated slug and returns the configured product.
func (repository *repositoryStub) GetBySlug(
	_ context.Context,
	productSlug string,
) (Product, error) {
	repository.getBySlugCallCount++
	repository.getProductSlug = productSlug
	return repository.getProduct, repository.getError
}

// ListAdmin records the normalized admin filter.
func (repository *adminRepositoryStub) ListAdmin(
	_ context.Context,
	filter AdminListFilter,
) (AdminProductPage, error) {
	repository.listCallCount++
	repository.listFilter = filter
	return repository.listPage, repository.listError
}

// GetAdminBySlug records the edit-detail slug.
func (repository *adminRepositoryStub) GetAdminBySlug(
	_ context.Context,
	productSlug string,
) (AdminProduct, error) {
	repository.getCallCount++
	repository.getSlug = productSlug
	return repository.getProduct, repository.getError
}

// AddImage records an appended image key.
func (repository *adminRepositoryStub) AddImage(
	_ context.Context,
	productSlug string,
	imageKey string,
) (AdminProduct, error) {
	repository.addImageCallCount++
	repository.addImageSlug = productSlug
	repository.addImageKey = imageKey
	return repository.addImageProduct, repository.addImageError
}

// CreateAdmin records the normalized create draft.
func (repository *adminRepositoryStub) CreateAdmin(
	_ context.Context,
	draft AdminProductDraft,
) (AdminProduct, error) {
	repository.createCallCount++
	repository.createDraft = draft
	return repository.createProduct, repository.createError
}

// DeleteImage records an image deletion.
func (repository *adminRepositoryStub) DeleteImage(
	_ context.Context,
	productSlug string,
	imageID string,
) (AdminProduct, error) {
	repository.deleteImageCallCount++
	repository.deleteImageSlug = productSlug
	repository.deleteImageID = imageID
	return repository.deleteImageProduct, repository.deleteImageError
}

// ReorderImages records a full image order mutation.
func (repository *adminRepositoryStub) ReorderImages(
	_ context.Context,
	productSlug string,
	imageIDs []string,
) (AdminProduct, error) {
	repository.reorderCallCount++
	repository.reorderSlug = productSlug
	repository.reorderImageIDs = imageIDs
	return repository.reorderProduct, repository.reorderError
}

// UpdateAdmin records the normalized update draft.
func (repository *adminRepositoryStub) UpdateAdmin(
	_ context.Context,
	productSlug string,
	draft AdminProductDraft,
) (AdminProduct, error) {
	repository.updateCallCount++
	repository.updateSlug = productSlug
	repository.updateDraft = draft
	return repository.updateProduct, repository.updateError
}

// DeleteAdmin records the permanent removal.
func (repository *adminRepositoryStub) DeleteAdmin(
	_ context.Context,
	productSlug string,
) error {
	repository.deleteAdminCallCount++
	repository.deleteAdminSlug = productSlug
	return repository.deleteAdminError
}

// UpdateActive records the active-state mutation.
func (repository *adminRepositoryStub) UpdateActive(
	_ context.Context,
	productSlug string,
	active bool,
) error {
	repository.updateActiveCallCount++
	repository.updateActiveSlug = productSlug
	repository.updateActiveValue = active
	return repository.updateActiveError
}

// Upload records one provider upload and returns the configured key.
func (uploader *imageUploaderStub) Upload(
	_ context.Context,
	productSlug string,
	upload AdminImageUpload,
) (string, error) {
	uploader.callCount++
	uploader.productSlug = productSlug
	uploader.upload = upload
	return uploader.imageKey, uploader.err
}

// TestServiceListValidatesFiltersAndBuildsImageURLs verifies the full list walkthrough.
func TestServiceListValidatesFiltersAndBuildsImageURLs(t *testing.T) {
	t.Parallel()

	repository := &repositoryStub{
		listProducts: []Product{{
			Slug: "linen",
			Images: []Image{{
				R2Key: "products/linen/main.webp",
			}},
		}},
	}
	service := NewService(repository, " https://media.example.com/ ")

	products, err := service.List(context.Background(), ListInput{
		Section: " casuals-men ",
		Query:   " linen ",
	})
	if err != nil {
		t.Fatalf("expected list success, got %v", err)
	}

	if repository.listFilter.Section != SectionCasualsMen {
		t.Fatalf("expected normalized section, got %q", repository.listFilter.Section)
	}
	if repository.listFilter.Query != "linen" {
		t.Fatalf("expected normalized query, got %q", repository.listFilter.Query)
	}
	if products[0].Images[0].URL != "https://media.example.com/products/linen/main.webp" {
		t.Fatalf("unexpected public image URL %q", products[0].Images[0].URL)
	}
}

// TestServiceListRejectsInvalidSectionBeforeRepositoryAccess verifies section validation.
func TestServiceListRejectsInvalidSectionBeforeRepositoryAccess(t *testing.T) {
	t.Parallel()

	repository := &repositoryStub{}
	service := NewService(repository, "")

	_, err := service.List(context.Background(), ListInput{Section: "invalid"})
	if !errors.Is(err, ErrInvalidProductSection) {
		t.Fatalf("expected invalid section error, got %v", err)
	}
	if repository.listCallCount != 0 {
		t.Fatalf("expected no repository calls, got %d", repository.listCallCount)
	}
}

// TestServiceListRejectsLongSearchBeforeRepositoryAccess verifies the query limit.
func TestServiceListRejectsLongSearchBeforeRepositoryAccess(t *testing.T) {
	t.Parallel()

	repository := &repositoryStub{}
	service := NewService(repository, "")

	_, err := service.List(context.Background(), ListInput{
		Query: strings.Repeat("a", MaximumSearchLength+1),
	})
	if !errors.Is(err, ErrSearchQueryTooLong) {
		t.Fatalf("expected long search error, got %v", err)
	}
	if repository.listCallCount != 0 {
		t.Fatalf("expected no repository calls, got %d", repository.listCallCount)
	}
}

// TestServiceGetBySlugRejectsMalformedSlug verifies slug validation before persistence.
func TestServiceGetBySlugRejectsMalformedSlug(t *testing.T) {
	t.Parallel()

	repository := &repositoryStub{}
	service := NewService(repository, "")

	_, err := service.GetBySlug(context.Background(), "Royal Agbada")
	if !errors.Is(err, ErrInvalidProductSlug) {
		t.Fatalf("expected invalid slug error, got %v", err)
	}
	if repository.getBySlugCallCount != 0 {
		t.Fatalf("expected no repository calls, got %d", repository.getBySlugCallCount)
	}
}

// TestAdminServiceListNormalizesFiltersAndBuildsImageURLs verifies owner list validation.
func TestAdminServiceListNormalizesFiltersAndBuildsImageURLs(t *testing.T) {
	t.Parallel()

	repository := &adminRepositoryStub{
		listPage: AdminProductPage{
			Products: []AdminProduct{{
				Slug: "gold-agbada",
				Images: []Image{{
					R2Key: "/oyecouture/asoke-look-01-1",
				}},
			}},
		},
	}
	service := NewAdminService(repository, "https://res.cloudinary.com/demo/image/upload/")

	page, err := service.List(context.Background(), AdminListInput{
		Section: " trad ",
		Query:   " agbada ",
		Active:  "false",
		Limit:   500,
		Offset:  -10,
	})
	if err != nil {
		t.Fatalf("expected admin list success, got %v", err)
	}

	if repository.listFilter.Section != SectionTrad {
		t.Fatalf("expected trad section, got %q", repository.listFilter.Section)
	}
	if repository.listFilter.Query != "agbada" {
		t.Fatalf("expected trimmed query, got %q", repository.listFilter.Query)
	}
	if repository.listFilter.Active == nil || *repository.listFilter.Active {
		t.Fatalf("expected inactive filter, got %#v", repository.listFilter.Active)
	}
	if repository.listFilter.Limit != maximumAdminProductLimit {
		t.Fatalf("expected clamped limit, got %d", repository.listFilter.Limit)
	}
	if repository.listFilter.Offset != 0 {
		t.Fatalf("expected clamped offset, got %d", repository.listFilter.Offset)
	}
	if page.Products[0].Images[0].URL != "https://res.cloudinary.com/demo/image/upload/oyecouture/asoke-look-01-1" {
		t.Fatalf("unexpected admin image URL %q", page.Products[0].Images[0].URL)
	}
}

// TestAdminServiceListRejectsInvalidActiveFilter verifies active query validation.
func TestAdminServiceListRejectsInvalidActiveFilter(t *testing.T) {
	t.Parallel()

	repository := &adminRepositoryStub{}
	service := NewAdminService(repository, "")

	_, err := service.List(context.Background(), AdminListInput{Active: "archived"})
	if !errors.Is(err, ErrInvalidActiveFilter) {
		t.Fatalf("expected invalid active filter, got %v", err)
	}
	if repository.listCallCount != 0 {
		t.Fatalf("expected no repository calls, got %d", repository.listCallCount)
	}
}

// TestAdminServiceUpdateActiveValidatesSlug verifies mutation validation.
func TestAdminServiceUpdateActiveValidatesSlug(t *testing.T) {
	t.Parallel()

	repository := &adminRepositoryStub{}
	service := NewAdminService(repository, "")

	err := service.UpdateActive(context.Background(), " Gold Agbada ", AdminActiveUpdate{Active: false})
	if !errors.Is(err, ErrInvalidProductSlug) {
		t.Fatalf("expected invalid slug error, got %v", err)
	}
	if repository.updateActiveCallCount != 0 {
		t.Fatalf("expected no repository calls, got %d", repository.updateActiveCallCount)
	}
}

// TestAdminServiceUpdateActiveCallsRepository verifies valid mutations pass through.
func TestAdminServiceUpdateActiveCallsRepository(t *testing.T) {
	t.Parallel()

	repository := &adminRepositoryStub{}
	service := NewAdminService(repository, "")

	err := service.UpdateActive(context.Background(), " gold-agbada ", AdminActiveUpdate{Active: false})
	if err != nil {
		t.Fatalf("expected update success, got %v", err)
	}
	if repository.updateActiveSlug != "gold-agbada" {
		t.Fatalf("expected normalized slug, got %q", repository.updateActiveSlug)
	}
	if repository.updateActiveValue {
		t.Fatalf("expected active=false mutation")
	}
}

// TestAdminServiceAddImageUploadsAndPersistsKey verifies upload-provider wiring.
func TestAdminServiceAddImageUploadsAndPersistsKey(t *testing.T) {
	t.Parallel()

	repository := &adminRepositoryStub{
		getProduct: AdminProduct{Slug: "gold-agbada"},
		addImageProduct: AdminProduct{
			Slug: "gold-agbada",
			Images: []Image{{
				R2Key: "oyecouture/gold-agbada-1",
			}},
		},
	}
	uploader := &imageUploaderStub{imageKey: "oyecouture/gold-agbada-1"}
	service := NewAdminServiceWithUploader(repository, uploader, "https://res.cloudinary.com/demo/image/upload")

	product, err := service.AddImage(context.Background(), " gold-agbada ", AdminImageUpload{
		FileName:    "look.webp",
		ContentType: "image/webp",
		Content:     []byte("image"),
	})
	if err != nil {
		t.Fatalf("expected image upload success, got %v", err)
	}
	if uploader.productSlug != "gold-agbada" || repository.addImageSlug != "gold-agbada" {
		t.Fatalf("expected normalized slug, got uploader=%q repository=%q",
			uploader.productSlug, repository.addImageSlug)
	}
	if repository.addImageKey != "oyecouture/gold-agbada-1" {
		t.Fatalf("expected uploaded key to be persisted, got %q", repository.addImageKey)
	}
	if product.Images[0].URL != "https://res.cloudinary.com/demo/image/upload/oyecouture/gold-agbada-1" {
		t.Fatalf("expected public image URL, got %q", product.Images[0].URL)
	}
}

// TestAdminServiceAddImageRequiresUploader verifies config failures are explicit.
func TestAdminServiceAddImageRequiresUploader(t *testing.T) {
	t.Parallel()

	repository := &adminRepositoryStub{}
	service := NewAdminService(repository, "")

	_, err := service.AddImage(context.Background(), "gold-agbada", AdminImageUpload{
		FileName:    "look.jpg",
		ContentType: "image/jpeg",
		Content:     []byte("image"),
	})
	if !errors.Is(err, ErrImageUploadUnavailable) {
		t.Fatalf("expected unavailable upload error, got %v", err)
	}
	if repository.getCallCount != 0 {
		t.Fatalf("expected no repository read before uploader config, got %d", repository.getCallCount)
	}
}

// TestAdminServiceAddImageRejectsFullGallery verifies the 12-image cap.
func TestAdminServiceAddImageRejectsFullGallery(t *testing.T) {
	t.Parallel()

	repository := &adminRepositoryStub{
		getProduct: AdminProduct{
			Slug:   "gold-agbada",
			Images: make([]Image, maximumProductImages),
		},
	}
	uploader := &imageUploaderStub{imageKey: "oyecouture/gold-agbada-13"}
	service := NewAdminServiceWithUploader(repository, uploader, "")

	_, err := service.AddImage(context.Background(), "gold-agbada", AdminImageUpload{
		FileName:    "look.png",
		ContentType: "image/png",
		Content:     []byte("image"),
	})
	if !errors.Is(err, ErrImageLimitReached) {
		t.Fatalf("expected image limit error, got %v", err)
	}
	if uploader.callCount != 0 || repository.addImageCallCount != 0 {
		t.Fatalf("expected no upload or insert, got uploader=%d insert=%d",
			uploader.callCount, repository.addImageCallCount)
	}
}

// TestAdminServiceReorderImagesNormalizesIDs verifies image order validation.
func TestAdminServiceReorderImagesNormalizesIDs(t *testing.T) {
	t.Parallel()

	repository := &adminRepositoryStub{
		reorderProduct: AdminProduct{
			Slug: "gold-agbada",
			Images: []Image{{
				ID:    "image-2",
				R2Key: "oyecouture/two",
			}},
		},
	}
	service := NewAdminService(repository, "https://media.example.com")

	product, err := service.ReorderImages(context.Background(), " gold-agbada ", AdminImageReorder{
		ImageIDs: []string{" image-2 ", "image-1"},
	})
	if err != nil {
		t.Fatalf("expected reorder success, got %v", err)
	}
	if repository.reorderSlug != "gold-agbada" ||
		repository.reorderImageIDs[0] != "image-2" ||
		repository.reorderImageIDs[1] != "image-1" {
		t.Fatalf("unexpected reorder call: %+v", repository.reorderImageIDs)
	}
	if product.Images[0].URL != "https://media.example.com/oyecouture/two" {
		t.Fatalf("expected public image URL, got %q", product.Images[0].URL)
	}
}

// TestAdminServiceReorderImagesRejectsDuplicates verifies exact-order inputs.
func TestAdminServiceReorderImagesRejectsDuplicates(t *testing.T) {
	t.Parallel()

	repository := &adminRepositoryStub{}
	service := NewAdminService(repository, "")

	_, err := service.ReorderImages(context.Background(), "gold-agbada", AdminImageReorder{
		ImageIDs: []string{"image-1", "image-1"},
	})
	if !errors.Is(err, ErrInvalidImageReorder) {
		t.Fatalf("expected invalid reorder error, got %v", err)
	}
	if repository.reorderCallCount != 0 {
		t.Fatalf("expected no repository call, got %d", repository.reorderCallCount)
	}
}

// TestAdminServiceCreateNormalizesDraft verifies owner create validation.
func TestAdminServiceCreateNormalizesDraft(t *testing.T) {
	t.Parallel()

	repository := &adminRepositoryStub{
		createProduct: AdminProduct{Slug: "linen-shirt"},
	}
	service := NewAdminService(repository, "")

	_, err := service.Create(context.Background(), AdminProductDraft{
		Slug:        " linen-shirt ",
		Name:        " Linen Shirt ",
		Section:     string(SectionCasualsMen),
		Category:    " Shirts ",
		PriceKobo:   12_000_000,
		Description: " Lightweight linen \u2014 cut for warm days. ",
		Tag:         " New ",
		HasSizes:    true,
		Active:      true,
		SortOrder:   4,
	})
	if err != nil {
		t.Fatalf("expected create success, got %v", err)
	}
	if repository.createDraft.Slug != "linen-shirt" ||
		repository.createDraft.Name != "Linen Shirt" ||
		repository.createDraft.Category != "Shirts" ||
		repository.createDraft.Description != "Lightweight linen, cut for warm days." ||
		repository.createDraft.Tag != "New" {
		t.Fatalf("draft was not normalized: %+v", repository.createDraft)
	}
}

// TestAdminServiceCreateRejectsIncompleteDraft verifies required-field validation.
func TestAdminServiceCreateRejectsIncompleteDraft(t *testing.T) {
	t.Parallel()

	repository := &adminRepositoryStub{}
	service := NewAdminService(repository, "")

	_, err := service.Create(context.Background(), AdminProductDraft{
		Slug:      "linen-shirt",
		Section:   string(SectionCasualsMen),
		PriceKobo: -1,
	})
	if !errors.Is(err, ErrInvalidProductDraft) {
		t.Fatalf("expected invalid draft error, got %v", err)
	}
	if repository.createCallCount != 0 {
		t.Fatalf("expected no repository calls, got %d", repository.createCallCount)
	}
}

// TestAdminServiceUpdateNormalizesSlugAndDraft verifies update validation.
func TestAdminServiceUpdateNormalizesSlugAndDraft(t *testing.T) {
	t.Parallel()

	repository := &adminRepositoryStub{
		updateProduct: AdminProduct{Slug: "new-linen-shirt"},
	}
	service := NewAdminService(repository, "")

	_, err := service.Update(context.Background(), " linen-shirt ", AdminProductDraft{
		Slug:        "new-linen-shirt",
		Name:        "New Linen Shirt",
		Section:     string(SectionCasualsWomen),
		Category:    "Sets",
		PriceKobo:   10_000_000,
		Description: "Edited description.",
		Active:      true,
	})
	if err != nil {
		t.Fatalf("expected update success, got %v", err)
	}
	if repository.updateSlug != "linen-shirt" {
		t.Fatalf("expected normalized current slug, got %q", repository.updateSlug)
	}
	if repository.updateDraft.Section != string(SectionCasualsWomen) {
		t.Fatalf("expected normalized section, got %q", repository.updateDraft.Section)
	}
}

// TestSlugIsDerivedFromNameWhenBlank covers the optional slug. The owner should
// not have to invent a URL-safe string for every piece, and a hand-typed one is
// usually worse than the generated version.
func TestSlugIsDerivedFromNameWhenBlank(t *testing.T) {
	t.Parallel()

	testCases := []struct{ name, productName, wantSlug string }{
		{name: "spaces become hyphens", productName: "The Sovereign Suit", wantSlug: "the-sovereign-suit"},
		{name: "punctuation is dropped", productName: "Aso-Oke Two-Piece", wantSlug: "aso-oke-two-piece"},
		{name: "runs collapse", productName: "Ivory   &   Black Set", wantSlug: "ivory-black-set"},
		{name: "edges are trimmed", productName: "  Gold Agbada!  ", wantSlug: "gold-agbada"},
		{name: "digits survive", productName: "Look 12 Kaftan", wantSlug: "look-12-kaftan"},
		{name: "accents are dropped rather than mangled", productName: "Café Blazer", wantSlug: "caf-blazer"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			normalized, err := normalizeProductDraft(AdminProductDraft{
				Name:        testCase.productName,
				Section:     string(SectionSuits),
				Category:    "Two-Piece",
				Description: "A description.",
				PriceKobo:   100,
			})
			if err != nil {
				t.Fatalf("normalizeProductDraft() error = %v", err)
			}
			if normalized.Slug != testCase.wantSlug {
				t.Fatalf("slug = %q, want %q", normalized.Slug, testCase.wantSlug)
			}
		})
	}
}

// TestExplicitSlugIsKept confirms a supplied slug still wins, so existing
// products and deliberate URLs are untouched by the generator.
func TestExplicitSlugIsKept(t *testing.T) {
	t.Parallel()

	normalized, err := normalizeProductDraft(AdminProductDraft{
		Slug:        "chosen-by-hand",
		Name:        "Something Else Entirely",
		Section:     string(SectionSuits),
		Category:    "Two-Piece",
		Description: "A description.",
		PriceKobo:   100,
	})
	if err != nil {
		t.Fatalf("normalizeProductDraft() error = %v", err)
	}
	if normalized.Slug != "chosen-by-hand" {
		t.Fatalf("slug = %q, want it left alone", normalized.Slug)
	}
}

// TestNameWithNothingUsableIsRejected keeps a meaningless slug out of the URL
// space rather than storing an empty or hyphen-only one.
func TestNameWithNothingUsableIsRejected(t *testing.T) {
	t.Parallel()

	_, err := normalizeProductDraft(AdminProductDraft{
		Name:        "!!! ???",
		Section:     string(SectionSuits),
		Category:    "Two-Piece",
		Description: "A description.",
		PriceKobo:   100,
	})
	if !errors.Is(err, ErrInvalidProductSlug) {
		t.Fatalf("error = %v, want ErrInvalidProductSlug", err)
	}
}
