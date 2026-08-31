package orders

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	ErrNoItems           = errors.New("order has no items")
	ErrTooManyItems      = errors.New("order has too many items")
	ErrInvalidCustomer   = errors.New("customer details are invalid")
	ErrInvalidItem       = errors.New("order item is invalid")
	ErrUnknownProduct    = errors.New("product is unavailable")
	ErrSizeRequired      = errors.New("size is required for this product")
	ErrSizeNotApplicable = errors.New("size is not applicable to this product")
	ErrDeliveryRequired  = errors.New("delivery details are required")
	ErrUnknownZone       = errors.New("delivery zone is not available")
)

const (
	maximumItemsPerOrder    = 20
	maximumNameLength       = 120
	maximumNotesLength      = 1000
	maximumMeasurementValue = 200
	payTokenBytes           = 24
)

var (
	productSlugPattern       = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	emailPattern             = regexp.MustCompile(`^[^@\s]+@[^@\s.]+\.[^@\s]+$`)
	phonePattern             = regexp.MustCompile(`^\+?[0-9][0-9\s-]{6,19}$`)
	validSizes               = map[string]struct{}{"S": {}, "M": {}, "L": {}, "XL": {}}
	validMeasurementProfiles = map[string]struct{}{"suit": {}, "native": {}, "mixed": {}}
	suitCategoryTerms        = []string{"suit", "tuxedo", "blazer", "two-piece", "three-piece"}
	nativeCategoryTerms      = []string{"agbada", "aso-oke", "kaftan", "buba", "native"}
	validMeasurementFields   = map[string]struct{}{
		"back": {}, "sleeve": {}, "round-sleeve": {}, "chest": {}, "stomach": {},
		"jacket-length": {}, "top-length": {}, "trouser-length": {}, "waist": {},
		"thigh": {}, "calf": {}, "ankle": {}, "neck": {}, "cap": {},
	}
	requiredMeasurementFields = map[string][]string{
		"suit": {
			"back", "sleeve", "round-sleeve", "chest", "stomach", "jacket-length",
			"trouser-length", "waist", "thigh", "calf", "ankle", "neck",
		},
		"native": {
			"back", "sleeve", "round-sleeve", "chest", "stomach", "top-length",
			"trouser-length", "waist", "thigh", "calf", "ankle", "neck", "cap",
		},
		"mixed": {
			"back", "sleeve", "round-sleeve", "chest", "stomach", "jacket-length",
			"top-length", "trouser-length", "waist", "thigh", "calf", "ankle", "neck", "cap",
		},
	}
)

// CatalogueProduct is the trusted product record used for pricing.
type CatalogueProduct struct {
	ID        string
	Slug      string
	Name      string
	Section   string
	Category  string
	PriceKobo int64
	HasSizes  bool
}

// Repository is the persistence contract for order creation.
// DeliveryRates resolves the fee for a zone, so checkout cannot set its own.
type DeliveryRates interface {
	FeeForZone(context.Context, string) (int64, bool, error)
}

type Repository interface {
	ProductsBySlugs(context.Context, []string) (map[string]CatalogueProduct, error)
	Create(context.Context, *Order, []Item, *Measurements, CreateOrderInput) error
}

type Service struct {
	repository    Repository
	deliveryRates DeliveryRates
}

func NewService(repository Repository, deliveryRates DeliveryRates) *Service {
	return &Service{repository: repository, deliveryRates: deliveryRates}
}

