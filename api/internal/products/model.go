package products

import "time"

// Section identifies a storefront catalogue section.
type Section string

const (
	SectionSuits             Section = "suits"
	SectionTrad              Section = "trad"
	SectionCasualsMen        Section = "casuals-men"
	SectionCasualsWomen      Section = "casuals-women"
	MaximumSearchLength              = 100
	defaultAdminProductLimit         = 25
	maximumAdminProductLimit         = 100
)

// Image represents an ordered public product image.
type Image struct {
	ID       string `json:"id"`
	R2Key    string `json:"-"`
	URL      string `json:"url"`
	Position int32  `json:"position"`
}

// Product represents a storefront-safe catalogue product.
type Product struct {
	ID          string    `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Section     Section   `json:"section"`
	Category    string    `json:"category"`
	PriceKobo   int64     `json:"priceKobo"`
	Description string    `json:"description"`
	Tag         string    `json:"tag"`
	HasSizes    bool      `json:"hasSizes"`
	SortOrder   int32     `json:"sortOrder"`
	CreatedAt   time.Time `json:"createdAt"`
	Images      []Image   `json:"images"`
}

// ListInput contains untrusted storefront product filters.
type ListInput struct {
	Section string
	Query   string
}

// ListFilter contains validated repository product filters.
type ListFilter struct {
	Section Section
	Query   string
}

// AdminProduct represents catalogue data the owner can manage.
type AdminProduct struct {
	ID          string    `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Section     Section   `json:"section"`
	Category    string    `json:"category"`
	PriceKobo   int64     `json:"priceKobo"`
	Description string    `json:"description"`
	Tag         string    `json:"tag"`
	HasSizes    bool      `json:"hasSizes"`
	Active      bool      `json:"active"`
	SortOrder   int32     `json:"sortOrder"`
	CreatedAt   time.Time `json:"createdAt"`
	Images      []Image   `json:"images"`
}

// AdminListInput contains untrusted admin catalogue filters.
type AdminListInput struct {
	Section string
	Query   string
	Active  string
	Limit   int
	Offset  int
}

// AdminListFilter contains normalized admin catalogue filters.
type AdminListFilter struct {
	Section Section
	Query   string
	Active  *bool
	Limit   int
	Offset  int
}

// AdminProductPage is a paged admin catalogue response.
type AdminProductPage struct {
	Products []AdminProduct `json:"products"`
	Total    int            `json:"total"`
	Limit    int            `json:"limit"`
	Offset   int            `json:"offset"`
}

// AdminActiveUpdate is the availability mutation exposed to the owner.
type AdminActiveUpdate struct {
	Active bool `json:"active"`
}
