package admin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// storeStub serves one known account so the handler can be exercised without a
// database. Limiter mechanics live in internal/ratelimit and are tested there;
// what matters here is how the sign-in endpoint behaves.
type storeStub struct {
	email        string
	passwordHash string
}

func (stub *storeStub) UserByEmail(_ context.Context, email string) (User, string, error) {
	if email != stub.email {
		return User{}, "", errors.New("not found")
	}

	return User{ID: "11111111-1111-1111-1111-111111111111", Email: stub.email, Role: "owner"},
		stub.passwordHash, nil
}

func (stub *storeStub) CreateSession(context.Context, string, string, time.Time) error { return nil }
func (stub *storeStub) SessionByToken(context.Context, string) (Session, error) {
	return Session{}, errors.New("not found")
}
func (stub *storeStub) DeleteSession(context.Context, string) error { return nil }
func (stub *storeStub) DeleteExpiredSessions(context.Context) error { return nil }

func newTestHandler(t *testing.T) *Handler {
	t.Helper()

	hash, err := HashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	return NewHandler(NewAuthService(&storeStub{
		email:        "owner@example.com",
		passwordHash: hash,
	}), false)
}

func postLogin(handler *Handler, email string, password string) *httptest.ResponseRecorder {
	body, _ := json.Marshal(loginRequest{Email: email, Password: password})
	request := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(string(body)))
	request.RemoteAddr = "198.51.100.7:54321"

	recorder := httptest.NewRecorder()
	handler.LogIn(recorder, request)

	return recorder
}

// TestSignInThrottlesBruteForce proves the endpoint stops repeated guessing,
// and, importantly, that the owner cannot lock
// themselves out by signing in normally.
func TestSignInThrottlesBruteForce(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t)

	// Successful sign-ins never charge the budget.
	for attempt := range loginAttemptsAllowed * 2 {
		if code := postLogin(handler, "owner@example.com", "correct-horse-battery-staple").Code; code != http.StatusOK {
			t.Fatalf("successful sign-in %d returned %d, want 200", attempt+1, code)
		}
	}

	for attempt := range loginAttemptsAllowed {
		if code := postLogin(handler, "owner@example.com", "wrong").Code; code != http.StatusUnauthorized {
			t.Fatalf("failed sign-in %d returned %d, want 401", attempt+1, code)
		}
	}

	if code := postLogin(handler, "owner@example.com", "wrong").Code; code != http.StatusTooManyRequests {
		t.Fatalf("attempt past the limit returned %d, want 429", code)
	}

	// Even the right password is refused while locked out.
	if code := postLogin(handler, "owner@example.com", "correct-horse-battery-staple").Code; code != http.StatusTooManyRequests {
		t.Fatalf("correct password during lockout returned %d, want 429", code)
	}
}

// TestSignInDoesNotRevealWhichAccountsExist keeps an unknown email and a wrong
// password indistinguishable to a caller.
func TestSignInDoesNotRevealWhichAccountsExist(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t)

	unknown := postLogin(handler, "nobody@example.com", "whatever-long-enough")
	wrong := postLogin(handler, "owner@example.com", "wrong-password")

	if unknown.Code != wrong.Code {
		t.Fatalf("status differed: unknown email %d, wrong password %d", unknown.Code, wrong.Code)
	}
	if unknown.Body.String() != wrong.Body.String() {
		t.Fatalf("body differed:\n unknown: %s\n wrong:   %s", unknown.Body.String(), wrong.Body.String())
	}
}