// Create validates the request, re-prices every line from the database and
// persists the order. Client-supplied prices are never used: the storefront
// only says which product, which size and how many.
func (service *Service) Create(
	requestContext context.Context,
	input CreateOrderInput,
) (Order, error) {
	if err := validateCustomer(input); err != nil {
		return Order{}, err
	}
	if len(input.Items) == 0 {
		return Order{}, ErrNoItems
	}
	if len(input.Items) > maximumItemsPerOrder {
		return Order{}, ErrTooManyItems
	}

	// Collect the slugs so the whole bag is priced in one query.
	requestedSlugs := make([]string, 0, len(input.Items))
	for _, item := range input.Items {
		if !productSlugPattern.MatchString(item.Slug) {
			return Order{}, ErrInvalidItem
		}
		if item.Quantity < 1 || item.Quantity > MaximumItemQuantity {
			return Order{}, ErrInvalidItem
		}
		requestedSlugs = append(requestedSlugs, item.Slug)
	}

	catalogue, err := service.repository.ProductsBySlugs(requestContext, requestedSlugs)
	if err != nil {
		return Order{}, fmt.Errorf("load products: %w", err)
	}

	// Build priced lines and decide how the order must be paid for.
	items := make([]Item, 0, len(input.Items))
	var subtotalKobo int64
	orderKind := KindReadyToWear
	bespokeProfiles := make(map[string]struct{})

	for _, requested := range input.Items {
		product, isKnown := catalogue[requested.Slug]
		if !isKnown {
			return Order{}, ErrUnknownProduct
		}

		size := strings.TrimSpace(requested.Size)
		if product.HasSizes {
			if _, isValidSize := validSizes[size]; !isValidSize {
				return Order{}, ErrSizeRequired
			}
		} else {
			if requested.Quantity != 1 {
				return Order{}, ErrInvalidItem
			}
			if size != "" {
				return Order{}, ErrSizeNotApplicable
			}
			// Anything without a size is made to measure, so the whole order
			// uses the bespoke confirmation and measurement rules.
			orderKind = KindBespoke
			bespokeProfiles[measurementProfileForProduct(product)] = struct{}{}
		}

		lineKobo := product.PriceKobo * int64(requested.Quantity)
		subtotalKobo += lineKobo

		items = append(items, Item{
			Slug:      product.Slug,
			Name:      product.Name,
			Size:      size,
			Quantity:  requested.Quantity,
			PriceKobo: product.PriceKobo,
			LineKobo:  lineKobo,
			productID: product.ID,
		})
	}

	// Every order ships to a configured zone. The server resolves the fee so a
	// client can never suppress or replace it, including on bespoke orders.
	deliveryZone := strings.TrimSpace(input.DeliveryZone)
	address := strings.TrimSpace(input.Address)

	if address == "" || deliveryZone == "" {
		return Order{}, ErrDeliveryRequired
	}

	deliveryFeeKobo, zoneExists, err := service.deliveryRates.FeeForZone(
		requestContext,
		deliveryZone,
	)
	if err != nil {
		return Order{}, fmt.Errorf("resolve delivery zone: %w", err)
	}
	if !zoneExists {
		return Order{}, ErrUnknownZone
	}

	payToken, err := newPayToken()
	if err != nil {
		return Order{}, fmt.Errorf("generate pay token: %w", err)
	}

	// Complete measurements make a bespoke order immediately payable. Orders
	// without them remain available for the existing WhatsApp confirmation flow.
	sanitizedMeasurements := sanitizeMeasurements(input.Measurements)
	status := StatusAwaitingPayment
	if orderKind == KindBespoke && !hasCompleteMeasurements(
		sanitizedMeasurements,
		requiredProfile(bespokeProfiles),
	) {
		status = StatusPendingConfirmation
	}

	order := Order{
		Kind:            orderKind,
		Status:          status,
		CustomerName:    strings.TrimSpace(input.CustomerName),
		Items:           items,
		SubtotalKobo:    subtotalKobo,
		DeliveryZone:    deliveryZone,
		DeliveryFeeKobo: deliveryFeeKobo,
		TotalKobo:       subtotalKobo + deliveryFeeKobo,
		PayToken:        payToken,
	}

	if err := service.repository.Create(
		requestContext, &order, items, sanitizedMeasurements, input,
	); err != nil {
		return Order{}, fmt.Errorf("create order: %w", err)
	}

	return order, nil
}

