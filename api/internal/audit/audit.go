// Package audit records every change made through the admin, so the shop owner
// can verify what was done to their store rather than take anyone's word for
// it, the developer included.
package audit

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/uzor-co/oyecouture/api/internal/admin"
)

// Entry is one recorded admin action.
type Entry struct {
	ID         string     `json:"id"`
	ActorEmail string     `json:"actorEmail"`
	ActorRole  admin.Role `json:"actorRole"`
	Method     string     `json:"method"`
	Path       string     `json:"path"`
	StatusCode int        `json:"statusCode"`
	OccurredAt time.Time  `json:"occurredAt"`
}

type Page struct {
	Entries []Entry `json:"entries"`
	Total   int     `json:"total"`
}

// Store persists and reads the log. There is deliberately no update or delete:
// the table refuses both at the database level, and offering the methods here
// would suggest otherwise.
type Store interface {
	Record(context.Context, Entry, string) error
	List(context.Context, int, int) (Page, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (service *Service) List(
	requestContext context.Context,
	limit int,
	offset int,
) (Page, error) {
	if limit <= 0 || limit > maximumPageSize {
		limit = defaultPageSize
	}
	if offset < 0 {
		offset = 0
	}

	return service.store.List(requestContext, limit, offset)
}

const (
	defaultPageSize = 50
	maximumPageSize = 200
)

// statusRecorder captures the status a handler wrote, so refused attempts are
// logged alongside successful ones.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (recorder *statusRecorder) WriteHeader(status int) {
	recorder.status = status
	recorder.ResponseWriter.WriteHeader(status)
}

func (recorder *statusRecorder) Write(payload []byte) (int, error) {
	// A handler that writes a body without calling WriteHeader has implicitly
	// answered 200; without this the entry would record a zero status.
	if recorder.status == 0 {
		recorder.status = http.StatusOK
	}

	return recorder.ResponseWriter.Write(payload)
}

// Record wraps authenticated admin routes and logs every request that changes
// something.
//
// It is middleware rather than a call inside each handler because coverage is
// the whole value: a log with gaps in it is worse than none, since it invites
// the belief that anything absent did not happen. Wrapping the mux means a
// route added later is audited without anyone remembering to add it.
func (service *Service) Record(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet || request.Method == http.MethodHead {
			next.ServeHTTP(responseWriter, request)
			return
		}

		recorder := &statusRecorder{ResponseWriter: responseWriter}
		next.ServeHTTP(recorder, request)

		session, authenticated := admin.SessionFrom(request.Context())
		if !authenticated {
			// Unauthenticated requests never reached anything to change, and
			// have no actor to attribute. Sign-in attempts are throttled and
			// logged separately.
			return
		}

		entry := Entry{
			ActorEmail: session.User.Email,
			ActorRole:  session.User.Role,
			Method:     request.Method,
			Path:       request.URL.Path,
			StatusCode: recorder.status,
			OccurredAt: time.Now().UTC(),
		}

		// Deliberately not in the request's context: it is cancelled the moment
		// the response is written, which would drop the entry for exactly the
		// fast-returning requests there are most of.
		recordContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := service.store.Record(recordContext, entry, session.User.ID); err != nil {
			// The action itself already happened and its response has been
			// sent, so this cannot fail the request. It is logged loudly
			// because a silent gap in an audit trail is the worst outcome here.
			log.Printf("AUDIT WRITE FAILED actor=%s %s %s: %v",
				entry.ActorEmail, entry.Method, entry.Path, err)
		}
	})
}
