package orders

import "time"

// Kind decides how an order is paid for. A bag containing any made-to-measure
// piece is confirmed on WhatsApp first, because measurements must be agreed
// before money moves.
type Kind string

const (
	KindBespoke     Kind = "bespoke"
	KindReadyToWear Kind = "ready-to-wear"
)

// Status is the order lifecycle. Bespoke orders start at PendingConfirmation;
// ready-to-wear orders start at AwaitingPayment.
type Status string

const (
	StatusPendingConfirmation Status = "pending_confirmation"
	StatusAwaitingPayment     Status = "awaiting_payment"
	StatusPaid                Status = "paid"
	StatusInProgress          Status = "in_progress"
	StatusShipped             Status = "shipped"
	StatusDelivered           Status = "delivered"
	StatusCancelled           Status = "cancelled"
	StatusRefunded            Status = "refunded"
)

// MaximumItemQuantity mirrors the storefront cart cap and the database CHECK.
const MaximumItemQuantity = 10

// RequestedItem is one line as submitted by the storefront. Only the slug, size
// and quantity are trusted — the price is always re-read from the database.
type RequestedItem struct {
	Slug     string `json:"slug"`
	Size     string `json:"size"`
	Quantity int    `json:"quantity"`
}

// Measurements are optional; customers who already know theirs can skip the
// back-and-forth on WhatsApp.
type Measurements struct {
	NeckCm   *float64 `json:"neckCm"`
	ChestCm  *float64 `json:"chestCm"`
	WaistCm  *float64 `json:"waistCm"`
	HipsCm   *float64 `json:"hipsCm"`
	SleeveCm *float64 `json:"sleeveCm"`
	InseamCm *float64 `json:"inseamCm"`
	LengthCm *float64 `json:"lengthCm"`
	HeightCm *float64 `json:"heightCm"`
	WeightKg *float64 `json:"weightKg"`
	Notes    string   `json:"notes"`
}

// CreateOrderInput is the untrusted request body.
type CreateOrderInput struct {
	CustomerName string          `json:"customerName"`
	Email        string          `json:"email"`
	Phone        string          `json:"phone"`
	Items        []RequestedItem `json:"items"`
	Measurements *Measurements   `json:"measurements"`
	// Ready-to-wear only: bespoke delivery is agreed on WhatsApp.
	Address      string `json:"address"`
	DeliveryZone string `json:"deliveryZone"`
}

// Item is a persisted order line, priced by the server.
type Item struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Size      string `json:"size,omitempty"`
	Quantity  int    `json:"quantity"`
	PriceKobo int64  `json:"priceKobo"`
	LineKobo  int64  `json:"lineKobo"`
	productID string
}

// Order is the created order as returned to the storefront.
type Order struct {
	OrderNo         string    `json:"orderNo"`
	Kind            Kind      `json:"kind"`
	Status          Status    `json:"status"`
	CustomerName    string    `json:"customerName"`
	Items           []Item    `json:"items"`
	SubtotalKobo    int64     `json:"subtotalKobo"`
	DeliveryZone    string    `json:"deliveryZone,omitempty"`
	DeliveryFeeKobo int64     `json:"deliveryFeeKobo"`
	TotalKobo       int64     `json:"totalKobo"`
	PayToken        string    `json:"payToken"`
	CreatedAt       time.Time `json:"createdAt"`
}
