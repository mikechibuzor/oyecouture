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
	listFilter            AdminListFilter
	listPage              AdminProductPage
	listError             error
	updateActiveSlug      string
	updateActiveValue     bool
	updateActiveError     error
	listCallCount         int
	updateActiveCallCount int
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
