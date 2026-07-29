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
