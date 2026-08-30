package orders

import (
	"context"
	"errors"
	"testing"
)

type orderRepositoryStub struct {
	products map[string]CatalogueProduct
	created  *Order
}

func (repository *orderRepositoryStub) ProductsBySlugs(
	context.Context,
	[]string,
) (map[string]CatalogueProduct, error) {
	return repository.products, nil
}

func (repository *orderRepositoryStub) Create(
	_ context.Context,
	order *Order,
	_ []Item,
	_ *Measurements,
	_ CreateOrderInput,
) error {
	repository.created = order
	return nil
}

type deliveryRatesStub struct {
	feeKobo int64
	exists  bool
	zone    string
}

func (rates *deliveryRatesStub) FeeForZone(
	_ context.Context,
	zone string,
) (int64, bool, error) {
	rates.zone = zone
	return rates.feeKobo, rates.exists, nil
}

func TestCreateAppliesDeliveryFeeToBespokeOrder(t *testing.T) {
	t.Parallel()

	repository := &orderRepositoryStub{products: map[string]CatalogueProduct{
		"gold-agbada": {
			ID:        "product-1",
			Slug:      "gold-agbada",
			Name:      "Gold Agbada",
			Section:   "trad",
			Category:  "Agbada",
			PriceKobo: 78000000,
			HasSizes:  false,
		},
	}}
	rates := &deliveryRatesStub{feeKobo: 350000, exists: true}
	service := NewService(repository, rates)

	order, err := service.Create(context.Background(), validCreateOrderInput())
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if order.Kind != KindBespoke || order.Status != StatusPendingConfirmation {
		t.Fatalf("kind/status = %q/%q, want bespoke/pending_confirmation", order.Kind, order.Status)
	}
	if order.DeliveryZone != "Ibadan" || rates.zone != "Ibadan" {
		t.Fatalf("delivery zone = %q, resolved zone = %q, want Ibadan", order.DeliveryZone, rates.zone)
	}
	if order.DeliveryFeeKobo != 350000 || order.TotalKobo != 78350000 {
		t.Fatalf(
			"delivery/total = %d/%d, want 350000/78350000",
			order.DeliveryFeeKobo,
			order.TotalKobo,
		)
	}
	if repository.created == nil || repository.created.DeliveryFeeKobo != order.DeliveryFeeKobo {
		t.Fatal("Create() did not persist the server-resolved delivery fee")
	}
}

func TestCreateMakesBespokeOrderPayableWithCompleteMeasurements(t *testing.T) {
	t.Parallel()

	repository := &orderRepositoryStub{products: map[string]CatalogueProduct{
		"gold-agbada": {
			ID:        "product-1",
			Slug:      "gold-agbada",
			Name:      "Gold Agbada",
			Section:   "trad",
			Category:  "Agbada",
			PriceKobo: 78000000,
			HasSizes:  false,
		},
	}}
	service := NewService(repository, &deliveryRatesStub{feeKobo: 350000, exists: true})
	input := validCreateOrderInput()
	input.Measurements = completeNativeMeasurements()

	order, err := service.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if order.Kind != KindBespoke || order.Status != StatusAwaitingPayment {
		t.Fatalf("kind/status = %q/%q, want bespoke/awaiting_payment", order.Kind, order.Status)
	}
}

func TestCreateKeepsPartiallyMeasuredBespokeOrderPending(t *testing.T) {
	t.Parallel()

	repository := &orderRepositoryStub{products: map[string]CatalogueProduct{
		"gold-agbada": {
			ID:        "product-1",
			Slug:      "gold-agbada",
			Name:      "Gold Agbada",
			Section:   "trad",
			Category:  "Agbada",
			PriceKobo: 78000000,
			HasSizes:  false,
		},
	}}
	service := NewService(repository, &deliveryRatesStub{feeKobo: 350000, exists: true})
	input := validCreateOrderInput()
	input.Measurements = &Measurements{
		Profile: "native",
		Unit:    "in",
		Values:  map[string]float64{"chest": 42},
	}

	order, err := service.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if order.Status != StatusPendingConfirmation {
		t.Fatalf("status = %q, want pending_confirmation", order.Status)
	}
}

func TestCreateRejectsCompleteMeasurementsForWrongProductProfile(t *testing.T) {
	t.Parallel()

	repository := &orderRepositoryStub{products: map[string]CatalogueProduct{
		"gold-agbada": {
			ID:        "product-1",
			Slug:      "gold-agbada",
			Name:      "Gold Agbada",
			Section:   "trad",
			Category:  "Agbada",
			PriceKobo: 78000000,
			HasSizes:  false,
		},
	}}
	service := NewService(repository, &deliveryRatesStub{feeKobo: 350000, exists: true})
	input := validCreateOrderInput()
	input.Measurements = completeNativeMeasurements()
	input.Measurements.Profile = "suit"

	order, err := service.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if order.Status != StatusPendingConfirmation {
		t.Fatalf("status = %q, want pending_confirmation", order.Status)
	}
}

func TestCreateRequiresDeliveryDetailsForBespokeOrder(t *testing.T) {
	t.Parallel()

	repository := &orderRepositoryStub{products: map[string]CatalogueProduct{
		"gold-agbada": {
			ID:        "product-1",
			Slug:      "gold-agbada",
			Name:      "Gold Agbada",
			Section:   "trad",
			Category:  "Agbada",
			PriceKobo: 78000000,
			HasSizes:  false,
		},
	}}
	rates := &deliveryRatesStub{feeKobo: 350000, exists: true}
	service := NewService(repository, rates)
	input := validCreateOrderInput()
	input.Address = ""

	_, err := service.Create(context.Background(), input)
	if !errors.Is(err, ErrDeliveryRequired) {
		t.Fatalf("Create() error = %v, want ErrDeliveryRequired", err)
	}
	if rates.zone != "" {
		t.Fatalf("delivery lookup ran for incomplete details with zone %q", rates.zone)
	}
}

func validCreateOrderInput() CreateOrderInput {
	return CreateOrderInput{
		CustomerName: "Chibuzor Iwuagwu",
		Email:        "iwuagwuchibuzor@gmail.com",
		Phone:        "+2348094939264",
		Address:      "12 Example Road",
		DeliveryZone: "Ibadan",
		Items: []RequestedItem{{
			Slug:     "gold-agbada",
			Quantity: 1,
		}},
	}
}

func completeNativeMeasurements() *Measurements {
	return &Measurements{
		Profile: "native",
		Unit:    "in",
		Values: map[string]float64{
			"back": 18, "sleeve": 25, "round-sleeve": 15, "chest": 42,
			"stomach": 40, "top-length": 35, "trouser-length": 42, "waist": 36,
			"thigh": 24, "calf": 16, "ankle": 14, "neck": 16, "cap": 23,
		},
	}
}
