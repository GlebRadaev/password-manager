package api_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/GlebRadaev/password-manager/internal/data/api"
	"github.com/GlebRadaev/password-manager/internal/data/models"
	"github.com/GlebRadaev/password-manager/internal/data/service"
	"github.com/GlebRadaev/password-manager/pkg/data"
)

func setupMocks(t *testing.T) (context.Context, *api.API, *api.MockService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	srv := api.NewMockService(ctrl)
	api := api.New(srv)

	return context.Background(), api, srv
}

func toModelMetadata(protoMetadata []*data.Metadata) []models.Metadata {
	var metadata []models.Metadata
	for _, m := range protoMetadata {
		metadata = append(metadata, models.Metadata{
			Key:   m.Key,
			Value: m.Value,
		})
	}
	return metadata
}

func toProtoMetadata(metadata []models.Metadata) []*data.Metadata {
	var protoMetadata []*data.Metadata
	for _, m := range metadata {
		protoMetadata = append(protoMetadata, &data.Metadata{
			Key:   m.Key,
			Value: m.Value,
		})
	}
	return protoMetadata
}

func TestAddData(t *testing.T) {
	ctx, api, srv := setupMocks(t)

	userID := uuid.NewString()
	dataID := uuid.NewString()
	testData := []byte("test data")
	metadata := []*data.Metadata{
		{Key: "key1", Value: "value1"},
	}

	tests := []struct {
		name      string
		req       *data.AddDataRequest
		setupMock func()
		want      *data.AddDataResponse
		wantErr   bool
		errCode   codes.Code
	}{
		{
			name: "successful add",
			req: &data.AddDataRequest{
				UserId:   userID,
				Data:     testData,
				Type:     data.DataType_LOGIN_PASSWORD,
				Metadata: metadata,
			},
			setupMock: func() {
				srv.EXPECT().AddData(ctx, gomock.AssignableToTypeOf(models.DataEntry{})).DoAndReturn(
					func(_ context.Context, entry models.DataEntry) (string, error) {
						assert.Equal(t, userID, entry.UserID)
						assert.Equal(t, models.DataType(data.DataType_LOGIN_PASSWORD), entry.Type)
						assert.Equal(t, testData, entry.Data)
						assert.Equal(t, toModelMetadata(metadata), entry.Metadata)
						return dataID, nil
					})
			},
			want: &data.AddDataResponse{
				DataId: dataID,
			},
			wantErr: false,
		},
		{
			name: "validation error",
			req: &data.AddDataRequest{
				UserId: "invalid-uuid",
				Data:   testData,
				Type:   data.DataType_LOGIN_PASSWORD,
			},
			setupMock: func() {},
			want:      nil,
			wantErr:   true,
			errCode:   codes.InvalidArgument,
		},
		{
			name: "service error",
			req: &data.AddDataRequest{
				UserId: userID,
				Data:   testData,
				Type:   data.DataType_LOGIN_PASSWORD,
			},
			setupMock: func() {
				srv.EXPECT().AddData(ctx, gomock.Any()).Return("", errors.New("service error"))
			},
			want:    nil,
			wantErr: true,
			errCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := api.AddData(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errCode, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, resp)
			}
		})
	}
}

