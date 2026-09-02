package admin

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrInvalidCredentials is deliberately returned for both an unknown email
	// and a wrong password, so the response never reveals which accounts exist.
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotAuthenticated   = errors.New("not authenticated")
)

const (
	// Sessions are long enough to be practical for a shop owner checking orders
	// through the day, short enough that a stolen cookie expires.
	SessionLifetime = 12 * time.Hour
	SessionCookie   = "oyecouture_admin_session"
	sessionBytes    = 32
	// bcrypt's default is 10; 12 costs a little more per login and a lot more
	// per offline guess.
	bcryptCost = 12
)

type User struct {
	ID    string
	Email string
	Role  string
}

type Session struct {
	Token     string
	ExpiresAt time.Time
	User      User
}

// Store is the persistence contract for admin authentication.
type Store interface {
	UserByEmail(context.Context, string) (User, string, error)
	CreateSession(context.Context, string, string, time.Time) error
	SessionByToken(context.Context, string) (Session, error)
	DeleteSession(context.Context, string) error
	DeleteExpiredSessions(context.Context) error
}

type AuthService struct {
	store Store
}

func NewAuthService(store Store) *AuthService {
	return &AuthService{store: store}
}

// LogIn verifies the password and issues a session.
func (service *AuthService) LogIn(
	requestContext context.Context,
	email string,
	password string,
) (Session, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	user, passwordHash, err := service.store.UserByEmail(requestContext, normalizedEmail)
	if err != nil {
		// Hash a dummy value so a missing account takes the same time as a wrong
		// password, leaving no timing signal about which emails exist.
		_ = bcrypt.CompareHashAndPassword(
			[]byte("$2a$12$vRFyqZ0mJm8Gy7Nn1qDcMuWnGqk1S1V1eYcM6M0y1nS9pQ9m8p1zK"),
			[]byte(password),
		)
		return Session{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return Session{}, ErrInvalidCredentials
	}

	token, err := newSessionToken()
	if err != nil {
		return Session{}, fmt.Errorf("generate session token: %w", err)
	}

	expiresAt := time.Now().Add(SessionLifetime)
	if err := service.store.CreateSession(requestContext, token, user.ID, expiresAt); err != nil {
		return Session{}, fmt.Errorf("create session: %w", err)
	}

	return Session{Token: token, ExpiresAt: expiresAt, User: user}, nil
}

// Authenticate resolves a session token, rejecting expired sessions.
func (service *AuthService) Authenticate(
	requestContext context.Context,
	token string,
) (Session, error) {
	if token == "" {
		return Session{}, ErrNotAuthenticated
	}

	session, err := service.store.SessionByToken(requestContext, token)
	if err != nil {
		return Session{}, ErrNotAuthenticated
	}
	if time.Now().After(session.ExpiresAt) {
		_ = service.store.DeleteSession(requestContext, token)
		return Session{}, ErrNotAuthenticated
	}

	return session, nil
}

func (service *AuthService) LogOut(requestContext context.Context, token string) error {
	if token == "" {
		return nil
	}

	return service.store.DeleteSession(requestContext, token)
}

// passwordMatches reports whether a plaintext password matches a stored hash.
// bcrypt's comparison is constant time for a given hash, so this leaks nothing
// beyond the answer itself.
func passwordMatches(passwordHash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) == nil
}

// HashPassword is used by the admin CLI when creating a user.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	return string(hash), nil
}

func newSessionToken() (string, error) {
	tokenBytes := make([]byte, sessionBytes)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(tokenBytes), nil
}
