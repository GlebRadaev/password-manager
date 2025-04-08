package service_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/GlebRadaev/password-manager/internal/sync/models"
	"github.com/GlebRadaev/password-manager/internal/sync/service"
	"github.com/GlebRadaev/password-manager/pkg/data"
)

func setupMocks(t *testing.T) (context.Context, *service.Service, *service.MockRepo, *service.MockDataClient) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repoMock := service.NewMockRepo(ctrl)
	dataClientMock := service.NewMockDataClient(ctrl)
	srv := service.New(repoMock, dataClientMock)

	return context.Background(), srv, repoMock, dataClientMock
}

func TestSyncData(t *testing.T) {
	ctx, srv, repoMock, dataClientMock := setupMocks(t)

	userID := uuid.NewString()
	now := time.Now()
	clientData := []models.ClientData{
		{
			DataID:    "data1",
			Operation: models.Add,
			Data:      []byte("client data 1"),
			UpdatedAt: now,
		},
		{
			DataID:    "data2",
			Operation: models.Update,
			Data:      []byte("client data 2"),
			UpdatedAt: now.Add(time.Hour),
		},
		{
			DataID:    "data3",
			Operation: models.Delete,
		},
	}

	serverEntries := []models.DataEntry{
		{
			DataID:    "data1",
			Data:      []byte("server data 1"),
			UpdatedAt: now.Add(-time.Hour),
		},
		{
			DataID:    "data2",
			Data:      []byte("server data 2"),
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name      string
		setupMock func()
		want      []models.Conflict
		wantErr   error
	}{
		{
			name: "successful sync with conflicts",
			setupMock: func() {
				dataClientMock.EXPECT().ListData(ctx, userID).
					Return(serverEntries, nil).Times(1)

				dataClientMock.EXPECT().CreateUpdateOperation(userID, gomock.Any()).
					DoAndReturn(func(userID string, entry models.ClientData) *data.DataOperation {
						assert.Equal(t, "data1", entry.DataID)
						return &data.DataOperation{
							Operation: &data.DataOperation_Update{
								Update: &data.UpdateDataRequest{
									DataId: "data1",
									Data:   entry.Data,
								},
							},
						}
					}).Times(1)

				dataClientMock.EXPECT().CreateUpdateOperation(userID, gomock.Any()).
					DoAndReturn(func(userID string, entry models.ClientData) *data.DataOperation {
						assert.Equal(t, "data2", entry.DataID)
						return &data.DataOperation{
							Operation: &data.DataOperation_Update{
								Update: &data.UpdateDataRequest{
									DataId: "data2",
									Data:   entry.Data,
								},
							},
						}
					}).Times(1)

				dataClientMock.EXPECT().BatchProcess(ctx, userID, gomock.Any()).
					DoAndReturn(func(ctx context.Context, userID string, ops []*data.DataOperation) (string, error) {
						if len(ops) != 2 {
							return "", fmt.Errorf("expected 2 operations, got %d", len(ops))
						}

						var hasUpdate1, hasUpdate2 bool
						for _, op := range ops {
							switch v := op.Operation.(type) {
							case *data.DataOperation_Update:
								if v.Update.DataId == "data1" {
									hasUpdate1 = true
								}
								if v.Update.DataId == "data2" {
									hasUpdate2 = true
								}
							}
						}

						if !hasUpdate1 || !hasUpdate2 {
							return "", errors.New("missing expected operations")
						}
						return "batch-id", nil
					}).Times(1)
			},
			want:    nil,
			wantErr: nil,
		},

		{
			name: "failed to list server data",
			setupMock: func() {
				dataClientMock.EXPECT().ListData(ctx, userID).
					Return(nil, errors.New("list error")).Times(1)
			},
			want:    nil,
			wantErr: errors.New("failed to list server data: list error"),
		},

		{
			name: "failed to add conflicts",
			setupMock: func() {
				conflictServerEntries := []models.DataEntry{
					{
						DataID:    "data1",
						Data:      []byte("server data 1"),
						UpdatedAt: now.Add(-time.Hour),
					},
					{
						DataID:    "data2",
						Data:      []byte("server data 2"),
						UpdatedAt: now.Add(2 * time.Hour),
					},
					{
						DataID:    "data3",
						Data:      []byte("server data 3"),
						UpdatedAt: now.Add(-time.Hour),
					},
				}

				dataClientMock.EXPECT().ListData(ctx, userID).
					Return(conflictServerEntries, nil).Times(1)

				dataClientMock.EXPECT().CreateUpdateOperation(userID, gomock.Any()).
					DoAndReturn(func(userID string, entry models.ClientData) *data.DataOperation {
						assert.Equal(t, "data1", entry.DataID)
						return &data.DataOperation{
							Operation: &data.DataOperation_Update{
								Update: &data.UpdateDataRequest{
									DataId: "data1",
									Data:   entry.Data,
								},
							},
						}
					}).Times(1)

				dataClientMock.EXPECT().CreateDeleteOperation(userID, "data3").
					Return(&data.DataOperation{
						Operation: &data.DataOperation_Delete{
							Delete: &data.DeleteDataRequest{
								DataId: "data3",
							},
						},
					}).Times(1)

				repoMock.EXPECT().AddConflicts(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, conflicts []models.Conflict) error {
						if len(conflicts) != 1 {
							return fmt.Errorf("expected 1 conflict, got %d", len(conflicts))
						}
						assert.Equal(t, userID, conflicts[0].UserID)
						assert.Equal(t, "data2", conflicts[0].DataID)
						assert.True(t, bytes.Equal([]byte("client data 2"), conflicts[0].ClientData))
						assert.True(t, bytes.Equal([]byte("server data 2"), conflicts[0].ServerData))
						return errors.New("add conflict error")
					}).Times(1)
			},
			want:    nil,
			wantErr: errors.New("failed to add conflicts: add conflict error"),
		},
		{
			name: "failed batch process with rollback (no conflicts)",
			setupMock: func() {
				dataClientMock.EXPECT().ListData(ctx, userID).
					Return(serverEntries, nil).Times(1)

				dataClientMock.EXPECT().CreateUpdateOperation(userID, gomock.Any()).
					Return(&data.DataOperation{
						Operation: &data.DataOperation_Update{
							Update: &data.UpdateDataRequest{
								DataId: "data1",
								Data:   []byte("client data 1"),
							},
						},
					}).Times(1)

				dataClientMock.EXPECT().CreateUpdateOperation(userID, gomock.Any()).
					Return(&data.DataOperation{
						Operation: &data.DataOperation_Update{
							Update: &data.UpdateDataRequest{
								DataId: "data2",
								Data:   []byte("client data 2"),
							},
						},
					}).Times(1)

				dataClientMock.EXPECT().BatchProcess(ctx, userID, gomock.Any()).
					Return("", errors.New("batch error")).Times(1)
			},
			want:    nil,
			wantErr: errors.New("failed to process batch operations: batch error"),
		},
		{
			name: "failed batch process with rollback error",
			setupMock: func() {
				conflictServerEntries := []models.DataEntry{
					{
						DataID:    "data1",
						Data:      []byte("server data 1"),
						UpdatedAt: now.Add(2 * time.Hour),
					},
					{
						DataID:    "data2",
						Data:      []byte("server data 2"),
						UpdatedAt: now.Add(-time.Hour),
					},
				}

				dataClientMock.EXPECT().ListData(ctx, userID).
					Return(conflictServerEntries, nil).Times(1)

				dataClientMock.EXPECT().CreateUpdateOperation(userID, gomock.Any()).
					DoAndReturn(func(userID string, entry models.ClientData) *data.DataOperation {
						assert.Equal(t, "data2", entry.DataID)
						return &data.DataOperation{
							Operation: &data.DataOperation_Update{
								Update: &data.UpdateDataRequest{
									DataId: "data2",
									Data:   entry.Data,
								},
							},
						}
					}).Times(1)

				repoMock.EXPECT().AddConflicts(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, conflicts []models.Conflict) error {
						if len(conflicts) != 1 {
							return fmt.Errorf("expected 1 conflict, got %d", len(conflicts))
						}
						assert.Equal(t, "data1", conflicts[0].DataID)
						return nil
					}).Times(1)

				dataClientMock.EXPECT().BatchProcess(ctx, userID, gomock.Any()).
					DoAndReturn(func(ctx context.Context, userID string, ops []*data.DataOperation) (string, error) {
						if len(ops) != 1 {
							return "", fmt.Errorf("expected 1 operation, got %d", len(ops))
						}
						return "", errors.New("batch error")
					}).Times(1)

				repoMock.EXPECT().DeleteConflicts(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, conflictIDs []string) error {
						if len(conflictIDs) != 1 {
							return fmt.Errorf("expected 1 conflict ID, got %d", len(conflictIDs))
						}
						return errors.New("rollback error")
					}).Times(1)
			},
			wantErr: errors.New("failed to process batch operations and rollback conflicts: batch error (rollback error: failed to rollback inserted conflicts: rollback error)"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			got, err := srv.SyncData(ctx, userID, clientData)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				if len(tt.want) > 0 {
					assert.Equal(t, tt.want[0].UserID, got[0].UserID)
					assert.Equal(t, tt.want[0].DataID, got[0].DataID)
					assert.True(t, bytes.Equal(tt.want[0].ClientData, got[0].ClientData))
					assert.True(t, bytes.Equal(tt.want[0].ServerData, got[0].ServerData))
				} else {
					assert.Empty(t, got)
				}
			}
		})
	}
}

