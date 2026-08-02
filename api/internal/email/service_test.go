package email

import (
	"context"
	"testing"
)

func TestSendPaidOrderEmailsSendsConfirmationAndReceipt(t *testing.T) {
	t.Parallel()

	repository := &repositoryStub{reserved: true}
	provider := &providerStub{}
	service := NewService(repository, provider, "Oyecouture <orders@example.com>", "https://oyecouture.test")

	err := service.SendPaidOrderEmails(context.Background(), testOrder())
	if err != nil {
		t.Fatalf("SendPaidOrderEmails() error = %v", err)
	}

	if provider.sentCount != 2 {
		t.Fatalf("sent emails = %d, want 2", provider.sentCount)
	}
	if repository.sentCount != 2 {
		t.Fatalf("marked sent emails = %d, want 2", repository.sentCount)
	}
	if provider.messages[0].IdempotencyKey == provider.messages[1].IdempotencyKey {
		t.Fatal("confirmation and receipt used the same idempotency key")
	}
	if provider.messages[0].HTML == "" || provider.messages[0].Text == "" {
		t.Fatal("confirmation email did not include html and text bodies")
	}
}

func TestSendPaidOrderEmailsSkipsAlreadyReservedEmails(t *testing.T) {
	t.Parallel()

	repository := &repositoryStub{reserved: false}
	provider := &providerStub{}
	service := NewService(repository, provider, "Oyecouture <orders@example.com>", "https://oyecouture.test")

	err := service.SendPaidOrderEmails(context.Background(), testOrder())
	if err != nil {
		t.Fatalf("SendPaidOrderEmails() error = %v", err)
	}

	if provider.sentCount != 0 {
		t.Fatalf("sent emails = %d, want 0", provider.sentCount)
	}
}

func testOrder() Order {
	return Order{
		ID:              "order-id",
		OrderNo:         "OYE-1001",
		CustomerName:    "Ada",
		Email:           "ada@example.com",
		Phone:           "07000000000",
		Address:         "1 Lagos Street",
		DeliveryZone:    "Lagos",
		SubtotalKobo:    8600000,
		DeliveryFeeKobo: 500000,
		TotalKobo:       9100000,
		Items: []Item{{
			Name:      "The Oxford Overshirt",
			Size:      "M",
			Quantity:  1,
			PriceKobo: 8600000,
			LineKobo:  8600000,
		}},
		Payment: Payment{
			Provider:  "paystack",
			Reference: "OYE-1001-reference",
			PaidAt:    "2 Aug 2026, 16:18 WAT",
		},
	}
}

type repositoryStub struct {
	reserved  bool
	sentCount int
}

func (repository *repositoryStub) Reserve(context.Context, Event) (bool, error) {
	if !repository.reserved {
		return false, nil
	}
	return true, nil
}

func (repository *repositoryStub) MarkSent(context.Context, string, SendResult, []byte) error {
	repository.sentCount++
	return nil
}

func (repository *repositoryStub) MarkFailed(context.Context, string, string, []byte) error {
	return nil
}

type providerStub struct {
	sentCount int
	messages  []Message
}

func (provider *providerStub) Name() string { return ProviderSendByte }

func (provider *providerStub) Send(
	_ context.Context,
	message Message,
) (SendResult, []byte, error) {
	provider.sentCount++
	provider.messages = append(provider.messages, message)
	return SendResult{MessageID: "message-id"}, []byte(`{"id":"message-id"}`), nil
}
