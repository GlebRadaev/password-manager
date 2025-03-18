package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/GlebRadaev/password-manager/internal/common/pg"
	"github.com/GlebRadaev/password-manager/internal/data/models"
	"github.com/jackc/pgx/v5"
)

var (
	ErrDataNotFound = errors.New("data not found")
)

type Repository struct {
	db pg.Database
}

func New(db pg.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateData(ctx context.Context, data models.Data) error {
	query := `
		INSERT INTO data.data (id, user_id, type, data, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	metadataJSON, err := json.Marshal(data.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = r.db.Exec(ctx, query,
		data.ID,
		data.UserID,
		data.Type,
		data.Data,
		metadataJSON,
		data.CreatedAt,
		data.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create data: %w", err)
	}

	return nil
}

func (r *Repository) GetData(ctx context.Context, userID, dataID string) (models.Data, error) {
	query := `
		SELECT id, user_id, type, data, metadata, created_at, updated_at
		FROM data.data
		WHERE id = $1 AND user_id = $2`

	var data models.Data
	var metadataJSON []byte

	err := r.db.QueryRow(ctx, query, dataID, userID).Scan(
		&data.ID,
		&data.UserID,
		&data.Type,
		&data.Data,
		&metadataJSON,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Data{}, ErrDataNotFound
		}
		return models.Data{}, fmt.Errorf("failed to get data: %w", err)
	}

	if err := json.Unmarshal(metadataJSON, &data.Metadata); err != nil {
		return models.Data{}, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return data, nil
}

func (r *Repository) UpdateData(ctx context.Context, data models.Data) error {
	query := `
		UPDATE data.data
		SET type = $1, data = $2, metadata = $3, updated_at = $4
		WHERE id = $5 AND user_id = $6`

	metadataJSON, err := json.Marshal(data.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = r.db.Exec(ctx, query,
		data.Type,
		data.Data,
		metadataJSON,
		data.UpdatedAt,
		data.ID,
		data.UserID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrDataNotFound
		}
		return fmt.Errorf("failed to update data: %w", err)
	}

	return nil
}

func (r *Repository) DeleteData(ctx context.Context, userID, dataID string) error {
	query := `
		DELETE FROM data.data
		WHERE id = $1 AND user_id = $2`

	_, err := r.db.Exec(ctx, query, dataID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrDataNotFound
		}
		return fmt.Errorf("failed to delete data: %w", err)
	}

	return nil
}

func (r *Repository) ListData(ctx context.Context, userID string) ([]models.Data, error) {
	query := `
		SELECT id, user_id, type, data, metadata, created_at, updated_at
		FROM data.data
		WHERE user_id = $1`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list data: %w", err)
	}
	defer rows.Close()

	var dataList []models.Data
	for rows.Next() {
		var data models.Data
		var metadataJSON []byte

		err := rows.Scan(
			&data.ID,
			&data.UserID,
			&data.Type,
			&data.Data,
			&metadataJSON,
			&data.CreatedAt,
			&data.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan data: %w", err)
		}

		if err := json.Unmarshal(metadataJSON, &data.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		dataList = append(dataList, data)
	}

	return dataList, nil
}
