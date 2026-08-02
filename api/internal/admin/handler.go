package admin

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type contextKey string

const sessionContextKey contextKey = "admin-session"

const (
	maximumLoginBytes    = 4 << 10
	loginAttemptWindow   = 5 * time.Minute
	loginAttemptsAllowed = 8
)

type Handler struct {
	authService *AuthService
	// Production serves the admin over HTTPS; the Secure cookie flag is only
	// dropped in development, where the browser would otherwise discard it.
	useSecureCookie bool
	limiter         *loginLimiter
}

func NewHandler(authService *AuthService, useSecureCookie bool) *Handler {
	return &Handler{
		authService:     authService,
		useSecureCookie: useSecureCookie,
		limiter:         newLoginLimiter(),
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type sessionResponse struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// LogIn authenticates and sets the session cookie.
func (handler *Handler) LogIn(responseWriter http.ResponseWriter, request *http.Request) {
	// Throttle by client address so the login form cannot be brute forced.
	address := clientAddress(request)
	if !handler.limiter.allow(address) {
		writeJSON(responseWriter, http.StatusTooManyRequests,
			map[string]string{"error": "too many attempts, try again shortly"})
		return
	}

	var credentials loginRequest
	decoder := json.NewDecoder(http.MaxBytesReader(responseWriter, request.Body, maximumLoginBytes))
	if err := decoder.Decode(&credentials); err != nil {
		writeJSON(responseWriter, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	session, err := handler.authService.LogIn(
		request.Context(), credentials.Email, credentials.Password,
	)
	if err != nil {
		// Only failures count, so signing in and out through the day never
		// locks the owner out of their own shop.
		handler.limiter.recordFailure(address)

		// One message for every failure: no hint about which part was wrong.
		writeJSON(responseWriter, http.StatusUnauthorized,
			map[string]string{"error": "invalid email or password"})
		return
	}

	http.SetCookie(responseWriter, &http.Cookie{
		Name:     SessionCookie,
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   handler.useSecureCookie,
		SameSite: http.SameSiteLaxMode,
	})

	writeJSON(responseWriter, http.StatusOK, sessionResponse{
		Email: session.User.Email,
		Role:  session.User.Role,
	})
}

// LogOut clears the session both server-side and in the browser.
func (handler *Handler) LogOut(responseWriter http.ResponseWriter, request *http.Request) {
	if cookie, err := request.Cookie(SessionCookie); err == nil {
		_ = handler.authService.LogOut(request.Context(), cookie.Value)
	}

	http.SetCookie(responseWriter, &http.Cookie{
		Name:     SessionCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   handler.useSecureCookie,
		SameSite: http.SameSiteLaxMode,
	})

	responseWriter.WriteHeader(http.StatusNoContent)
}

// Me lets the admin UI check whether the current session is still valid.
func (handler *Handler) Me(responseWriter http.ResponseWriter, request *http.Request) {
	session, ok := SessionFrom(request.Context())
	if !ok {
		writeJSON(responseWriter, http.StatusUnauthorized, map[string]string{"error": "not signed in"})
		return
	}

	writeJSON(responseWriter, http.StatusOK, sessionResponse{
		Email: session.User.Email,
		Role:  session.User.Role,
	})
}

// RequireSession gates every admin endpoint behind a valid session.
func (handler *Handler) RequireSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie(SessionCookie)
		if err != nil {
			writeJSON(responseWriter, http.StatusUnauthorized, map[string]string{"error": "not signed in"})
			return
		}

		session, err := handler.authService.Authenticate(request.Context(), cookie.Value)
		if err != nil {
			if !errors.Is(err, ErrNotAuthenticated) {
				log.Printf("authenticate admin session: %v", err)
			}
			writeJSON(responseWriter, http.StatusUnauthorized, map[string]string{"error": "not signed in"})
			return
		}

		next.ServeHTTP(responseWriter, request.WithContext(
			context.WithValue(request.Context(), sessionContextKey, session),
		))
	})
}

// SessionFrom reads the authenticated session placed by RequireSession.
func SessionFrom(requestContext context.Context) (Session, bool) {
	session, ok := requestContext.Value(sessionContextKey).(Session)

	return session, ok
}

// loginLimiter is a small in-memory throttle. A single admin does not justify
// shared state; the process restarting simply resets the window.
type loginLimiter struct {
	mutex    sync.Mutex
	attempts map[string][]time.Time
}

func newLoginLimiter() *loginLimiter {
	return &loginLimiter{attempts: make(map[string][]time.Time)}
}

// allow reports whether the caller still has failed attempts left in the window.
func (limiter *loginLimiter) allow(key string) bool {
	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	recent := limiter.pruneLocked(key)

	return len(recent) < loginAttemptsAllowed
}

// recordFailure charges one attempt against the window.
func (limiter *loginLimiter) recordFailure(key string) {
	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	limiter.attempts[key] = append(limiter.pruneLocked(key), time.Now())
}

// pruneLocked drops attempts that have aged out. Callers hold the mutex.
func (limiter *loginLimiter) pruneLocked(key string) []time.Time {
	cutoff := time.Now().Add(-loginAttemptWindow)
	recent := make([]time.Time, 0, loginAttemptsAllowed)
	for _, attempt := range limiter.attempts[key] {
		if attempt.After(cutoff) {
			recent = append(recent, attempt)
		}
	}
	limiter.attempts[key] = recent

	return recent
}

// clientAddress keys the throttle by IP only. RemoteAddr carries a fresh port
// per connection, so including it would give every attempt its own bucket and
// the limiter would never fire.
func clientAddress(request *http.Request) string {
	if forwarded := request.Header.Get("X-Forwarded-For"); forwarded != "" {
		// Caddy appends; the original client is first.
		if comma := strings.IndexByte(forwarded, ','); comma > 0 {
			return strings.TrimSpace(forwarded[:comma])
		}

		return strings.TrimSpace(forwarded)
	}

	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		return request.RemoteAddr
	}

	return host
}

func writeJSON(responseWriter http.ResponseWriter, statusCode int, payload any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)

	if err := json.NewEncoder(responseWriter).Encode(payload); err != nil {
		log.Printf("encode admin response: %v", err)
	}
}
