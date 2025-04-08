// Package pg provides PostgreSQL transaction management
package pg

//go:generate mockgen -destination=manager_mock.go -source=manager.go -package=pg
import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// TXManager defines transaction management interface
type TXManager interface {
	// Begin executes fn in transaction, handling commit/rollback
	Begin(ctx context.Context, fn TransactionalFn) error
}

// Manager implements TXManager for PostgreSQL
type Manager struct {
	db Database
}

// NewTXManager creates new transaction manager
func NewTXManager(db Database) *Manager {
	return &Manager{db: db}
}

// TransactionalFn represents function to run in transaction
type TransactionalFn func(ctx context.Context) (err error)

// Begin executes fn in transaction
func (mgr *Manager) Begin(ctx context.Context, fn TransactionalFn) (err error) {
	_, ok := TXFrom(ctx)
	if ok {
		return fn(ctx)
	}

	tx, err := mgr.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pg: can't begin tx: %w", err)
	}
	ctx = With(ctx, WithTransaction(tx))
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
				err = fmt.Errorf("pg: can't rollback tx: %w, original error: %w", rollbackErr, err)
			}
		}
	}()
	err = fn(ctx)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("pg: can't commit tx: %w", err)
	}
	return
}
