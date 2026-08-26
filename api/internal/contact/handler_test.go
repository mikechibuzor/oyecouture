package contact

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// submitterStub counts how many submissions actually reached the service, which
// is what decides whether an email would have been sent.
type submitterStub struct {
	err   error
	calls int
}

func (stub *submitterStub) Submit(context.Context, Submission) error {
	stub.calls++
	return stub.err
}

func postContact(handler *Handler, address string, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(body))
	request.RemoteAddr = address

	recorder := httptest.NewRecorder()
	handler.Submit(recorder, request)

	return recorder
}

func validBody() string {
	body, _ := json.Marshal(Submission{
		Name:    "Ada Obi",
		Email:   "ada@example.com",
		Topic:   "Bespoke",
		Message: "I would like to book a fitting.",
	})

	return string(body)
}

// TestSubmitThrottlesFlooding is the point of the limiter: the form is public
// and every accepted message sends an email, so unbounded submissions would
// bury the owner's inbox, burn the email quota, and put the sending domain's
// reputation and therefore order emails at risk.
func TestSubmitThrottlesFlooding(t *testing.T) {
	t.Parallel()

	stub := &submitterStub{}
	handler := NewHandler(stub)

	for attempt := range submissionsPerHour {
		if code := postContact(handler, "198.51.100.7:5000", validBody()).Code; code != http.StatusCreated {
			t.Fatalf("submission %d returned %d, want 201", attempt+1, code)
		}
	}

	recorder := postContact(handler, "198.51.100.7:5000", validBody())
	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("submission past the limit returned %d, want 429", recorder.Code)
	}
	if stub.calls != submissionsPerHour {
		t.Fatalf("service was called %d times, want %d; a throttled message still sent an email",
			stub.calls, submissionsPerHour)
	}

	// The window is per address; one flooder must not silence real enquiries.
	if code := postContact(handler, "203.0.113.4:5000", validBody()).Code; code != http.StatusCreated {
		t.Fatalf("a different address returned %d, want 201", code)
	}
}

// TestSubmitDoesNotChargeRejectedMessages keeps a mistyped form from eating a
// real person's allowance. Only messages that actually send are counted.
func TestSubmitDoesNotChargeRejectedMessages(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&submitterStub{})

	for range submissionsPerHour * 3 {
		if code := postContact(handler, "198.51.100.9:5000", "{not json").Code; code != http.StatusBadRequest {
			t.Fatalf("malformed body returned %d, want 400", code)
		}
	}

	if code := postContact(handler, "198.51.100.9:5000", validBody()).Code; code != http.StatusCreated {
		t.Fatalf("a valid message after malformed ones returned %d, want 201", code)
	}
}

// TestSubmitThrottlesPerAddressNotPerPort covers the trap that the ephemeral
// source port would give every request its own bucket.
func TestSubmitThrottlesPerAddressNotPerPort(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&submitterStub{})

	for attempt := range submissionsPerHour {
		// A different source port each time, as a real client would use.
		address := "198.51.100.11:" + string(rune('0'+attempt)) + "0000"
		if code := postContact(handler, address, validBody()).Code; code != http.StatusCreated {
			t.Fatalf("submission %d returned %d, want 201", attempt+1, code)
		}
	}

	if code := postContact(handler, "198.51.100.11:59999", validBody()).Code; code != http.StatusTooManyRequests {
		t.Fatalf("changing source port bypassed the limit: got %d, want 429", code)
	}
}
