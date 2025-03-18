package repo

import (
	"context"
	"fmt"

	"github.com/GlebRadaev/password-manager/internal/common/pg"
	"github.com/GlebRadaev/password-manager/internal/sync/models"
)

type Repository struct {
	db pg.Database
}

func New(db pg.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveChanges(ctx context.Context, changes []models.DataChange) error {
	query := `
		INSERT INTO sync.changes (id, user_id, data_id, type, data, metadata, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	for _, change := range changes {
		_, err := r.db.Exec(ctx, query, change.ID, change.UserID, change.DataID, change.Type, change.Data, change.Metadata, change.Timestamp.Unix())
		if err != nil {
			return fmt.Errorf("failed to save change: %w", err)
		}
	}

	return nil
}

func (r *Repository) GetChanges(ctx context.Context, userID string, lastSyncTime int64) ([]models.DataChange, error) {
	query := `
		SELECT id, data_id, type, data, metadata, timestamp
		FROM sync.changes
		WHERE user_id = $1 AND timestamp > $2`

	rows, err := r.db.Query(ctx, query, userID, lastSyncTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get changes: %w", err)
	}
	defer rows.Close()

	var changes []models.DataChange
	for rows.Next() {
		var change models.DataChange
		if err := rows.Scan(&change.ID, &change.DataID, &change.Type, &change.Data, &change.Metadata, &change.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan change: %w", err)
		}
		changes = append(changes, change)
	}

	return changes, nil
}

func (r *Repository) GetChangesByDataID(ctx context.Context, userID, dataID string) ([]models.DataChange, error) {
	query := `
		SELECT id, data_id, type, data, metadata, timestamp
		FROM sync.changes
		WHERE user_id = $1 AND data_id = $2`

	rows, err := r.db.Query(ctx, query, userID, dataID)
	if err != nil {
		return nil, fmt.Errorf("failed to get changes by data_id: %w", err)
	}
	defer rows.Close()

	var changes []models.DataChange
	for rows.Next() {
		var change models.DataChange
		if err := rows.Scan(&change.ID, &change.DataID, &change.Type, &change.Data, &change.Metadata, &change.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan change: %w", err)
		}
		changes = append(changes, change)
	}

	return changes, nil
}