func TestResolveConflict(t *testing.T) {
	ctx, srv, repoMock, dataClientMock := setupMocks(t)

	conflictID := uuid.NewString()
	userID := uuid.NewString()
	dataID := uuid.NewString()
	conflict := models.Conflict{
		ID:         conflictID,
		UserID:     userID,
		DataID:     dataID,
		ClientData: []byte("client data"),
		ServerData: []byte("server data"),
	}

	tests := []struct {
		name      string
		strategy  models.ResolutionStrategy
		setupMock func()
		wantErr   error
	}{
		{
			name:     "successful resolve with client version",
			strategy: models.UseClientVersion,
			setupMock: func() {
				repoMock.EXPECT().GetConflict(ctx, conflictID).
					Return(conflict, nil)

				dataClientMock.EXPECT().UpdateData(ctx, userID, models.ClientData{
					DataID: dataID,
					Data:   conflict.ClientData,
				}).Return(nil)

				repoMock.EXPECT().ResolveConflict(ctx, conflictID).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:     "successful resolve with server version",
			strategy: models.UseServerVersion,
			setupMock: func() {
				repoMock.EXPECT().GetConflict(ctx, conflictID).
					Return(conflict, nil)

				dataClientMock.EXPECT().UpdateData(ctx, userID, models.ClientData{
					DataID: dataID,
					Data:   conflict.ServerData,
				}).Return(nil)

				repoMock.EXPECT().ResolveConflict(ctx, conflictID).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:     "conflict not found",
			strategy: models.UseClientVersion,
			setupMock: func() {
				repoMock.EXPECT().GetConflict(ctx, conflictID).
					Return(models.Conflict{}, service.ErrConflictNotFound)
			},
			wantErr: service.ErrConflictNotFound,
		},
		{
			name:     "failed to update data",
			strategy: models.UseClientVersion,
			setupMock: func() {
				repoMock.EXPECT().GetConflict(ctx, conflictID).
					Return(conflict, nil)

				dataClientMock.EXPECT().UpdateData(ctx, userID, gomock.Any()).
					Return(errors.New("update error"))
			},
			wantErr: errors.New("failed to update data: update error"),
		},
		{
			name:     "failed to mark as resolved",
			strategy: models.UseClientVersion,
			setupMock: func() {
				repoMock.EXPECT().GetConflict(ctx, conflictID).
					Return(conflict, nil)

				dataClientMock.EXPECT().UpdateData(ctx, userID, gomock.Any()).
					Return(nil)

				repoMock.EXPECT().ResolveConflict(ctx, conflictID).
					Return(errors.New("resolve error"))
			},
			wantErr: errors.New("failed to resolve conflict: resolve error"),
		},
		{
			name:     "failed to get conflict with other error",
			strategy: models.UseClientVersion,
			setupMock: func() {
				repoMock.EXPECT().GetConflict(ctx, conflictID).
					Return(models.Conflict{}, errors.New("some db error"))
			},
			wantErr: errors.New("failed to get conflict: some db error"),
		},
		{
			name:     "unknown resolution strategy",
			strategy: models.ResolutionStrategy(999),
			setupMock: func() {
				repoMock.EXPECT().GetConflict(ctx, conflictID).
					Return(conflict, nil)
			},
			wantErr: errors.New("unknown resolution strategy"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := srv.ResolveConflict(ctx, conflictID, tt.strategy)

			if tt.wantErr != nil {
				assert.Error(t, err)
				if errors.Is(tt.wantErr, service.ErrConflictNotFound) {
					assert.ErrorIs(t, err, service.ErrConflictNotFound)
				} else {
					assert.Contains(t, err.Error(), tt.wantErr.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListConflicts(t *testing.T) {
	ctx, srv, repoMock, _ := setupMocks(t)

	userID := uuid.NewString()
	conflicts := []models.Conflict{
		{
			ID:         uuid.NewString(),
			UserID:     userID,
			DataID:     "data1",
			ClientData: []byte("client data"),
			ServerData: []byte("server data"),
		},
	}

	tests := []struct {
		name      string
		setupMock func()
		want      []models.Conflict
		wantErr   error
	}{
		{
			name: "successful list",
			setupMock: func() {
				repoMock.EXPECT().GetUnresolvedConflicts(ctx, userID).
					Return(conflicts, nil)
			},
			want:    conflicts,
			wantErr: nil,
		},
		{
			name: "repository error",
			setupMock: func() {
				repoMock.EXPECT().GetUnresolvedConflicts(ctx, userID).
					Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: errors.New("failed to get unresolved conflicts: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			got, err := srv.ListConflicts(ctx, userID)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRollbackInsertedConflicts(t *testing.T) {
	ctx, srv, repoMock, _ := setupMocks(t)

	userID := uuid.NewString()
	conflict1 := models.Conflict{
		ID:     uuid.NewString(),
		UserID: userID,
		DataID: "data1",
	}
	conflict2 := models.Conflict{
		ID:     uuid.NewString(),
		UserID: userID,
		DataID: "data2",
	}

	tests := []struct {
		name      string
		conflicts []models.Conflict
		setupMock func()
		wantErr   error
	}{
		{
			name:      "successful rollback with conflicts",
			conflicts: []models.Conflict{conflict1, conflict2},
			setupMock: func() {
				repoMock.EXPECT().DeleteConflicts(ctx, []string{conflict1.ID, conflict2.ID}).
					Return(nil).
					Times(1)
			},
			wantErr: nil,
		},
		{
			name:      "empty conflicts - no operation",
			conflicts: []models.Conflict{},
			setupMock: func() {
			},
			wantErr: nil,
		},
		{
			name:      "failed to delete conflicts",
			conflicts: []models.Conflict{conflict1},
			setupMock: func() {
				repoMock.EXPECT().DeleteConflicts(ctx, []string{conflict1.ID}).
					Return(errors.New("delete error")).
					Times(1)
			},
			wantErr: errors.New("failed to rollback inserted conflicts: delete error"),
		},
		{
			name:      "nil conflicts - no operation",
			conflicts: nil,
			setupMock: func() {
			},
			wantErr: nil,
		},
		{
			name:      "single conflict rollback",
			conflicts: []models.Conflict{conflict2},
			setupMock: func() {
				repoMock.EXPECT().DeleteConflicts(ctx, []string{conflict2.ID}).
					Return(nil).
					Times(1)
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := srv.RollbackInsertedConflicts(ctx, tt.conflicts)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMergeData(t *testing.T) {
	tests := []struct {
		name       string
		clientData []byte
		serverData []byte
		dataType   models.DataType
		want       []byte
		wantErr    error
	}{
		{
			name:       "text merge",
			clientData: []byte("client"),
			serverData: []byte("server"),
			dataType:   models.Text,
			want:       []byte("client\nserver"),
			wantErr:    nil,
		},
		{
			name:       "binary merge",
			clientData: []byte("client"),
			serverData: []byte("server"),
			dataType:   models.Binary,
			want:       []byte("client"),
			wantErr:    nil,
		},
		{
			name:       "login/password merge",
			clientData: []byte("client"),
			serverData: []byte("server"),
			dataType:   models.LoginPassword,
			want:       []byte("client"),
			wantErr:    nil,
		},
		{
			name:       "card merge",
			clientData: []byte("client"),
			serverData: []byte("server"),
			dataType:   models.Card,
			want:       []byte("client"),
			wantErr:    nil,
		},
		{
			name:       "unsupported type",
			clientData: []byte("client"),
			serverData: []byte("server"),
			dataType:   models.DataType(999),
			want:       nil,
			wantErr:    fmt.Errorf("unsupported data type for merging: %v", models.DataType(999)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.MergeData(tt.clientData, tt.serverData, tt.dataType)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
