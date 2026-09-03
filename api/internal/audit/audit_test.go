package audit

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/uzor-co/oyecouture/api/internal/admin"
)

type storeStub struct {
	entries []Entry
	actors  []string
	err     error
}

func (stub *storeStub) Record(_ context.Context, entry Entry, actorAdminID string) error {
	stub.entries = append(stub.entries, entry)
	stub.actors = append(stub.actors, actorAdminID)

	return stub.err
}

func (stub *storeStub) List(context.Context, int, int) (Page, error) {
	return Page{Entries: stub.entries}, nil
}

// withSession runs the recorder as it runs in production: inside RequireSession,
// so the session is already in the context.
func withSession(handler http.Handler, session admin.Session) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(responseWriter, admin.WithSession(request, session))
	})
}

func testSession() admin.Session {
	return admin.Session{
		User: admin.User{
			ID:    "11111111-1111-1111-1111-111111111111",
			Email: "owner@example.com",
			Role:  admin.RoleOwner,
		},
	}
}

// TestRecordsEveryMutation is the property the whole feature rests on. A log
// with gaps is worse than none, because anything absent from it reads as
// something that did not happen.
func TestRecordsEveryMutation(t *testing.T) {
	t.Parallel()

	for _, method := range []string{
		http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete,
	} {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			store := &storeStub{}
			handler := withSession(
				NewService(store).Record(http.HandlerFunc(
					func(responseWriter http.ResponseWriter, _ *http.Request) {
						responseWriter.WriteHeader(http.StatusOK)
					},
				)),
				testSession(),
			)

			request := httptest.NewRequest(method, "/admin/products/monarch", nil)
			handler.ServeHTTP(httptest.NewRecorder(), request)

			if len(store.entries) != 1 {
				t.Fatalf("entries = %d, want 1", len(store.entries))
			}
			entry := store.entries[0]
			if entry.Method != method || entry.Path != "/admin/products/monarch" {
				t.Fatalf("recorded %s %s, want %s /admin/products/monarch",
					entry.Method, entry.Path, method)
			}
			if entry.ActorEmail != "owner@example.com" || entry.ActorRole != admin.RoleOwner {
				t.Fatalf("actor = %s/%s, want owner@example.com/owner",
					entry.ActorEmail, entry.ActorRole)
			}
		})
	}
}

// TestReadsAreNotRecorded keeps the log to changes. Logging every page view
// would bury the actions worth noticing.
func TestReadsAreNotRecorded(t *testing.T) {
	t.Parallel()

	store := &storeStub{}
	handler := withSession(
		NewService(store).Record(http.HandlerFunc(
			func(responseWriter http.ResponseWriter, _ *http.Request) {
				responseWriter.WriteHeader(http.StatusOK)
			},
		)),
		testSession(),
	)

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/admin/orders", nil))

	if len(store.entries) != 0 {
		t.Fatalf("entries = %d, want none for a read", len(store.entries))
	}
}

// TestRefusedAttemptsAreRecorded covers the more interesting half of an audit
// trail: someone trying to reach something they cannot.
func TestRefusedAttemptsAreRecorded(t *testing.T) {
	t.Parallel()

	store := &storeStub{}
	handler := withSession(
		NewService(store).Record(http.HandlerFunc(
			func(responseWriter http.ResponseWriter, _ *http.Request) {
				responseWriter.WriteHeader(http.StatusForbidden)
			},
		)),
		testSession(),
	)

	handler.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest(http.MethodDelete, "/admin/team/22222222-2222-2222-2222-222222222222", nil),
	)

	if len(store.entries) != 1 {
		t.Fatalf("entries = %d, want the refused attempt recorded", len(store.entries))
	}
	if store.entries[0].StatusCode != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", store.entries[0].StatusCode)
	}
}

// TestUnauthenticatedRequestsAreNotRecorded avoids entries with no actor. They
// never reached anything to change, and sign-in attempts are throttled and
// logged separately.
func TestUnauthenticatedRequestsAreNotRecorded(t *testing.T) {
	t.Parallel()

	store := &storeStub{}
	handler := NewService(store).Record(http.HandlerFunc(
		func(responseWriter http.ResponseWriter, _ *http.Request) {
			responseWriter.WriteHeader(http.StatusUnauthorized)
		},
	))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/admin/team", nil))

	if len(store.entries) != 0 {
		t.Fatalf("entries = %d, want none without an actor", len(store.entries))
	}
}

// TestAWriteFailureDoesNotBreakTheRequest keeps a logging problem from
// undoing work that already succeeded and was already answered.
func TestAWriteFailureDoesNotBreakTheRequest(t *testing.T) {
	t.Parallel()

	store := &storeStub{err: errors.New("database unavailable")}
	handler := withSession(
		NewService(store).Record(http.HandlerFunc(
			func(responseWriter http.ResponseWriter, _ *http.Request) {
				responseWriter.WriteHeader(http.StatusCreated)
			},
		)),
		testSession(),
	)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/admin/products", nil))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want the handler's 201 despite the audit failure", recorder.Code)
	}
}

// TestImplicitOkIsRecorded covers a handler that writes a body without calling
// WriteHeader, which would otherwise be logged with a zero status.
func TestImplicitOkIsRecorded(t *testing.T) {
	t.Parallel()

	store := &storeStub{}
	handler := withSession(
		NewService(store).Record(http.HandlerFunc(
			func(responseWriter http.ResponseWriter, _ *http.Request) {
				_, _ = responseWriter.Write([]byte(`{"ok":true}`))
			},
		)),
		testSession(),
	)

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPatch, "/admin/settings", nil))

	if len(store.entries) != 1 {
		t.Fatalf("entries = %d, want 1", len(store.entries))
	}
	if store.entries[0].StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", store.entries[0].StatusCode)
	}
}
