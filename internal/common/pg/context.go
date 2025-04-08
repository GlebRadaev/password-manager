// Package pg provides PostgreSQL transaction context utilities
package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type key int

const (
	keyWithTX key = iota // context key for transaction
)

// TXFrom extracts transaction from context if present
// Returns (tx, true) if found, (nil, false) otherwise
func TXFrom(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(keyWithTX).(pgx.Tx)
	return tx, ok
}

// ContextOpt defines context modification function
type ContextOpt func(ctx context.Context) context.Context

// With applies all context options to the given context
func With(ctx context.Context, opts ...ContextOpt) context.Context {
	for _, o := range opts {
		ctx = o(ctx)
	}
	return ctx
}

// WithTransaction creates option to add transaction to context
func WithTransaction(tx pgx.Tx) ContextOpt {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, keyWithTX, tx)
	}
}
