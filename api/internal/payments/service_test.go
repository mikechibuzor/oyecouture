package payments

import (
	"context"
	"errors"
	"testing"
)

func TestVerifyConfirmsSuccessfulPayment(t *testing.T) {
	t.Parallel()

	repository := newPaymentRepositoryStub()
	provider := &paymentProviderStub{
		verification: Verification{
			Reference:  "OYE-1001-test",
			Status:     StatusSuccess,
			AmountKobo: repository.order.TotalKobo,
			Currency:   CurrencyNGN,
		},
	}
	service := NewService(repository, provider, "http://localhost:3000")

	result, err := service.Verify(context.Background(), repository.order.PayToken, "OYE-1001-test")
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}

	if !result.Paid || result.Order.Status != OrderStatusPaid {
		t.Fatalf("Verify() paid = %v, status = %q", result.Paid, result.Order.Status)
	}
	if !repository.confirmed {
		t.Fatal("Verify() did not confirm the successful payment")
	}
}

func TestVerifyRecordsIncompletePaymentWithoutConfirming(t *testing.T) {
	t.Parallel()

	repository := newPaymentRepositoryStub()
	provider := &paymentProviderStub{
		verification: Verification{
			Reference:  "OYE-1001-test",
			Status:     "pending",
			AmountKobo: repository.order.TotalKobo,
			Currency:   CurrencyNGN,
		},
	}
	service := NewService(repository, provider, "http://localhost:3000")

	result, err := service.Verify(context.Background(), repository.order.PayToken, "OYE-1001-test")
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}

	if result.Paid {
		t.Fatal("Verify() marked an incomplete payment as paid")
	}
	if repository.confirmed {
		t.Fatal("Verify() confirmed an incomplete payment")
	}
	if repository.recordedStatus != "pending" {
		t.Fatalf("recorded status = %q, want pending", repository.recordedStatus)
	}
}

func TestVerifyRejectsReferenceMismatch(t *testing.T) {
	t.Parallel()

	repository := newPaymentRepositoryStub()
	provider := &paymentProviderStub{
		verification: Verification{
			Reference:  "different-reference",
			Status:     StatusSuccess,
			AmountKobo: repository.order.TotalKobo,
			Currency:   CurrencyNGN,
		},
	}
	service := NewService(repository, provider, "http://localhost:3000")

	_, err := service.Verify(context.Background(), repository.order.PayToken, "OYE-1001-test")
	if !errors.Is(err, ErrReferenceMismatch) {
		t.Fatalf("Verify() error = %v, want ErrReferenceMismatch", err)
	}
	if repository.confirmed {
		t.Fatal("Verify() confirmed a mismatched reference")
	}
}

func TestVerifyBindsReferenceToPayTokenOrder(t *testing.T) {
	t.Parallel()

	repository := newPaymentRepositoryStub()
	repository.confirmError = ErrReferenceMismatch
	provider := &paymentProviderStub{
		verification: Verification{
			Reference:  "OYE-9999-test",
			Status:     StatusSuccess,
			AmountKobo: repository.order.TotalKobo,
			Currency:   CurrencyNGN,
		},
	}
	service := NewService(repository, provider, "http://localhost:3000")

	_, err := service.Verify(context.Background(), repository.order.PayToken, "OYE-9999-test")
	if !errors.Is(err, ErrReferenceMismatch) {
		t.Fatalf("Verify() error = %v, want ErrReferenceMismatch", err)
	}
	if !repository.confirmed {
		t.Fatal("Verify() did not ask the repository to enforce the order binding")
	}
}

type paymentRepositoryStub struct {
	order          Order
	confirmed      bool
	confirmError   error
	recordedStatus string
}

func newPaymentRepositoryStub() *paymentRepositoryStub {
	return &paymentRepositoryStub{
		order: Order{
			ID:           "order-id",
			OrderNo:      "OYE-1001",
			PayToken:     "pay-token",
			Status:       OrderStatusAwaitingPayment,
			CustomerName: "Ada",
			Email:        "ada@example.com",
			TotalKobo:    8300000,
		},
	}
}

func (repository *paymentRepositoryStub) OrderByPayToken(context.Context, string) (Order, error) {
	return repository.order, nil
}

func (repository *paymentRepositoryStub) CreateAttempt(context.Context, Attempt) error { return nil }

func (repository *paymentRepositoryStub) MarkAttemptInitialized(
	context.Context,
	string,
	InitializeResult,
	[]byte,
) error {
	return nil
}

func (repository *paymentRepositoryStub) MarkAttemptFailed(context.Context, string, []byte) error {
	return nil
}

func (repository *paymentRepositoryStub) RecordVerification(
	_ context.Context,
	_ string,
	verification Verification,
) error {
	repository.recordedStatus = verification.Status
	return nil
}

func (repository *paymentRepositoryStub) ConfirmSuccessfulPayment(
	context.Context,
	string,
	Verification,
) (Order, error) {
	repository.confirmed = true
	if repository.confirmError != nil {
		return Order{}, repository.confirmError
	}
	repository.order.Status = OrderStatusPaid
	return repository.order, nil
}

type paymentProviderStub struct {
	verification Verification
}

func (provider *paymentProviderStub) Name() string { return ProviderPaystack }

func (provider *paymentProviderStub) Initialize(
	context.Context,
	InitializeRequest,
) (InitializeResult, []byte, error) {
	return InitializeResult{}, nil, nil
}

func (provider *paymentProviderStub) Verify(context.Context, string) (Verification, error) {
	return provider.verification, nil
}
