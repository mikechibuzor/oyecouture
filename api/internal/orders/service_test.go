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

func TestCreateRequiresDeliveryDetailsForBespokeOrder(t *testing.T) {
	t.Parallel()

	repository := &orderRepositoryStub{products: map[string]CatalogueProduct{
		"gold-agbada": {
			ID:        "product-1",
			Slug:      "gold-agbada",
			Name:      "Gold Agbada",
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