// validateCustomer checks the contact details the owner needs to reach them.
func validateCustomer(input CreateOrderInput) error {
	name := strings.TrimSpace(input.CustomerName)
	email := strings.TrimSpace(input.Email)
	phone := strings.TrimSpace(input.Phone)

	if name == "" || utf8.RuneCountInString(name) > maximumNameLength {
		return ErrInvalidCustomer
	}
	if !emailPattern.MatchString(email) || len(email) > 254 {
		return ErrInvalidCustomer
	}
	if !phonePattern.MatchString(phone) {
		return ErrInvalidCustomer
	}

	return nil
}

// sanitizeMeasurements drops an all-empty submission and trims free text.
func sanitizeMeasurements(measurements *Measurements) *Measurements {
	if measurements == nil {
		return nil
	}

	notes := strings.TrimSpace(measurements.Notes)
	if utf8.RuneCountInString(notes) > maximumNotesLength {
		notes = string([]rune(notes)[:maximumNotesLength])
	}
	measurements.Notes = notes

	_, validProfile := validMeasurementProfiles[measurements.Profile]
	if !validProfile || measurements.Unit != "in" {
		measurements.Profile = ""
		measurements.Unit = ""
		measurements.Values = nil
	} else {
		sanitizedValues := make(map[string]float64, len(measurements.Values))
		for field, value := range measurements.Values {
			_, validField := validMeasurementFields[field]
			if validField && !math.IsNaN(value) && !math.IsInf(value, 0) && value > 0 && value <= maximumMeasurementValue {
				sanitizedValues[field] = value
			}
		}
		if len(sanitizedValues) == 0 {
			measurements.Values = nil
		} else {
			measurements.Values = sanitizedValues
		}
	}

	hasAnyValue := notes != "" || len(measurements.Values) > 0
	for _, value := range []*float64{
		measurements.NeckCm, measurements.ChestCm, measurements.WaistCm,
		measurements.HipsCm, measurements.SleeveCm, measurements.InseamCm,
		measurements.LengthCm, measurements.HeightCm, measurements.WeightKg,
	} {
		if value != nil {
			hasAnyValue = true
			break
		}
	}

	if !hasAnyValue {
		return nil
	}

	return measurements
}

// hasCompleteMeasurements keeps payment eligibility server-authoritative by
// requiring every field defined for the submitted garment profile.
func hasCompleteMeasurements(measurements *Measurements, expectedProfile string) bool {
	if measurements == nil || measurements.Unit != "in" || measurements.Profile != expectedProfile {
		return false
	}

	requiredFields, profileExists := requiredMeasurementFields[measurements.Profile]
	if !profileExists {
		return false
	}

	for _, field := range requiredFields {
		if _, fieldExists := measurements.Values[field]; !fieldExists {
			return false
		}
	}

	return true
}

// requiredProfile collapses one or more trusted product profiles into the
// worksheet profile the customer must complete.
func requiredProfile(profiles map[string]struct{}) string {
	if len(profiles) != 1 {
		return "mixed"
	}

	for profile := range profiles {
		return profile
	}

	return "mixed"
}

// measurementProfileForProduct derives tailoring requirements from trusted
// catalogue fields rather than accepting a client-selected shortcut.
func measurementProfileForProduct(product CatalogueProduct) string {
	if product.Section == "suits" {
		return "suit"
	}
	if product.Section == "trad" {
		return "native"
	}

	description := strings.ToLower(product.Category + " " + product.Name)
	for _, term := range suitCategoryTerms {
		if strings.Contains(description, term) {
			return "suit"
		}
	}
	for _, term := range nativeCategoryTerms {
		if strings.Contains(description, term) {
			return "native"
		}
	}

	return "mixed"
}

// newPayToken generates the unguessable token behind the /pay/<token> link.
func newPayToken() (string, error) {
	tokenBytes := make([]byte, payTokenBytes)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(tokenBytes), nil
}
