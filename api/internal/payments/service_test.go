package payments

import (
	"context"
	"errors"
	"testing"
)

func TestInitializeReusesExistingInitializedAttempt(t *testing.T) {
	t.Parallel()

	repository := newPaymentRepositoryStub()
	repository.reusableAttempt = InitializeResult{
		AuthorizationURL: "https://checkout.paystack.com/reusable",
		AccessCode:       "access-code",
		Reference:        "OYE-1001-reusable",
	}
	provider := &paymentProviderStub{}
	service := NewService(repository, provider, "http://localhost:3000")

	result, err := service.Initialize(context.Background(), repository.order.PayToken)
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	if result.Reference != repository.reusableAttempt.Reference {
		t.Fatalf("Initialize() reference = %q, want reusable reference", result.Reference)
	}
	if repository.createAttemptCount != 0 {
		t.Fatalf("CreateAttempt() calls = %d, want 0", repository.createAttemptCount)
	}
	if provider.initializeCount != 0 {
		t.Fatalf("provider Initialize() calls = %d, want 0", provider.initializeCount)
	}
}

func TestInitializeFallsBackToReusableAttemptOnInitializedConflict(t *testing.T) {
	t.Parallel()

	repository := newPaymentRepositoryStub()
	repository.initializeError = ErrPaymentInProgress
	repository.reusableAfterInitializeConflict = InitializeResult{
		AuthorizationURL: "https://checkout.paystack.com/winner",
		AccessCode:       "winner-code",
		Reference:        "OYE-1001-winner",
	}
	provider := &paymentProviderStub{
		initializeResult: InitializeResult{
			AuthorizationURL: "https://checkout.paystack.com/loser",
			AccessCode:       "loser-code",
			Reference:        "generated-reference",
		},
	}
	service := NewService(repository, provider, "http://localhost:3000")

	result, err := service.Initialize(context.Background(), repository.order.PayToken)
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	if result.Reference != repository.reusableAfterInitializeConflict.Reference {
		t.Fatalf("Initialize() reference = %q, want winning reusable reference", result.Reference)
	}
	if repository.failedAttemptCount != 1 {
		t.Fatalf("MarkAttemptFailed() calls = %d, want 1", repository.failedAttemptCount)
	}
}

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
	order                                   Order
	reusableAttempt                         InitializeResult
	reusableAfterInitializeConflict         InitializeResult
	createAttemptCount                      int
	failedAttemptCount                      int
	confirmed                               bool
	confirmError                            error
	initializeError                         error
	recordedStatus                          string
	reusableInitializedAttemptLookupCounter int
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

func (repository *paymentRepositoryStub) ReusableInitializedAttempt(
	context.Context,
	string,
	string,
	int64,
	string,
) (InitializeResult, bool, error) {
	repository.reusableInitializedAttemptLookupCounter++
	if repository.reusableInitializedAttemptLookupCounter > 1 &&
		repository.reusableAfterInitializeConflict.Reference != "" {
		return repository.reusableAfterInitializeConflict, true, nil
	}
	if repository.reusableAttempt.Reference == "" {
		return InitializeResult{}, false, nil
	}
	return repository.reusableAttempt, true, nil
}

func (repository *paymentRepositoryStub) CreateAttempt(context.Context, Attempt) error {
	repository.createAttemptCount++
	return nil
}

func (repository *paymentRepositoryStub) MarkAttemptInitialized(
	context.Context,
	string,
	InitializeResult,
	[]byte,
) error {
	return repository.initializeError
}

func (repository *paymentRepositoryStub) MarkAttemptFailed(context.Context, string, []byte) error {
	repository.failedAttemptCount++
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
	initializeResult InitializeResult
	initializeCount  int
	verification     Verification
}

func (provider *paymentProviderStub) Name() string { return ProviderPaystack }

func (provider *paymentProviderStub) Initialize(
	_ context.Context,
	request InitializeRequest,
) (InitializeResult, []byte, error) {
	provider.initializeCount++
	if provider.initializeResult.Reference == "" {
		return InitializeResult{
			AuthorizationURL: "https://checkout.paystack.com/generated",
			AccessCode:       "generated-code",
			Reference:        request.Reference,
		}, nil, nil
	}

	result := provider.initializeResult
	if result.Reference == "generated-reference" {
		result.Reference = request.Reference
	}
	return result, nil, nil
}

func (provider *paymentProviderStub) Verify(context.Context, string) (Verification, error) {
	return provider.verification, nil
}
