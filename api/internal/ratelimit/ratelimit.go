// Package ratelimit throttles abuse-prone endpoints by client address.
//
// It lives on its own because two very different endpoints need it: admin
// sign-in and the public contact form. Duplicating the address handling
// is how a fixed bug comes back. Keying on RemoteAddr including its ephemeral
// port once gave every request its own bucket, so nothing was ever limited.
package ratelimit

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Limiter is a small in-memory throttle. Shared state across processes is not
// worth the complexity here; a restart simply resets the window.
//
// Checking and charging are separate on purpose. Sign-in charges only failed
// attempts, so an owner is never locked out by working normally, while the
// contact form charges every submission, because the abuse there is volume of
// successful sends rather than guessing.
type Limiter struct {
	mutex    sync.Mutex
	attempts map[string][]time.Time
	window   time.Duration
	allowed  int
}

func New(window time.Duration, allowed int) *Limiter {
	return &Limiter{
		attempts: make(map[string][]time.Time),
		window:   window,
		allowed:  allowed,
	}
}

// Allow reports whether the caller still has attempts left in the window. It
// does not charge one; call Record for that.
func (limiter *Limiter) Allow(key string) bool {
	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	return len(limiter.pruneLocked(key)) < limiter.allowed
}

// Record charges one attempt against the window.
func (limiter *Limiter) Record(key string) {
	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	limiter.attempts[key] = append(limiter.pruneLocked(key), time.Now())
}

// pruneLocked drops attempts that have aged out. Callers hold the mutex.
func (limiter *Limiter) pruneLocked(key string) []time.Time {
	cutoff := time.Now().Add(-limiter.window)
	recent := make([]time.Time, 0, limiter.allowed)
	for _, attempt := range limiter.attempts[key] {
		if attempt.After(cutoff) {
			recent = append(recent, attempt)
		}
	}

	if len(recent) == 0 {
		// Do not let the map grow a permanent entry per address seen.
		delete(limiter.attempts, key)
	} else {
		limiter.attempts[key] = recent
	}

	return recent
}

// ClientAddress keys a throttle by IP only. RemoteAddr carries a fresh port per
// connection, so including it would give every request its own bucket and the
// limiter would never fire.
func ClientAddress(request *http.Request) string {
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
