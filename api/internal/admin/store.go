package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

// UserByEmail returns the user and their password hash for verification.
func (store *PostgresStore) UserByEmail(
	requestContext context.Context,
	email string,
) (User, string, error) {
	const query = `
		SELECT id, email, role, password_hash
		FROM admin_users
		WHERE lower(email) = lower($1)
	`

	var user User
	var passwordHash string
	if err := store.pool.QueryRow(requestContext, query, email).
		Scan(&user.ID, &user.Email, &user.Role, &passwordHash); err != nil {
		return User{}, "", fmt.Errorf("load admin user: %w", err)
	}

	return user, passwordHash, nil
}

func (store *PostgresStore) CreateSession(
	requestContext context.Context,
	token string,
	adminUserID string,
	expiresAt time.Time,
) error {
	const query = `
		INSERT INTO admin_sessions (token, admin_user_id, expires_at)
		VALUES ($1, $2, $3)
	`

	if _, err := store.pool.Exec(requestContext, query, token, adminUserID, expiresAt); err != nil {
		return fmt.Errorf("insert session: %w", err)
	}

	return nil
}

func (store *PostgresStore) SessionByToken(
	requestContext context.Context,
	token string,
) (Session, error) {
	const query = `
		SELECT s.token, s.expires_at, u.id, u.email, u.role
		FROM admin_sessions s
		JOIN admin_users u ON u.id = s.admin_user_id
		WHERE s.token = $1
	`

	var session Session
	if err := store.pool.QueryRow(requestContext, query, token).Scan(
		&session.Token, &session.ExpiresAt,
		&session.User.ID, &session.User.Email, &session.User.Role,
	); err != nil {
		return Session{}, fmt.Errorf("load session: %w", err)
	}

	return session, nil
}

func (store *PostgresStore) DeleteSession(requestContext context.Context, token string) error {
	if _, err := store.pool.Exec(
		requestContext, `DELETE FROM admin_sessions WHERE token = $1`, token,
	); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

// DeleteExpiredSessions is called on startup so dead rows do not accumulate.
func (store *PostgresStore) DeleteExpiredSessions(requestContext context.Context) error {
	if _, err := store.pool.Exec(
		requestContext, `DELETE FROM admin_sessions WHERE expires_at < NOW()`,
	); err != nil {
		return fmt.Errorf("delete expired sessions: %w", err)
	}

	return nil
}

// CreateUser is used by the admin CLI. Returns the new user's id.
func (store *PostgresStore) CreateUser(
	requestContext context.Context,
	email string,
	passwordHash string,
	role string,
) (string, error) {
	const query = `
		INSERT INTO admin_users (email, password_hash, role)
		VALUES (lower($1), $2, $3)
		ON CONFLICT (lower(email)) DO UPDATE SET password_hash = EXCLUDED.password_hash
		RETURNING id
	`

	var id string
	if err := store.pool.QueryRow(requestContext, query, email, passwordHash, role).Scan(&id); err != nil {
		return "", fmt.Errorf("create admin user: %w", err)
	}

	return id, nil
}
