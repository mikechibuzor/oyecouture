package contact

import (
	"context"
	"errors"
	"testing"

	"github.com/uzor-co/oyecouture/api/internal/email"
)

type stubProvider struct {
	message email.Message
	err     error
}

func (provider *stubProvider) Name() string { return "stub" }

func (provider *stubProvider) Send(
	_ context.Context,
	message email.Message,
) (email.SendResult, []byte, error) {
	provider.message = message
	return email.SendResult{MessageID: "msg_123"}, nil, provider.err
}

// staticRecipient is the old fixed-address behaviour, for the cases that do not care where
// the enquiry goes.
func staticRecipient(address string) func(context.Context) (string, error) {
	return func(context.Context) (string, error) { return address, nil }
}

func TestServiceSubmitSendsContactEmail(t *testing.T) {
	provider := &stubProvider{}
	service := NewService(provider, "Oyecouture <orders@example.com>", staticRecipient("owner@example.com"))

	err := service.Submit(context.Background(), Submission{
		Name:    " Chibuzor ",
		Email:   " chibuzor@example.com ",
		Topic:   "Feedback",
		Rating:  5,
		Message: "Lovely cut.\nMore linen please.",
	})
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}

	if provider.message.To != "owner@example.com" {
		t.Fatalf("message.To = %q", provider.message.To)
	}
	if provider.message.Subject != "Oyecouture contact: Feedback from Chibuzor" {
		t.Fatalf("message.Subject = %q", provider.message.Subject)
	}
	if provider.message.IdempotencyKey == "" {
		t.Fatal("expected idempotency key")
	}
	if provider.message.Text == "" || provider.message.HTML == "" {
		t.Fatal("expected text and html bodies")
	}
}

func TestServiceSubmitRejectsInvalidSubmission(t *testing.T) {
	service := NewService(&stubProvider{}, "from@example.com", staticRecipient("owner@example.com"))

	err := service.Submit(context.Background(), Submission{
		Name:    "Chibuzor",
		Email:   "not-an-email",
		Message: "Hello",
	})
	if !errors.Is(err, ErrInvalidSubmission) {
		t.Fatalf("Submit() error = %v, want %v", err, ErrInvalidSubmission)
	}
}

func TestServiceSubmitRequiresSenderConfiguration(t *testing.T) {
	service := NewService(nil, "", staticRecipient("owner@example.com"))

	err := service.Submit(context.Background(), Submission{
		Name:    "Chibuzor",
		Email:   "chibuzor@example.com",
		Message: "Hello",
	})
	if !errors.Is(err, ErrSenderUnavailable) {
		t.Fatalf("Submit() error = %v, want %v", err, ErrSenderUnavailable)
	}
}
