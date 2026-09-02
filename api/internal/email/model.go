package email

import (
	"context"
	"encoding/json"
	"time"
)

const (
	ProviderSendByte = "sendbyte"

	TypeOrderConfirmation = "order_confirmation"
	TypePaymentReceipt    = "payment_receipt"
	TypeFulfillmentUpdate = "fulfillment_update"
	TypeOrderDispatched   = "order_dispatched"
	TypeOrderDelivered    = "order_delivered"

	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusSent       = "sent"
	StatusFailed     = "failed"
)

type Provider interface {
	Name() string
	Send(context.Context, Message) (SendResult, []byte, error)
}

type Message struct {
	From           string
	To             string
	Subject        string
	HTML           string
	Text           string
	IdempotencyKey string
}

type SendResult struct {
	MessageID string
}

type Event struct {
	OrderID            string
	Type               string
	Provider           string
	IdempotencyKey     string
	RecipientEmail     string
	Subject            string
	FulfillmentEventID string
	Payload            json.RawMessage
}

// QueuedEvent is one fulfillment email claimed by the retry worker.
type QueuedEvent struct {
	IdempotencyKey string
	Type           string
	RecipientEmail string
	Subject        string
	Payload        json.RawMessage
}

// FulfillmentPayload contains only the customer-facing data needed to render
// an order-progress email. It intentionally excludes measurements and payment data.
type FulfillmentPayload struct {
	CustomerName          string    `json:"customerName"`
	OrderNo               string    `json:"orderNo"`
	Stage                 string    `json:"stage"`
	StageLabel            string    `json:"stageLabel"`
	PublicNote            string    `json:"publicNote,omitempty"`
	TrackingURL           string    `json:"trackingUrl,omitempty"`
	CarrierName           string    `json:"carrierName,omitempty"`
	CarrierTrackingNumber string    `json:"carrierTrackingNumber,omitempty"`
	CarrierTrackingURL    string    `json:"carrierTrackingUrl,omitempty"`
	EstimatedDeliveryDate string    `json:"estimatedDeliveryDate,omitempty"`
	OccurredAt            time.Time `json:"occurredAt"`
}

type Order struct {
	ID      string
	OrderNo string
	// PayToken and the payment reference together rebuild the customer's own
	// confirmation URL, which is what the design's "View your order" button
	// points at. Without the reference that page cannot verify and shows a
	// failure to someone who has already paid, so both are required.
	PayToken        string
	OrderDate       string
	CustomerName    string
	Email           string
	Phone           string
	Address         string
	DeliveryZone    string
	SubtotalKobo    int64
	DeliveryFeeKobo int64
	TotalKobo       int64
	Items           []Item
	Payment         Payment
}

type Item struct {
	Name      string
	Size      string
	Quantity  int
	PriceKobo int64
	LineKobo  int64
}

type Payment struct {
	Provider  string
	Reference string
	PaidAt    string
}
