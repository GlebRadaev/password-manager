// Package repo provides data access operations for password entries.
// It handles database interactions including CRUD operations and batch processing.
package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	pgx "github.com/jackc/pgx/v5"

	"github.com/GlebRadaev/password-manager/internal/common/pg"
	"github.com/GlebRadaev/password-manager/internal/data/models"
)

// ErrDataNotFound is returned when requested data is not found in the database.
var ErrDataNotFound = errors.New("data not found")

// Repo provides methods for interacting with the data storage.
type Repo struct {
	db pg.Database
}

// New creates a new Repo instance with the given database connection.
func New(db pg.Database) *Repo {
	return &Repo{db: db}
}

// AddList inserts multiple data entries in a batch and returns their IDs.
func (r *Repo) AddList(ctx context.Context, entries []models.DataEntry) ([]string, error) {
	batch := &pgx.Batch{}
	for _, entry := range entries {
		query := `
			INSERT INTO data.entries (id, user_id, type, data, metadata, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id`
		batch.Queue(query,
			entry.ID,
			entry.UserID,
			entry.Type,
			entry.Data,
			entry.Metadata,
			time.Now(),
			time.Now(),
		)
	}

	results := r.db.SendBatch(ctx, batch)
	defer results.Close()

	ids := make([]string, 0, len(entries))
	for i := 0; i < len(entries); i++ {
		var id string
		if err := results.QueryRow().Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to add data list: %w", err)
		}
		ids = append(ids, id)
	}

	if err := results.Close(); err != nil {
		return nil, fmt.Errorf("failed to close batch results: %w", err)
	}

	return ids, nil
}

// UpdateData modifies an existing data entry.
// Returns ErrDataNotFound if no matching entry exists.
func (r *Repo) UpdateData(ctx context.Context, entry models.DataEntry) error {
	query := `
		UPDATE data.entries
		SET type = $1, data = $2, metadata = $3, updated_at = $4
		WHERE id = $5 AND user_id = $6`
	result, err := r.db.Exec(ctx, query,
		entry.Type,
		entry.Data,
		entry.Metadata,
		time.Now(),
		entry.ID,
		entry.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update data: %w", err)
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrDataNotFound
	}
	return nil
}

// DeleteList removes multiple entries by their IDs for a specific user.
// Returns ErrDataNotFound if no matching entries exist.
func (r *Repo) DeleteList(ctx context.Context, userID string, dataIDs []string) error {
	if len(dataIDs) == 0 {
		return nil
	}

	query := `
		DELETE FROM data.entries
		WHERE user_id = $1 AND id = ANY($2)`

	result, err := r.db.Exec(ctx, query, userID, dataIDs)
	if err != nil {
		return fmt.Errorf("failed to delete data list: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrDataNotFound
	}

	return nil
}

// ListData retrieves all entries for a specific user.
func (r *Repo) ListData(ctx context.Context, userID string) ([]models.DataEntry, error) {
	query := `
        SELECT id, type, data, metadata, created_at, updated_at
        FROM data.entries
        WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list data: %w", err)
	}
	defer rows.Close()

	var entries []models.DataEntry
	for rows.Next() {
		var entry models.DataEntry
		var metadataJSON []byte

		if err := rows.Scan(
			&entry.ID,
			&entry.Type,
			&entry.Data,
			&metadataJSON,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan data entry: %w", err)
		}

		if err := json.Unmarshal(metadataJSON, &entry.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
