// Package pg provides PostgreSQL database abstractions
package pg

//go:generate mockgen -destination=spec_mock.go -source=spec.go -package=pg
import (
	"context"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Database defines core PostgreSQL operations
type Database interface {
	Begin(ctx context.Context) (pgx.Tx, error)                                    // Starts transaction
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)         // Executes query
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row                // Executes query returning single row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) // Executes command
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults                 // Executes batch
}

// TX represents a database transaction (aliases pgx.Tx)
type TX interface {
	pgx.Tx
}

// Rows represents query results (aliases pgx.Rows)
type Rows interface {
	pgx.Rows
}

// Row represents single query result row (aliases pgx.Row)
type Row interface {
	pgx.Row
}
