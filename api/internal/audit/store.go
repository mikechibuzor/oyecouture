package audit

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

// Record appends one entry. actor_email is stored alongside the id so the trail
// still names who acted after that admin is deleted.
func (store *PostgresStore) Record(
	requestContext context.Context,
	entry Entry,
	actorAdminID string,
) error {
	const query = `
		INSERT INTO admin_audit_log (
			actor_admin_id, actor_email, actor_role, method, path, status_code, occurred_at
		) VALUES (NULLIF($1, '')::uuid, $2, $3, $4, $5, $6, $7)
	`

	if _, err := store.pool.Exec(
		requestContext, query,
		actorAdminID, entry.ActorEmail, entry.ActorRole,
		entry.Method, entry.Path, entry.StatusCode, entry.OccurredAt,
	); err != nil {
		return fmt.Errorf("record audit entry: %w", err)
	}

	return nil
}

func (store *PostgresStore) List(
	requestContext context.Context,
	limit int,
	offset int,
) (Page, error) {
	const query = `
		SELECT id, actor_email, actor_role, method, path, status_code, occurred_at,
			COUNT(*) OVER() AS total
		FROM admin_audit_log
		ORDER BY occurred_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := store.pool.Query(requestContext, query, limit, offset)
	if err != nil {
		return Page{}, fmt.Errorf("list audit entries: %w", err)
	}
	defer rows.Close()

	page := Page{Entries: []Entry{}}
	for rows.Next() {
		var entry Entry
		if err := rows.Scan(
			&entry.ID, &entry.ActorEmail, &entry.ActorRole,
			&entry.Method, &entry.Path, &entry.StatusCode, &entry.OccurredAt,
			&page.Total,
		); err != nil {
			return Page{}, fmt.Errorf("scan audit entry: %w", err)
		}
		page.Entries = append(page.Entries, entry)
	}
	if err := rows.Err(); err != nil {
		return Page{}, fmt.Errorf("read audit entries: %w", err)
	}

	// COUNT(*) OVER() only produces a value on a row, so an empty page past the
	// end would report zero and make the list look empty rather than exhausted.
	if len(page.Entries) == 0 && offset > 0 {
		if err := store.pool.QueryRow(
			requestContext, `SELECT COUNT(*) FROM admin_audit_log`,
		).Scan(&page.Total); err != nil {
			return Page{}, fmt.Errorf("count audit entries: %w", err)
		}
	}

	return page, nil
}
