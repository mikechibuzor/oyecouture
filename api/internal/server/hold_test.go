package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// holdGateCall records what one pass through the gate did, so each case can assert on both
// whether the route ran and what the visitor was told.
type holdGateCall struct {
	reachedRoute bool
	statusCode   int
}

// exercise runs one request through the gate with the given session and hold answers.
func exercise(
	t *testing.T,
	hasSession bool,
	isHeld bool,
	holdError error,
) holdGateCall {
	t.Helper()

	call := holdGateCall{}
	gate := holdGate(
		func(*http.Request) bool { return hasSession },
		func(context.Context) (bool, error) { return isHeld, holdError },
	)

	handler := gate(func(responseWriter http.ResponseWriter, _ *http.Request) {
		call.reachedRoute = true
		responseWriter.WriteHeader(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	handler(recorder, httptest.NewRequest(http.MethodGet, "/products", nil))
	call.statusCode = recorder.Code

	return call
}

func TestHoldGateServesTheShopWhenItIsNotHeld(t *testing.T) {
	call := exercise(t, false, false, nil)

	if !call.reachedRoute {
		t.Fatal("expected an open shop to serve the route")
	}
	if call.statusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", call.statusCode)
	}
}

func TestHoldGateRefusesVisitorsWhileHeld(t *testing.T) {
	call := exercise(t, false, true, nil)

	if call.reachedRoute {
		t.Fatal("a held shop must not run the route; the catalogue would still be readable")
	}
	if call.statusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", call.statusCode)
	}
}

func TestHoldGateLetsASignedInAdminThrough(t *testing.T) {
	call := exercise(t, true, true, nil)

	if !call.reachedRoute {
		t.Fatal("expected an admin to see the real shop while it is held")
	}
	if call.statusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", call.statusCode)
	}
}

// A database hiccup must not close a working shop: the hold is turned on deliberately.
func TestHoldGateStaysOpenWhenTheHoldCannotBeRead(t *testing.T) {
	call := exercise(t, false, false, errors.New("database unavailable"))

	if !call.reachedRoute {
		t.Fatal("expected a failed hold check to leave the shop open")
	}
	if call.statusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", call.statusCode)
	}
}

// If the flag reads true the shop closes even when the read also reported a problem, so a
// partial failure cannot quietly reopen a shop the owner has closed.
func TestHoldGateRefusesWhenHeldIsTrueDespiteAnError(t *testing.T) {
	call := exercise(t, false, true, errors.New("database unavailable"))

	if call.reachedRoute {
		t.Fatal("expected a held shop to stay closed")
	}
	if call.statusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", call.statusCode)
	}
}
