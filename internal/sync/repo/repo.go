package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	pgx "github.com/jackc/pgx/v5"

	"github.com/GlebRadaev/password-manager/internal/common/pg"
	"github.com/GlebRadaev/password-manager/internal/sync/models"
)

var (
	ErrConflictNotFound = errors.New("conflict not found")
)

type Repository struct {
	db pg.Database
}

func New(db pg.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetConflict(ctx context.Context, conflictID string) (models.Conflict, error) {
	query := `
		SELECT id, user_id, data_id, client_data, server_data, resolved, created_at, updated_at
		FROM sync.conflicts
		WHERE id = $1`

	var conflict models.Conflict
	err := r.db.QueryRow(ctx, query, conflictID).Scan(
		&conflict.ID,
		&conflict.UserID,
		&conflict.DataID,
		&conflict.ClientData,
		&conflict.ServerData,
		&conflict.Resolved,
		&conflict.CreatedAt,
		&conflict.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Conflict{}, ErrConflictNotFound
		}
		return models.Conflict{}, fmt.Errorf("failed to get conflict: %w", err)
	}

	return conflict, nil
}

func (r *Repository) AddConflicts(ctx context.Context, conflicts []models.Conflict) error {
	if len(conflicts) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, conflict := range conflicts {
		query := `
			INSERT INTO sync.conflicts (id, user_id, data_id, client_data, server_data, resolved, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		batch.Queue(query,
			conflict.ID,
			conflict.UserID,
			conflict.DataID,
			conflict.ClientData,
			conflict.ServerData,
			false,
			time.Now(),
			time.Now(),
		)
	}

	results := r.db.SendBatch(ctx, batch)
	defer results.Close()

	for i := 0; i < len(conflicts); i++ {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("failed to add conflict at index %d: %w", i, err)
		}
	}

	if err := results.Close(); err != nil {
		return fmt.Errorf("failed to close batch results: %w", err)
	}

	return nil
}

func (r *Repository) GetUnresolvedConflicts(ctx context.Context, userID string) ([]models.Conflict, error) {
	query := `
		SELECT id, user_id, data_id, client_data, server_data, resolved, created_at, updated_at
		FROM sync.conflicts
		WHERE user_id = $1 AND resolved = FALSE`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unresolved conflicts: %w", err)
	}
	defer rows.Close()

	var conflicts []models.Conflict
	for rows.Next() {
		var conflict models.Conflict
		if err := rows.Scan(
			&conflict.ID,
			&conflict.UserID,
			&conflict.DataID,
			&conflict.ClientData,
			&conflict.ServerData,
			&conflict.Resolved,
			&conflict.CreatedAt,
			&conflict.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan conflict: %w", err)
		}
		conflicts = append(conflicts, conflict)
	}

	return conflicts, nil
}

func (r *Repository) ResolveConflict(ctx context.Context, conflictID string) error {
	query := `
		UPDATE sync.conflicts
		SET resolved = TRUE, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	_, err := r.db.Exec(ctx, query, conflictID)
	if err != nil {
		return fmt.Errorf("failed to resolve conflict: %w", err)
	}

	return nil
}

func (r *Repository) DeleteConflicts(ctx context.Context, conflictIDs []string) error {
	query := `
		DELETE FROM sync.conflicts
		WHERE id = ANY($1)`

	_, err := r.db.Exec(ctx, query, conflictIDs)
	if err != nil {
		return fmt.Errorf("failed to delete conflicts: %w", err)
	}

	return nil
}
