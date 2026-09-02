package email

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// TestSendPaidOrderEmailsSendsOneConfirmation pins the single-email decision.
// This used to send a confirmation and a receipt a second apart, which reached
// the customer as near-duplicates; the client's approved design has one
// "Payment received" message carrying everything both used to.
func TestSendPaidOrderEmailsSendsOneConfirmation(t *testing.T) {
	t.Parallel()

	repository := &repositoryStub{reserved: true}
	provider := &providerStub{}
	service := NewService(repository, provider, "Oyecouture <orders@example.com>", "https://oyecouture.test")

	err := service.SendPaidOrderEmails(context.Background(), testOrder())
	if err != nil {
		t.Fatalf("SendPaidOrderEmails() error = %v", err)
	}

	if provider.sentCount != 1 {
		t.Fatalf("sent emails = %d, want 1", provider.sentCount)
	}
	if repository.sentCount != 1 {
		t.Fatalf("marked sent emails = %d, want 1", repository.sentCount)
	}
	if provider.messages[0].HTML == "" || provider.messages[0].Text == "" {
		t.Fatal("confirmation email did not include html and text bodies")
	}
	// Mail clients block remote images by default, so a hosted logo means most
	// recipients open an email carrying no brand at all. It is set as type.
	if strings.Contains(provider.messages[0].HTML, "<img") {
		t.Fatal("paid order emails should not rely on externally hosted logo images")
	}
	if !strings.Contains(provider.messages[0].HTML, "Oyecouture") {
		t.Fatal("paid order emails should include the Oyecouture brand wordmark")
	}
	// The design's "View your order" button. Without the reference the
	// confirmation page cannot verify and reports failure to someone who paid.
	if !strings.Contains(provider.messages[0].HTML, "/callback?reference=") {
		t.Fatal("confirmation email should link to the customer's own order page")
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

func TestFulfillmentWorkerDeliversQueuedTrackingEmail(t *testing.T) {
	t.Parallel()

	payload, err := json.Marshal(FulfillmentPayload{
		CustomerName: "Ada",
		OrderNo:      "OYE-1001",
		StageLabel:   "Dispatched",
		TrackingURL:  "https://oyecouture.test/track/secure-token",
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	repository := &repositoryStub{queuedEvent: &QueuedEvent{
		IdempotencyKey: "order:OYE-1001:fulfillment:5:dispatched",
		Type:           TypeOrderDispatched,
		RecipientEmail: "ada@example.com",
		Subject:        "Your order is on its way",
		Payload:        payload,
	}}
	provider := &providerStub{}
	service := NewService(repository, provider, "Oyecouture <orders@example.com>", "https://oyecouture.test")

	service.drainFulfillmentQueue(context.Background())

	if provider.sentCount != 1 || repository.sentCount != 1 {
		t.Fatalf("sent = %d, recorded = %d, want 1 each", provider.sentCount, repository.sentCount)
	}
	if !strings.Contains(provider.messages[0].HTML, "https://oyecouture.test/track/secure-token") {
		t.Fatal("dispatch email did not contain the secure tracking link")
	}
}

func testOrder() Order {
	return Order{
		ID:              "order-id",
		OrderNo:         "OYE-1001",
		PayToken:        "paytoken123",
		OrderDate:       "1 September 2026",
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
	reserved    bool
	sentCount   int
	queuedEvent *QueuedEvent
}

func (repository *repositoryStub) Reserve(context.Context, Event) (bool, error) {
	if !repository.reserved {
		return false, nil
	}
	return true, nil
}

func (repository *repositoryStub) Enqueue(context.Context, Event) error {
	return nil
}

func (repository *repositoryStub) ClaimNext(context.Context) (QueuedEvent, bool, error) {
	if repository.queuedEvent == nil {
		return QueuedEvent{}, false, nil
	}

	event := *repository.queuedEvent
	repository.queuedEvent = nil
	return event, true, nil
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
