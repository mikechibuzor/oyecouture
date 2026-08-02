package admin

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestLoginLimiterChargesFailuresOnly verifies the throttle stops a brute force
// run without penalising an owner who signs in repeatedly through the day.
func TestLoginLimiterChargesFailuresOnly(t *testing.T) {
	t.Parallel()

	limiter := newLoginLimiter()

	// Successes never call recordFailure, so the budget stays untouched.
	for attempt := range loginAttemptsAllowed * 2 {
		if !limiter.allow("198.51.100.7") {
			t.Fatalf("successful login %d was throttled", attempt+1)
		}
	}

	for attempt := range loginAttemptsAllowed {
		if !limiter.allow("198.51.100.7") {
			t.Fatalf("failed login %d was throttled before the limit", attempt+1)
		}
		limiter.recordFailure("198.51.100.7")
	}

	if limiter.allow("198.51.100.7") {
		t.Fatal("attempt past the limit was allowed")
	}

	// The window is per address, so one attacker cannot lock out the owner.
	if !limiter.allow("203.0.113.4") {
		t.Fatal("a different address was caught by another address's lockout")
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

			request := httptest.NewRequest(http.MethodPost, "/admin/login", nil)
			request.RemoteAddr = testCase.remote
			if testCase.forwarded != "" {
				request.Header.Set("X-Forwarded-For", testCase.forwarded)
			}

			if got := clientAddress(request); got != testCase.want {
				t.Fatalf("clientAddress() = %q, want %q", got, testCase.want)
			}
		})
	}
}
