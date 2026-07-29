package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const migrationSourceName = "iofs"

// migrationFiles packages every versioned SQL file into the API binary.
//
//go:embed *.sql
var migrationFiles embed.FS

// Up applies every pending database migration in version order.
func Up(databaseURL string) error {
	// Build a migration source from the SQL files embedded with the API binary.
	migrationSource, err := iofs.New(migrationFiles, ".")
	if err != nil {
		return fmt.Errorf("create embedded migration source: %w", err)
	}

	// Open a dedicated migration connection so application pooling remains separate.
	databaseHandle, err := sql.Open("pgx", databaseURL)
	if err != nil {
		_ = migrationSource.Close()
		return fmt.Errorf("open migration database: %w", err)
	}
	databaseHandle.SetMaxOpenConns(1)

	// Adapt the pgx-backed SQL connection to the PostgreSQL migration driver.
	databaseDriver, err := postgres.WithInstance(databaseHandle, &postgres.Config{
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		_ = databaseHandle.Close()
		_ = migrationSource.Close()
		return fmt.Errorf("create postgres migration driver: %w", err)
	}

	// Pair the embedded source and PostgreSQL driver into one versioned runner.
	migrationRunner, err := migrate.NewWithInstance(
		migrationSourceName,
		migrationSource,
		"postgres",
		databaseDriver,
	)
	if err != nil {
		_ = databaseDriver.Close()
		_ = migrationSource.Close()
		return fmt.Errorf("create migration runner: %w", err)
	}

	// Release both source and database resources after migration completes.
	defer func() {
		_, _ = migrationRunner.Close()
	}()

	// Apply pending versions while treating an already-current schema as success.
	if err := migrationRunner.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply database migrations: %w", err)
	}

	return nil
}
