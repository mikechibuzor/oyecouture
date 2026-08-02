package paystack

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/uzor-co/oyecouture/api/internal/payments"
)

const testSecret = "test-secret"

// confirmStub fails webhook confirmation with a chosen error.
type confirmStub struct {
	err   error
	calls int
}

func (stub *confirmStub) ConfirmWebhookReference(context.Context, string) error {
	stub.calls++
	return stub.err
}

// signedRequest builds a webhook request Paystack would consider authentic.
func signedRequest(body string) *http.Request {
	mac := hmac.New(sha512.New, []byte(testSecret))
	mac.Write([]byte(body))

	request := httptest.NewRequest(http.MethodPost, "/webhooks/paystack", strings.NewReader(body))
	request.Header.Set("x-paystack-signature", hex.EncodeToString(mac.Sum(nil)))

	return request
}

// TestWebhookRetryBehaviour pins which failures Paystack should redeliver.
// Getting this wrong is costly in both directions: acknowledging a transient
// failure drops a real payment, and returning 500 for something unresolvable
// makes Paystack retry it forever.
func TestWebhookRetryBehaviour(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		confirmErr error
		wantStatus int
	}{
		{
			name:       "success is acknowledged",
			confirmErr: nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "a reference we do not know is not retryable",
			confirmErr: payments.ErrReferenceNotFound,
			wantStatus: http.StatusOK,
		},
		{
			name:       "a reference belonging to another order is not retryable",
			confirmErr: payments.ErrReferenceMismatch,
			wantStatus: http.StatusOK,
		},
		{
			name:       "an amount mismatch is not retryable",
			confirmErr: payments.ErrPaymentAmountMismatch,
			wantStatus: http.StatusOK,
		},
		{
			name:       "a currency mismatch is not retryable",
			confirmErr: payments.ErrPaymentCurrency,
			wantStatus: http.StatusOK,
		},
		{
			name:       "an order that cannot be paid is not retryable",
			confirmErr: payments.ErrOrderNotPayable,
			wantStatus: http.StatusOK,
		},
		{
			name:       "wrapped terminal errors are still recognised",
			confirmErr: fmt.Errorf("confirm payment: %w", payments.ErrReferenceNotFound),
			wantStatus: http.StatusOK,
		},
		{
			name:       "paystack rejecting the verification is not retryable",
			confirmErr: fmt.Errorf("%w: Transaction reference not found", ErrVerifyRejected),
			wantStatus: http.StatusOK,
		},
		{
			// A 400 will answer the same way every time; retrying it forever is
			// the storm this classification exists to prevent.
			name:       "a 4xx from paystack is not retryable",
			confirmErr: &APIError{StatusCode: 400, Status: "400 Bad Request"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "a 404 from paystack is not retryable",
			confirmErr: fmt.Errorf("verify: %w", &APIError{StatusCode: 404, Status: "404 Not Found"}),
			wantStatus: http.StatusOK,
		},
		{
			// Their side struggling is exactly what retries are for.
			name:       "a 5xx from paystack is retried",
			confirmErr: &APIError{StatusCode: 502, Status: "502 Bad Gateway"},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "a paystack timeout is retried",
			confirmErr: &APIError{StatusCode: 504, Status: "504 Gateway Timeout"},
			wantStatus: http.StatusInternalServerError,
		},
		{
			// The important direction: anything unrecognised must be retried,
			// because acknowledging it would silently drop a real payment.
			name:       "an unknown failure is retried",
			confirmErr: errors.New("dial tcp: connection refused"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "the provider being unreachable is retried",
			confirmErr: payments.ErrPaymentUnavailable,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			stub := &confirmStub{err: testCase.confirmErr}
			recorder := httptest.NewRecorder()

			NewWebhookHandler(stub, testSecret).Handle(
				recorder,
				signedRequest(`{"event":"charge.success","data":{"reference":"OYE-1001-abc"}}`),
			)

			if recorder.Code != testCase.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, testCase.wantStatus)
			}
			if stub.calls != 1 {
				t.Fatalf("confirm called %d times, want 1", stub.calls)
			}
		})
	}
}

// TestWebhookRejectsUnauthentic keeps the signature gate ahead of any work.
func TestWebhookRejectsUnauthentic(t *testing.T) {
	t.Parallel()

	body := `{"event":"charge.success","data":{"reference":"OYE-1001-abc"}}`

	testCases := []struct {
		name      string
		signature string
	}{
		{name: "missing signature", signature: ""},
		{name: "wrong signature", signature: "deadbeef"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			stub := &confirmStub{}
			request := httptest.NewRequest(http.MethodPost, "/webhooks/paystack", strings.NewReader(body))
			if testCase.signature != "" {
				request.Header.Set("x-paystack-signature", testCase.signature)
			}

			recorder := httptest.NewRecorder()
			NewWebhookHandler(stub, testSecret).Handle(recorder, request)

			if recorder.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
			}
			if stub.calls != 0 {
				t.Fatal("an unauthentic webhook reached the payment service")
			}
		})
	}
}
