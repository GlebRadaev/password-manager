package repo

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"

	"github.com/GlebRadaev/password-manager/internal/data/models"
)

func NewMock(t *testing.T) (*Repo, pgxmock.PgxPoolIface) {
	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	repo := New(mockDB)

	t.Cleanup(func() {
		mockDB.Close()
	})

	return repo, mockDB
}

func TestRepo_AddList(t *testing.T) {
	repo, mock := NewMock(t)

	entries := []models.DataEntry{
		{
			ID:     "id1",
			UserID: "user1",
			Type:   models.LoginPassword,
			Data:   []byte("data1"),
			Metadata: []models.Metadata{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
			},
		},
		{
			ID:     "id2",
			UserID: "user1",
			Type:   models.Text,
			Data:   []byte("data2"),
			Metadata: []models.Metadata{
				{Key: "key3", Value: "value3"},
			},
		},
	}

	tests := []struct {
		name      string
		mockSetup func(*pgxmock.PgxPoolIface)
		expectErr bool
	}{
		{
			name: "successful batch insert",
			mockSetup: func(mock *pgxmock.PgxPoolIface) {
				eb := (*mock).ExpectBatch()

				eb.ExpectQuery(regexp.QuoteMeta(`
                    INSERT INTO data.entries (id, user_id, type, data, metadata, created_at, updated_at)
                    VALUES ($1, $2, $3, $4, $5, $6, $7)
                    RETURNING id`)).
					WithArgs(
						entries[0].ID,
						entries[0].UserID,
						entries[0].Type,
						entries[0].Data,
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(entries[0].ID))

				eb.ExpectQuery(regexp.QuoteMeta(`
                    INSERT INTO data.entries (id, user_id, type, data, metadata, created_at, updated_at)
                    VALUES ($1, $2, $3, $4, $5, $6, $7)
                    RETURNING id`)).
					WithArgs(
						entries[1].ID,
						entries[1].UserID,
						entries[1].Type,
						entries[1].Data,
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(entries[1].ID))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(&mock)

			ids, err := repo.AddList(context.Background(), entries)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to add data list")
			} else {
				assert.NoError(t, err)
				assert.Len(t, ids, len(entries))
				assert.Equal(t, entries[0].ID, ids[0])
				assert.Equal(t, entries[1].ID, ids[1])
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepo_UpdateData(t *testing.T) {
	repo, mock := NewMock(t)

	entry := models.DataEntry{
		ID:     "id1",
		UserID: "user1",
		Type:   models.Card,
		Data:   []byte("newdata"),
		Metadata: []models.Metadata{
			{Key: "newkey", Value: "newvalue"},
		},
	}

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
	}{
		{
			name: "successful update",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`
                    UPDATE data.entries
                    SET type = $1, data = $2, metadata = $3, updated_at = $4
                    WHERE id = $5 AND user_id = $6`)).
					WithArgs(
						entry.Type,
						entry.Data,
						entry.Metadata,
						pgxmock.AnyArg(),
						entry.ID,
						entry.UserID,
					).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			expectErr: false,
		},
		{
			name: "data not found",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE data.entries`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))
			},
			expectErr: true,
		},
		{
			name: "database error",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE data.entries`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnError(errors.New("db error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.UpdateData(context.Background(), entry)

			if tt.expectErr {
				assert.Error(t, err)
				if errors.Is(err, ErrDataNotFound) {
					assert.Contains(t, err.Error(), "data not found")
				} else {
					assert.Contains(t, err.Error(), "failed to update data")
				}
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepo_DeleteList(t *testing.T) {
	repo, mock := NewMock(t)

	userID := "user1"
	dataIDs := []string{"id1", "id2"}

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
	}{
		{
			name: "successful delete",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`
					DELETE FROM data.entries
					WHERE user_id = $1 AND id = ANY($2)`)).
					WithArgs(userID, dataIDs).
					WillReturnResult(pgxmock.NewResult("DELETE", 2))
			},
			expectErr: false,
		},
		{
			name: "data not found",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM data.entries`)).
					WithArgs(userID, dataIDs).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))
			},
			expectErr: true,
		},
		{
			name: "database error",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM data.entries`)).
					WithArgs(userID, dataIDs).
					WillReturnError(errors.New("db error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.DeleteList(context.Background(), userID, dataIDs)

			if tt.expectErr {
				assert.Error(t, err)
				if errors.Is(err, ErrDataNotFound) {
					assert.Contains(t, err.Error(), "data not found")
				} else {
					assert.Contains(t, err.Error(), "failed to delete data list")
				}
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepo_ListData(t *testing.T) {
	repo, mock := NewMock(t)

	userID := "user1"
	now := time.Now()
	expectedEntries := []models.DataEntry{
		{
			ID:   "id1",
			Type: models.LoginPassword,
			Data: []byte("data1"),
			Metadata: []models.Metadata{
				{Key: "key1", Value: "value1"},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:   "id2",
			Type: models.Binary,
			Data: []byte("data2"),
			Metadata: []models.Metadata{
				{Key: "key2", Value: "value2"},
				{Key: "key3", Value: "value3"},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
		expected  []models.DataEntry
	}{
		{
			name: "successful list",
			mockSetup: func() {
				rows := pgxmock.NewRows([]string{"id", "type", "data", "metadata", "created_at", "updated_at"}).
					AddRow(
						expectedEntries[0].ID,
						expectedEntries[0].Type,
						expectedEntries[0].Data,
						[]byte(`[{"key":"key1","value":"value1"}]`),
						expectedEntries[0].CreatedAt,
						expectedEntries[0].UpdatedAt,
					).
					AddRow(
						expectedEntries[1].ID,
						expectedEntries[1].Type,
						expectedEntries[1].Data,
						[]byte(`[{"key":"key2","value":"value2"},{"key":"key3","value":"value3"}]`),
						expectedEntries[1].CreatedAt,
						expectedEntries[1].UpdatedAt,
					)
				mock.ExpectQuery(regexp.QuoteMeta(`
                    SELECT id, type, data, metadata, created_at, updated_at
                    FROM data.entries
                    WHERE user_id = $1`)).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectErr: false,
			expected:  expectedEntries,
		},
		{
			name: "no data found",
			mockSetup: func() {
				rows := pgxmock.NewRows([]string{"id", "type", "data", "metadata", "created_at", "updated_at"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectErr: false,
			expected:  nil,
		},
		{
			name: "database error",
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
					WithArgs(userID).
					WillReturnError(errors.New("db error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			entries, err := repo.ListData(context.Background(), userID)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to list data")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, entries)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
