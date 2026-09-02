package payments

import (
	"errors"
	"time"
)

const (
	ProviderPaystack = "paystack"
	CurrencyNGN      = "NGN"

	StatusInitializing = "initializing"
	StatusInitialized  = "initialized"
	StatusSuccess      = "success"
	StatusFailed       = "failed"
	StatusPaid         = "paid"

	OrderStatusAwaitingPayment = "awaiting_payment"
	OrderStatusPaid            = "paid"
)

var (
	ErrOrderNotFound         = errors.New("payment order not found")
	ErrOrderNotPayable       = errors.New("order is not payable")
	ErrPaymentAlreadyPaid    = errors.New("order is already paid")
	ErrPaymentInProgress     = errors.New("payment is already in progress")
	ErrPaymentUnavailable    = errors.New("payment provider is unavailable")
	ErrReferenceNotFound     = errors.New("payment reference not found")
	ErrReferenceMismatch     = errors.New("payment reference mismatch")
	ErrPaymentAmountMismatch = errors.New("payment amount mismatch")
	ErrPaymentCurrency       = errors.New("payment currency mismatch")
	ErrInvalidWebhook        = errors.New("invalid payment webhook")
)

type Order struct {
	ID              string    `json:"-"`
	OrderNo         string    `json:"orderNo"`
	PayToken        string    `json:"-"`
	Status          string    `json:"status"`
	Kind            string    `json:"kind"`
	CustomerName    string    `json:"customerName"`
	// Contact and delivery details are returned so the confirmation page can
	// show the customer what they just bought and where it is going, which is
	// the whole point of that page. The pay token in the URL is the access
	// secret, so anyone the link is forwarded to sees these too. Nothing here
	// is a credential, and none of it is exposed anywhere the token is absent.
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	Address         string    `json:"address"`
	DeliveryZone    string    `json:"deliveryZone"`
	SubtotalKobo    int64     `json:"subtotalKobo"`
	TotalKobo       int64     `json:"totalKobo"`
	DeliveryFeeKobo int64     `json:"deliveryFeeKobo"`
	PaymentProvider string    `json:"-"`
	PaymentRef      string    `json:"-"`
	Items           []Item    `json:"items"`
	CreatedAt       time.Time `json:"createdAt"`
}

type Item struct {
	Name      string `json:"name"`
	Size      string `json:"size,omitempty"`
	Quantity  int    `json:"quantity"`
	PriceKobo int64  `json:"-"`
	LineKobo  int64  `json:"lineKobo"`
}

type Attempt struct {
	OrderID          string
	Provider         string
	Reference        string
	AmountKobo       int64
	Currency         string
	Status           string
	AuthorizationURL string
	AccessCode       string
	ProviderPayload  []byte
}

type InitializeRequest struct {
	Email       string
	AmountKobo  int64
	Currency    string
	Reference   string
	CallbackURL string
	OrderNo     string
	Customer    string
}

type InitializeResult struct {
	AuthorizationURL string `json:"authorizationUrl"`
	AccessCode       string `json:"accessCode"`
	Reference        string `json:"reference"`
}

type Verification struct {
	Reference       string
	Status          string
	AmountKobo      int64
	Currency        string
	PaidAt          *time.Time
	GatewayResponse string
	// Channel is how the customer paid: card, bank_transfer, ussd. Carried
	// through for display only, never for any money decision.
	Channel         string
	ProviderPayload []byte
}

type Summary struct {
	Order   Order  `json:"order"`
	Payable bool   `json:"payable"`
	Message string `json:"message"`
}

type VerifyResult struct {
	Order     Order  `json:"order"`
	Reference string `json:"reference"`
	Status    string `json:"status"`
	Paid      bool   `json:"paid"`
	Message   string `json:"message"`
	// Channel drives the "Paystack · Card" line on the confirmation page.
	Channel string `json:"channel,omitempty"`
}
