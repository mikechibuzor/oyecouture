package admin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	// active is checked here so a revoked admin cannot sign in. Without it,
	// revoking access would only take effect when their current session
	// happened to expire.
	const query = `
		SELECT id, email, role, password_hash
		FROM admin_users
		WHERE lower(email) = lower($1)
			AND active = TRUE
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
	// A live session belonging to a revoked admin must stop working at once,
	// so this checks the account on every authenticated request rather than
	// trusting the session alone.
	const query = `
		SELECT s.token, s.expires_at, u.id, u.email, u.role
		FROM admin_sessions s
		JOIN admin_users u ON u.id = s.admin_user_id
		WHERE s.token = $1
			AND u.active = TRUE
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

// ListMembers returns every admin account, revoked ones included, so the team
// desk shows who once had access rather than quietly hiding them.
func (store *PostgresStore) ListMembers(requestContext context.Context) ([]TeamMember, error) {
	const query = `
		SELECT id, email, role, active, created_at, updated_at
		FROM admin_users
		ORDER BY active DESC, lower(email)
	`

	rows, err := store.pool.Query(requestContext, query)
	if err != nil {
		return nil, fmt.Errorf("list admin users: %w", err)
	}
	defer rows.Close()

	members := []TeamMember{}
	for rows.Next() {
		var member TeamMember
		if err := rows.Scan(
			&member.ID, &member.Email, &member.Role,
			&member.Active, &member.CreatedAt, &member.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan admin user: %w", err)
		}
		members = append(members, member)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read admin users: %w", err)
	}

	return members, nil
}

func (store *PostgresStore) MemberByID(
	requestContext context.Context,
	id string,
) (TeamMember, error) {
	const query = `
		SELECT id, email, role, active, created_at, updated_at
		FROM admin_users
		WHERE id = $1
	`

	var member TeamMember
	if err := store.pool.QueryRow(requestContext, query, id).Scan(
		&member.ID, &member.Email, &member.Role,
		&member.Active, &member.CreatedAt, &member.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TeamMember{}, ErrTeamNotFound
		}
		return TeamMember{}, fmt.Errorf("load admin user: %w", err)
	}

	return member, nil
}

// CreateMember refuses rather than overwrites on a duplicate email. CreateUser
// upserts because the CLI is also the password-reset path; here an accidental
// duplicate would silently replace a colleague's password.
func (store *PostgresStore) CreateMember(
	requestContext context.Context,
	email string,
	passwordHash string,
	role string,
) (TeamMember, error) {
	const query = `
		INSERT INTO admin_users (email, password_hash, role)
		VALUES (lower($1), $2, $3)
		RETURNING id, email, role, active, created_at, updated_at
	`

	var member TeamMember
	if err := store.pool.QueryRow(requestContext, query, email, passwordHash, role).Scan(
		&member.ID, &member.Email, &member.Role,
		&member.Active, &member.CreatedAt, &member.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return TeamMember{}, ErrEmailTaken
		}
		return TeamMember{}, fmt.Errorf("create admin user: %w", err)
	}

	return member, nil
}

func (store *PostgresStore) UpdateMember(
	requestContext context.Context,
	id string,
	role string,
	active bool,
) (TeamMember, error) {
	const query = `
		UPDATE admin_users
		SET role = $2, active = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, email, role, active, created_at, updated_at
	`

	var member TeamMember
	if err := store.pool.QueryRow(requestContext, query, id, role, active).Scan(
		&member.ID, &member.Email, &member.Role,
		&member.Active, &member.CreatedAt, &member.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TeamMember{}, ErrTeamNotFound
		}
		return TeamMember{}, fmt.Errorf("update admin user: %w", err)
	}

	return member, nil
}

// DeleteMember relies on admin_sessions.admin_user_id being ON DELETE CASCADE,
// so a deleted admin's sessions go with them.
func (store *PostgresStore) DeleteMember(requestContext context.Context, id string) error {
	tag, err := store.pool.Exec(requestContext, `DELETE FROM admin_users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete admin user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrTeamNotFound
	}

	return nil
}

func (store *PostgresStore) CountActiveOwners(requestContext context.Context) (int, error) {
	var owners int
	if err := store.pool.QueryRow(
		requestContext,
		`SELECT COUNT(*) FROM admin_users WHERE role = 'owner' AND active = TRUE`,
	).Scan(&owners); err != nil {
		return 0, fmt.Errorf("count owners: %w", err)
	}

	return owners, nil
}

func (store *PostgresStore) PasswordHashByID(
	requestContext context.Context,
	id string,
) (string, error) {
	var passwordHash string
	if err := store.pool.QueryRow(
		requestContext,
		`SELECT password_hash FROM admin_users WHERE id = $1 AND active = TRUE`,
		id,
	).Scan(&passwordHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrTeamNotFound
		}
		return "", fmt.Errorf("load password hash: %w", err)
	}

	return passwordHash, nil
}

func (store *PostgresStore) UpdatePassword(
	requestContext context.Context,
	id string,
	passwordHash string,
) error {
	tag, err := store.pool.Exec(
		requestContext,
		`UPDATE admin_users SET password_hash = $2, updated_at = NOW() WHERE id = $1`,
		id, passwordHash,
	)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrTeamNotFound
	}

	return nil
}

func (store *PostgresStore) DeleteSessionsForUser(
	requestContext context.Context,
	id string,
) error {
	if _, err := store.pool.Exec(
		requestContext,
		`DELETE FROM admin_sessions WHERE admin_user_id = $1`,
		id,
	); err != nil {
		return fmt.Errorf("delete sessions for user: %w", err)
	}

	return nil
}
