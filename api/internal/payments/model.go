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
	ErrPaymentUnavailable    = errors.New("payment provider is unavailable")
	ErrReferenceNotFound     = errors.New("payment reference not found")
	ErrReferenceMismatch     = errors.New("payment reference mismatch")
	ErrPaymentAmountMismatch = errors.New("payment amount mismatch")
	ErrPaymentCurrency       = errors.New("payment currency mismatch")
	ErrInvalidWebhook        = errors.New("invalid payment webhook")
)

type Order struct {
	ID           string    `json:"-"`
	OrderNo      string    `json:"orderNo"`
	PayToken     string    `json:"-"`
	Status       string    `json:"status"`
	Kind         string    `json:"kind"`
	CustomerName string    `json:"customerName"`
	Email        string    `json:"-"`
	Phone        string    `json:"-"`
	TotalKobo    int64     `json:"totalKobo"`
	Items        []Item    `json:"items"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Item struct {
	Name     string `json:"name"`
	Size     string `json:"size,omitempty"`
	Quantity int    `json:"quantity"`
	LineKobo int64  `json:"lineKobo"`
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
}
