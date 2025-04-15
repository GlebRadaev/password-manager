// Package migrations provides database migration utilities
package migrations

import (
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS // Embedded SQL migration files

// Exec runs all pending database migrations using the provided connection pool
// Returns error if any migration fails
func Exec(pool *pgxpool.Pool) error {
	goose.SetBaseFS(embedMigrations) // Use embedded migrations

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool) // Convert pool to standard database handle
	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	return nil
}