func TestUpdateData(t *testing.T) {
	ctx, api, srv := setupMocks(t)

	userID := uuid.NewString()
	dataID := uuid.NewString()
	testData := []byte("updated data")
	metadata := []*data.Metadata{
		{Key: "key1", Value: "value1"},
	}

	tests := []struct {
		name      string
		req       *data.UpdateDataRequest
		setupMock func()
		want      *data.UpdateDataResponse
		wantErr   bool
		errCode   codes.Code
	}{
		{
			name: "successful update",
			req: &data.UpdateDataRequest{
				UserId:   userID,
				DataId:   dataID,
				Data:     testData,
				Metadata: metadata,
			},
			setupMock: func() {
				srv.EXPECT().UpdateData(ctx, gomock.AssignableToTypeOf(models.DataEntry{})).DoAndReturn(
					func(_ context.Context, entry models.DataEntry) error {
						assert.Equal(t, userID, entry.UserID)
						assert.Equal(t, dataID, entry.ID)
						assert.Equal(t, testData, entry.Data)
						assert.Equal(t, toModelMetadata(metadata), entry.Metadata)
						return nil
					})
			},
			want: &data.UpdateDataResponse{
				Message: "Data updated successfully",
			},
			wantErr: false,
		},
		{
			name: "validation error",
			req: &data.UpdateDataRequest{
				UserId: "invalid-uuid",
				DataId: dataID,
				Data:   testData,
			},
			setupMock: func() {},
			want:      nil,
			wantErr:   true,
			errCode:   codes.InvalidArgument,
		},
		{
			name: "data not found",
			req: &data.UpdateDataRequest{
				UserId: userID,
				DataId: dataID,
				Data:   testData,
			},
			setupMock: func() {
				srv.EXPECT().UpdateData(ctx, gomock.Any()).Return(service.ErrDataNotFound)
			},
			want:    nil,
			wantErr: true,
			errCode: codes.NotFound,
		},
		{
			name: "service error",
			req: &data.UpdateDataRequest{
				UserId: userID,
				DataId: dataID,
				Data:   testData,
			},
			setupMock: func() {
				srv.EXPECT().UpdateData(ctx, gomock.Any()).Return(errors.New("service error"))
			},
			want:    nil,
			wantErr: true,
			errCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := api.UpdateData(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errCode, status.Code(err))
			} else {
				assert.Equal(t, tt.want, resp)
			}
		})
	}
}

func TestDeleteData(t *testing.T) {
	ctx, api, srv := setupMocks(t)

	userID := uuid.NewString()
	dataID := uuid.NewString()

	tests := []struct {
		name      string
		req       *data.DeleteDataRequest
		setupMock func()
		want      *data.DeleteDataResponse
		wantErr   bool
		errCode   codes.Code
		errMsg    string
	}{
		{
			name: "successful delete",
			req: &data.DeleteDataRequest{
				UserId: userID,
				DataId: dataID,
			},
			setupMock: func() {
				srv.EXPECT().DeleteData(ctx, userID, dataID).Return(nil)
			},
			want: &data.DeleteDataResponse{
				Message: "Data deleted successfully",
			},
			wantErr: false,
		},
		{
			name: "validation error - invalid user_id",
			req: &data.DeleteDataRequest{
				UserId: "invalid-uuid",
				DataId: dataID,
			},
			setupMock: func() {},
			want:      nil,
			wantErr:   true,
			errCode:   codes.InvalidArgument,
			errMsg:    "user_id must be a valid UUID",
		},
		{
			name: "validation error - invalid data_id",
			req: &data.DeleteDataRequest{
				UserId: userID,
				DataId: "invalid-uuid",
			},
			setupMock: func() {},
			want:      nil,
			wantErr:   true,
			errCode:   codes.InvalidArgument,
			errMsg:    "data_id must be a valid UUID",
		},
		{
			name: "data not found",
			req: &data.DeleteDataRequest{
				UserId: userID,
				DataId: dataID,
			},
			setupMock: func() {
				srv.EXPECT().DeleteData(ctx, userID, dataID).Return(service.ErrDataNotFound)
			},
			want:    nil,
			wantErr: true,
			errCode: codes.NotFound,
			errMsg:  "failed to delete data: data not found",
		},
		{
			name: "service error",
			req: &data.DeleteDataRequest{
				UserId: userID,
				DataId: dataID,
			},
			setupMock: func() {
				srv.EXPECT().DeleteData(ctx, userID, dataID).Return(errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
			errCode: codes.Internal,
			errMsg:  "failed to delete data: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := api.DeleteData(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errCode, status.Code(err))
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, resp)
			}
		})
	}
}

