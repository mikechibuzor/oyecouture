package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestAllowAndRecordAreSeparate covers the split both callers depend on:
// sign-in checks on every attempt but charges only failures, while the contact
// form charges every submission.
func TestAllowAndRecordAreSeparate(t *testing.T) {
	t.Parallel()

	limiter := New(5*time.Minute, 3)

	// Checking alone never consumes the budget.
	for attempt := range 10 {
		if !limiter.Allow("198.51.100.7") {
			t.Fatalf("check %d was throttled without any Record", attempt+1)
		}
	}

	for attempt := range 3 {
		if !limiter.Allow("198.51.100.7") {
			t.Fatalf("attempt %d was throttled before the limit", attempt+1)
		}
		limiter.Record("198.51.100.7")
	}

	if limiter.Allow("198.51.100.7") {
		t.Fatal("attempt past the limit was allowed")
	}

	// The window is per address, so one abuser cannot lock out everyone else.
	if !limiter.Allow("203.0.113.4") {
		t.Fatal("a different address was caught by another address's lockout")
	}
}

// TestAttemptsAgeOut proves the window slides rather than banning permanently.
func TestAttemptsAgeOut(t *testing.T) {
	t.Parallel()

	limiter := New(50*time.Millisecond, 1)
	limiter.Record("198.51.100.7")

	if limiter.Allow("198.51.100.7") {
		t.Fatal("expected the address to be throttled immediately after recording")
	}

	time.Sleep(80 * time.Millisecond)

	if !limiter.Allow("198.51.100.7") {
		t.Fatal("expected the attempt to age out of the window")
	}
}

// TestPruneDoesNotLeakEntries keeps an in-memory limiter from growing an entry
// per address forever on a public endpoint.
func TestPruneDoesNotLeakEntries(t *testing.T) {
	t.Parallel()

	limiter := New(20*time.Millisecond, 5)
	for _, address := range []string{"198.51.100.1", "198.51.100.2", "198.51.100.3"} {
		limiter.Record(address)
	}

	time.Sleep(50 * time.Millisecond)
	for _, address := range []string{"198.51.100.1", "198.51.100.2", "198.51.100.3"} {
		limiter.Allow(address)
	}

	limiter.mutex.Lock()
	remaining := len(limiter.attempts)
	limiter.mutex.Unlock()

	if remaining != 0 {
		t.Fatalf("limiter retained %d expired entries, want 0", remaining)
	}
}

// TestClientAddressIgnoresPort covers the bug where the ephemeral source port
// gave every attempt its own bucket and the limiter never fired.
func TestClientAddressIgnoresPort(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		remote    string
		forwarded string
		want      string
	}{
		{name: "ipv4 strips port", remote: "198.51.100.7:54321", want: "198.51.100.7"},
		{name: "ipv6 strips port", remote: "[2001:db8::1]:54321", want: "2001:db8::1"},
		{name: "unparseable falls back", remote: "198.51.100.7", want: "198.51.100.7"},
		{
			name:      "proxy chain keeps the original client",
			remote:    "10.0.0.2:443",
			forwarded: "198.51.100.7, 10.0.0.1",
			want:      "198.51.100.7",
		},
		{
			name:      "single proxy header is trimmed",
			remote:    "10.0.0.2:443",
			forwarded: " 198.51.100.7 ",
			want:      "198.51.100.7",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequest(http.MethodPost, "/", nil)
			request.RemoteAddr = testCase.remote
			if testCase.forwarded != "" {
				request.Header.Set("X-Forwarded-For", testCase.forwarded)
			}

			if got := ClientAddress(request); got != testCase.want {
				t.Fatalf("ClientAddress() = %q, want %q", got, testCase.want)
			}
		})
	}
}
