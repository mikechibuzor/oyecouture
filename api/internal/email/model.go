package email

import "context"

const (
	ProviderSendByte = "sendbyte"

	TypeOrderConfirmation = "order_confirmation"
	TypePaymentReceipt    = "payment_receipt"

	StatusPending = "pending"
	StatusSent    = "sent"
	StatusFailed  = "failed"
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
	OrderID        string
	Type           string
	Provider       string
	IdempotencyKey string
	RecipientEmail string
	Subject        string
}

type Order struct {
	ID              string
	OrderNo         string
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