func TestListData(t *testing.T) {
	ctx, api, srv := setupMocks(t)

	userID := uuid.NewString()
	now := time.Now()
	testEntries := []models.DataEntry{
		{
			ID:        uuid.NewString(),
			UserID:    userID,
			Type:      models.DataType(models.LoginPassword),
			Data:      []byte("data1"),
			Metadata:  []models.Metadata{{Key: "key1", Value: "value1"}},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name      string
		req       *data.ListDataRequest
		setupMock func()
		want      *data.ListDataResponse
		wantErr   bool
		errCode   codes.Code
	}{
		{
			name: "successful list",
			req: &data.ListDataRequest{
				UserId: userID,
			},
			setupMock: func() {
				srv.EXPECT().ListData(ctx, userID).Return(testEntries, nil)
			},
			want: &data.ListDataResponse{
				Entries: []*data.DataEntry{
					{
						DataId:    testEntries[0].ID,
						Type:      data.DataType_LOGIN_PASSWORD,
						Data:      testEntries[0].Data,
						Metadata:  toProtoMetadata(testEntries[0].Metadata),
						CreatedAt: now.Unix(),
						UpdatedAt: now.Unix(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "validation error",
			req: &data.ListDataRequest{
				UserId: "invalid-uuid",
			},
			setupMock: func() {},
			want:      nil,
			wantErr:   true,
			errCode:   codes.InvalidArgument,
		},
		{
			name: "service error",
			req: &data.ListDataRequest{
				UserId: userID,
			},
			setupMock: func() {
				srv.EXPECT().ListData(ctx, userID).Return(nil, errors.New("service error"))
			},
			want:    nil,
			wantErr: true,
			errCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := api.ListData(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errCode, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, resp)
			}
		})
	}
}

func TestBatchProcess(t *testing.T) {
	ctx, api, srv := setupMocks(t)

	userID := uuid.NewString()
	dataID := uuid.NewString()
	testData := []byte("test data")

	tests := []struct {
		name      string
		req       *data.BatchProcessRequest
		setupMock func()
		want      *data.BatchProcessResponse
		wantErr   bool
		errCode   codes.Code
	}{
		{
			name: "successful batch",
			req: &data.BatchProcessRequest{
				UserId: userID,
				Operations: []*data.DataOperation{
					{
						Operation: &data.DataOperation_Add{
							Add: &data.AddDataRequest{
								UserId: userID,
								Data:   testData,
								Type:   data.DataType_LOGIN_PASSWORD,
							},
						},
					},
				},
			},
			setupMock: func() {
				srv.EXPECT().BatchProcess(ctx, userID, gomock.Any()).Return([]*data.DataOperationResult{
					{
						Result: &data.DataOperationResult_Add{
							Add: &data.AddDataResponse{
								DataId: dataID,
							},
						},
					},
				}, nil)
			},
			want: &data.BatchProcessResponse{
				Results: []*data.DataOperationResult{
					{
						Result: &data.DataOperationResult_Add{
							Add: &data.AddDataResponse{
								DataId: dataID,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "validation error",
			req: &data.BatchProcessRequest{
				UserId:     "invalid-uuid",
				Operations: []*data.DataOperation{},
			},
			setupMock: func() {},
			want:      nil,
			wantErr:   true,
			errCode:   codes.InvalidArgument,
		},
		{
			name: "service error",
			req: &data.BatchProcessRequest{
				UserId: userID,
				Operations: []*data.DataOperation{
					{
						Operation: &data.DataOperation_Add{
							Add: &data.AddDataRequest{
								UserId: userID,
								Data:   testData,
								Type:   data.DataType_LOGIN_PASSWORD,
							},
						},
					},
				},
			},
			setupMock: func() {
				srv.EXPECT().BatchProcess(ctx, userID, gomock.Any()).Return(nil, errors.New("service error"))
			},
			want:    nil,
			wantErr: true,
			errCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := api.BatchProcess(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errCode, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, resp)
			}
		})
	}
}
