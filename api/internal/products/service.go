package products

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalidProductSection = errors.New("invalid product section")
	ErrInvalidProductSlug    = errors.New("invalid product slug")
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

// NewService creates a product service with a normalized public image base URL.
func NewService(productRepository Repository, r2PublicBase string) *Service {
	return &Service{
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

// buildPublicImageURL joins a configured R2 public base with a normalized object key.
func buildPublicImageURL(r2PublicBase string, r2Key string) string {
	if r2PublicBase == "" || strings.TrimSpace(r2Key) == "" {
		return ""
	}

	return r2PublicBase + "/" + strings.TrimLeft(r2Key, "/")
}
