package repo

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"

	"github.com/GlebRadaev/password-manager/internal/sync/models"
)

func NewMock(t *testing.T) (*Repository, pgxmock.PgxPoolIface) {
	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	repo := New(mockDB)

	t.Cleanup(func() {
		mockDB.Close()
	})

	return repo, mockDB
}

func TestRepository_GetConflict(t *testing.T) {
	repo, mock := NewMock(t)

	now := time.Now()
	conflictID := "conflict1"
	expectedConflict := models.Conflict{
		ID:         conflictID,
		UserID:     "user1",
		DataID:     "data1",
		ClientData: []byte("client data"),
		ServerData: []byte("server data"),
		Resolved:   false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
		expected  models.Conflict
	}{
		{
			name: "successful get",
			mockSetup: func() {
				rows := pgxmock.NewRows([]string{"id", "user_id", "data_id", "client_data", "server_data", "resolved", "created_at", "updated_at"}).
					AddRow(
						expectedConflict.ID,
						expectedConflict.UserID,
						expectedConflict.DataID,
						expectedConflict.ClientData,
						expectedConflict.ServerData,
						expectedConflict.Resolved,
						expectedConflict.CreatedAt,
						expectedConflict.UpdatedAt,
					)
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, user_id, data_id, client_data, server_data, resolved, created_at, updated_at
					FROM sync.conflicts
					WHERE id = $1`)).
					WithArgs(conflictID).
					WillReturnRows(rows)
			},
			expectErr: false,
			expected:  expectedConflict,
		},
		{
			name: "conflict not found",
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
					WithArgs(conflictID).
					WillReturnError(pgx.ErrNoRows)
			},
			expectErr: true,
		},
		{
			name: "database error",
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
					WithArgs(conflictID).
					WillReturnError(errors.New("db error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			conflict, err := repo.GetConflict(context.Background(), conflictID)

			if tt.expectErr {
				assert.Error(t, err)
				if errors.Is(err, ErrConflictNotFound) {
					assert.Contains(t, err.Error(), "conflict not found")
				} else {
					assert.Contains(t, err.Error(), "failed to get conflict")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, conflict)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_AddConflicts(t *testing.T) {
	repo, mock := NewMock(t)

	conflicts := []models.Conflict{
		{
			ID:         "conflict1",
			UserID:     "user1",
			DataID:     "data1",
			ClientData: []byte("client data 1"),
			ServerData: []byte("server data 1"),
		},
		{
			ID:         "conflict2",
			UserID:     "user1",
			DataID:     "data2",
			ClientData: []byte("client data 2"),
			ServerData: []byte("server data 2"),
		},
	}

	tests := []struct {
		name      string
		conflicts []models.Conflict
		mockSetup func([]models.Conflict)
		expectErr bool
	}{
		{
			name:      "successful batch insert",
			conflicts: conflicts,
			mockSetup: func(confs []models.Conflict) {
				batch := mock.ExpectBatch()
				for _, conflict := range confs {
					batch.ExpectExec(regexp.QuoteMeta(`
                        INSERT INTO sync.conflicts (id, user_id, data_id, client_data, server_data, resolved, created_at, updated_at)
                        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)).
						WithArgs(
							conflict.ID,
							conflict.UserID,
							conflict.DataID,
							conflict.ClientData,
							conflict.ServerData,
							false,
							pgxmock.AnyArg(),
							pgxmock.AnyArg(),
						).
						WillReturnResult(pgxmock.NewResult("INSERT", 1))
				}
			},
			expectErr: false,
		},
		{
			name:      "empty conflicts",
			conflicts: []models.Conflict{},
			mockSetup: func(confs []models.Conflict) {
			},
			expectErr: false,
		},
		{
			name:      "database error",
			conflicts: conflicts[:1],
			mockSetup: func(confs []models.Conflict) {
				batch := mock.ExpectBatch()
				batch.ExpectExec(regexp.QuoteMeta(`INSERT`)).
					WithArgs(
						confs[0].ID,
						confs[0].UserID,
						confs[0].DataID,
						confs[0].ClientData,
						confs[0].ServerData,
						false,
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
			tt.mockSetup(tt.conflicts)

			err := repo.AddConflicts(context.Background(), tt.conflicts)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to add conflict")
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetUnresolvedConflicts(t *testing.T) {
	repo, mock := NewMock(t)

	now := time.Now()
	userID := "user1"
	expectedConflicts := []models.Conflict{
		{
			ID:         "conflict1",
			UserID:     userID,
			DataID:     "data1",
			ClientData: []byte("client data 1"),
			ServerData: []byte("server data 1"),
			Resolved:   false,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "conflict2",
			UserID:     userID,
			DataID:     "data2",
			ClientData: []byte("client data 2"),
			ServerData: []byte("server data 2"),
			Resolved:   false,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
		expected  []models.Conflict
	}{
		{
			name: "successful get unresolved conflicts",
			mockSetup: func() {
				rows := pgxmock.NewRows([]string{"id", "user_id", "data_id", "client_data", "server_data", "resolved", "created_at", "updated_at"}).
					AddRow(
						expectedConflicts[0].ID,
						expectedConflicts[0].UserID,
						expectedConflicts[0].DataID,
						expectedConflicts[0].ClientData,
						expectedConflicts[0].ServerData,
						expectedConflicts[0].Resolved,
						expectedConflicts[0].CreatedAt,
						expectedConflicts[0].UpdatedAt,
					).
					AddRow(
						expectedConflicts[1].ID,
						expectedConflicts[1].UserID,
						expectedConflicts[1].DataID,
						expectedConflicts[1].ClientData,
						expectedConflicts[1].ServerData,
						expectedConflicts[1].Resolved,
						expectedConflicts[1].CreatedAt,
						expectedConflicts[1].UpdatedAt,
					)
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, user_id, data_id, client_data, server_data, resolved, created_at, updated_at
					FROM sync.conflicts
					WHERE user_id = $1 AND resolved = FALSE`)).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectErr: false,
			expected:  expectedConflicts,
		},
		{
			name: "no unresolved conflicts",
			mockSetup: func() {
				rows := pgxmock.NewRows([]string{"id", "user_id", "data_id", "client_data", "server_data", "resolved", "created_at", "updated_at"})
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

			conflicts, err := repo.GetUnresolvedConflicts(context.Background(), userID)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get unresolved conflicts")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, conflicts)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_ResolveConflict(t *testing.T) {
	repo, mock := NewMock(t)

	conflictID := "conflict1"

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
	}{
		{
			name: "successful resolve",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`
					UPDATE sync.conflicts
					SET resolved = TRUE, updated_at = CURRENT_TIMESTAMP
					WHERE id = $1`)).
					WithArgs(conflictID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			expectErr: false,
		},
		{
			name: "conflict not found",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE sync.conflicts`)).
					WithArgs(conflictID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))
			},
			expectErr: false, // Not considered an error if conflict not found
		},
		{
			name: "database error",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE sync.conflicts`)).
					WithArgs(conflictID).
					WillReturnError(errors.New("db error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.ResolveConflict(context.Background(), conflictID)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to resolve conflict")
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_DeleteConflicts(t *testing.T) {
	repo, mock := NewMock(t)

	conflictIDs := []string{"conflict1", "conflict2"}

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
	}{
		{
			name: "successful delete",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`
					DELETE FROM sync.conflicts
					WHERE id = ANY($1)`)).
					WithArgs(conflictIDs).
					WillReturnResult(pgxmock.NewResult("DELETE", 2))
			},
			expectErr: false,
		},
		{
			name: "no conflicts deleted",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM sync.conflicts`)).
					WithArgs(conflictIDs).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))
			},
			expectErr: false, // Not considered an error if no conflicts deleted
		},
		{
			name: "database error",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM sync.conflicts`)).
					WithArgs(conflictIDs).
					WillReturnError(errors.New("db error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.DeleteConflicts(context.Background(), conflictIDs)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete conflicts")
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
