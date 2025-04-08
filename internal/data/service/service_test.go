package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/GlebRadaev/password-manager/internal/common/pg"
	"github.com/GlebRadaev/password-manager/internal/data/models"
	"github.com/GlebRadaev/password-manager/internal/data/repo"
	"github.com/GlebRadaev/password-manager/internal/data/service"
	"github.com/GlebRadaev/password-manager/pkg/data"
)

func setupMocks(t *testing.T) (context.Context, *service.Service, *service.MockRepo, *service.MockTxManager) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repoMock := service.NewMockRepo(ctrl)
	txManagerMock := service.NewMockTxManager(ctrl)
	srv := service.New(repoMock, txManagerMock)

	return context.Background(), srv, repoMock, txManagerMock
}

func TestAddData(t *testing.T) {
	ctx, srv, repoMock, _ := setupMocks(t)

	userID := uuid.NewString()
	dataID := uuid.NewString()
	testEntry := models.DataEntry{
		UserID:   userID,
		Type:     models.LoginPassword,
		Data:     []byte("test data"),
		Metadata: []models.Metadata{{Key: "key", Value: "value"}},
	}

	tests := []struct {
		name      string
		entry     models.DataEntry
		setupMock func()
		want      string
		wantErr   error
	}{
		{
			name:  "successful add",
			entry: testEntry,
			setupMock: func() {
				repoMock.EXPECT().AddList(ctx, []models.DataEntry{testEntry}).
					Return([]string{dataID}, nil)
			},
			want:    dataID,
			wantErr: nil,
		},
		{
			name:  "empty result",
			entry: testEntry,
			setupMock: func() {
				repoMock.EXPECT().AddList(ctx, []models.DataEntry{testEntry}).
					Return([]string{}, nil)
			},
			want:    "",
			wantErr: errors.New("no data was added"),
		},
		{
			name:  "repository error",
			entry: testEntry,
			setupMock: func() {
				repoMock.EXPECT().AddList(ctx, []models.DataEntry{testEntry}).
					Return(nil, errors.New("db error"))
			},
			want:    "",
			wantErr: errors.New("failed to add data: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			got, err := srv.AddData(ctx, tt.entry)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUpdateData(t *testing.T) {
	ctx, srv, repoMock, _ := setupMocks(t)

	userID := uuid.NewString()
	dataID := uuid.NewString()
	testEntry := models.DataEntry{
		UserID: userID,
		ID:     dataID,
		Data:   []byte("updated data"),
	}

	tests := []struct {
		name      string
		entry     models.DataEntry
		setupMock func()
		wantErr   error
	}{
		{
			name:  "successful update",
			entry: testEntry,
			setupMock: func() {
				repoMock.EXPECT().UpdateData(ctx, testEntry).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:  "data not found",
			entry: testEntry,
			setupMock: func() {
				repoMock.EXPECT().UpdateData(ctx, testEntry).
					Return(repo.ErrDataNotFound)
			},
			wantErr: service.ErrDataNotFound,
		},
		{
			name:  "repository error",
			entry: testEntry,
			setupMock: func() {
				repoMock.EXPECT().UpdateData(ctx, testEntry).
					Return(errors.New("db error"))
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := srv.UpdateData(ctx, tt.entry)

			if tt.wantErr != nil {
				assert.Error(t, err)
				if errors.Is(tt.wantErr, service.ErrDataNotFound) {
					assert.ErrorIs(t, err, service.ErrDataNotFound)
				} else {
					assert.EqualError(t, err, tt.wantErr.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteData(t *testing.T) {
	ctx, srv, repoMock, _ := setupMocks(t)

	userID := uuid.NewString()
	dataID := uuid.NewString()

	tests := []struct {
		name      string
		userID    string
		dataID    string
		setupMock func()
		wantErr   error
	}{
		{
			name:   "successful delete",
			userID: userID,
			dataID: dataID,
			setupMock: func() {
				repoMock.EXPECT().DeleteList(ctx, userID, []string{dataID}).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:   "data not found",
			userID: userID,
			dataID: dataID,
			setupMock: func() {
				repoMock.EXPECT().DeleteList(ctx, userID, []string{dataID}).
					Return(repo.ErrDataNotFound)
			},
			wantErr: service.ErrDataNotFound,
		},
		{
			name:   "repository error",
			userID: userID,
			dataID: dataID,
			setupMock: func() {
				repoMock.EXPECT().DeleteList(ctx, userID, []string{dataID}).
					Return(errors.New("db error"))
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := srv.DeleteData(ctx, tt.userID, tt.dataID)

			if tt.wantErr != nil {
				assert.Error(t, err)
				if errors.Is(tt.wantErr, service.ErrDataNotFound) {
					assert.ErrorIs(t, err, service.ErrDataNotFound)
				} else {
					assert.EqualError(t, err, tt.wantErr.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListData(t *testing.T) {
	ctx, srv, repoMock, _ := setupMocks(t)

	userID := uuid.NewString()
	testEntries := []models.DataEntry{
		{
			ID:     uuid.NewString(),
			UserID: userID,
			Type:   models.LoginPassword,
			Data:   []byte("test data"),
		},
	}

	tests := []struct {
		name      string
		userID    string
		setupMock func()
		want      []models.DataEntry
		wantErr   error
	}{
		{
			name:   "successful list",
			userID: userID,
			setupMock: func() {
				repoMock.EXPECT().ListData(ctx, userID).
					Return(testEntries, nil)
			},
			want:    testEntries,
			wantErr: nil,
		},
		{
			name:   "repository error",
			userID: userID,
			setupMock: func() {
				repoMock.EXPECT().ListData(ctx, userID).
					Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			got, err := srv.ListData(ctx, tt.userID)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBatchProcess(t *testing.T) {
	ctx, srv, repoMock, txManagerMock := setupMocks(t)

	userID := uuid.NewString()
	addDataID := uuid.NewString()
	updateDataID := uuid.NewString()
	deleteDataID := uuid.NewString()

	operations := []*data.DataOperation{
		{
			Operation: &data.DataOperation_Add{
				Add: &data.AddDataRequest{
					Type: data.DataType_LOGIN_PASSWORD,
					Data: []byte("add data"),
				},
			},
		},
		{
			Operation: &data.DataOperation_Update{
				Update: &data.UpdateDataRequest{
					DataId: updateDataID,
					Data:   []byte("updated data"),
				},
			},
		},
		{
			Operation: &data.DataOperation_Delete{
				Delete: &data.DeleteDataRequest{
					DataId: deleteDataID,
				},
			},
		},
	}

	tests := []struct {
		name      string
		userID    string
		ops       []*data.DataOperation
		setupMock func()
		want      []*data.DataOperationResult
		wantErr   error
	}{
		{
			name:   "successful batch",
			userID: userID,
			ops:    operations,
			setupMock: func() {
				txManagerMock.EXPECT().Begin(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn pg.TransactionalFn) error {
						return fn(ctx)
					})

				// Add expectations
				repoMock.EXPECT().AddList(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, entries []models.DataEntry) ([]string, error) {
						assert.Equal(t, 1, len(entries))
						assert.Equal(t, userID, entries[0].UserID)
						assert.Equal(t, models.LoginPassword, entries[0].Type)
						return []string{addDataID}, nil
					})

				// Update expectations
				repoMock.EXPECT().UpdateData(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, entry models.DataEntry) error {
						assert.Equal(t, userID, entry.UserID)
						assert.Equal(t, updateDataID, entry.ID)
						return nil
					})

				// Delete expectations
				repoMock.EXPECT().DeleteList(ctx, userID, []string{deleteDataID}).
					Return(nil)
			},
			want: []*data.DataOperationResult{
				{
					Result: &data.DataOperationResult_Add{
						Add: &data.AddDataResponse{DataId: addDataID},
					},
				},
				{
					Result: &data.DataOperationResult_Update{
						Update: &data.UpdateDataResponse{Message: "Data updated successfully"},
					},
				},
				{
					Result: &data.DataOperationResult_Delete{
						Delete: &data.DeleteDataResponse{Message: "Data deleted successfully"},
					},
				},
			},
			wantErr: nil,
		},
		{
			name:   "transaction error",
			userID: userID,
			ops:    operations,
			setupMock: func() {
				txManagerMock.EXPECT().Begin(ctx, gomock.Any()).
					Return(errors.New("tx error"))
			},
			want:    nil,
			wantErr: errors.New("failed to process batch operations: tx error"),
		},
		{
			name:   "add error in batch",
			userID: userID,
			ops:    operations,
			setupMock: func() {
				txManagerMock.EXPECT().Begin(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn pg.TransactionalFn) error {
						return fn(ctx)
					})

				repoMock.EXPECT().AddList(ctx, gomock.Any()).
					Return(nil, errors.New("add error"))
			},
			want:    nil,
			wantErr: errors.New("failed to process batch operations: failed to add data: add error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			got, err := srv.BatchProcess(ctx, tt.userID, tt.ops)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.want), len(got))
				for i, res := range got {
					switch r := res.Result.(type) {
					case *data.DataOperationResult_Add:
						assert.Equal(t, tt.want[i].GetAdd().DataId, r.Add.DataId)
					case *data.DataOperationResult_Update:
						assert.Equal(t, tt.want[i].GetUpdate().Message, r.Update.Message)
					case *data.DataOperationResult_Delete:
						assert.Equal(t, tt.want[i].GetDelete().Message, r.Delete.Message)
					}
				}
			}
		})
	}
}
