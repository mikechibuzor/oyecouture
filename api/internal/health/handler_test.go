package health

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type databasePingerFunc func(context.Context) error

func (pinger databasePingerFunc) Ping(ctx context.Context) error {
	return pinger(ctx)
}

func TestHandlerReturnsHealthyServices(t *testing.T) {
	t.Parallel()

	handler := NewHandler(databasePingerFunc(func(context.Context) error {
		return nil
	}))
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	responseRecorder := httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}

	expectedBody := "{\"status\":\"ok\",\"services\":{\"api\":\"ok\",\"database\":\"ok\"}}\n"
	if responseRecorder.Body.String() != expectedBody {
		t.Fatalf("expected body %q, got %q", expectedBody, responseRecorder.Body.String())
	}
}

func TestHandlerReturnsDegradedWhenDatabaseIsUnavailable(t *testing.T) {
	t.Parallel()

	handler := NewHandler(databasePingerFunc(func(context.Context) error {
		return errors.New("database unavailable")
	}))
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	responseRecorder := httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, responseRecorder.Code)
	}

	expectedBody := "{\"status\":\"degraded\",\"services\":{\"api\":\"ok\",\"database\":\"unavailable\"}}\n"
	if responseRecorder.Body.String() != expectedBody {
		t.Fatalf("expected body %q, got %q", expectedBody, responseRecorder.Body.String())
	}
}
