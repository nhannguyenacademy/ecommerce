// Package migrate contains the database schema, migrations and seeding data.
package migrate

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
)

var (
	//go:embed seeds/seed.sql
	seedDoc string
)

// Migrate attempts to bring the database up to date with the migrations
// defined in this package.
func Migrate(ctx context.Context, db *sqlx.DB) error {
	if err := sqldb.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("construct postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/sdk/migrate/migrations",
		"postgres-migrate-up", // for logging purposes only
		driver,
	)
	if err != nil {
		return fmt.Errorf("construct migrate instance: %w", err)
	}

	return m.Up()
}

func MigrateDown(ctx context.Context, db *sqlx.DB) error {
	if err := sqldb.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("construct postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/sdk/migrate/migrations",
		"postgres-migrate-down", // for logging purposes only
		driver,
	)
	if err != nil {
		return fmt.Errorf("construct migrate instance: %w", err)
	}

	return m.Down()
}

// Seed runs the seed document defined in this package against db. The queries
// are run in a transaction and rolled back if any fail.
func Seed(ctx context.Context, db *sqlx.DB) (err error) {
	if err := sqldb.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if errTx := tx.Rollback(); errTx != nil {
			if errors.Is(errTx, sql.ErrTxDone) {
				return
			}

			err = fmt.Errorf("rollback: %w", errTx)
			return
		}
	}()

	if _, err := tx.Exec(seedDoc); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}