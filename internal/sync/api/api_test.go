package api

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

	"github.com/GlebRadaev/password-manager/internal/sync/models"
	"github.com/GlebRadaev/password-manager/internal/sync/service"
	"github.com/GlebRadaev/password-manager/pkg/sync"
)

func setupMocks(t *testing.T) (context.Context, *API, *MockService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	srv := NewMockService(ctrl)
	api := New(srv)

	return context.Background(), api, srv
}

func TestSyncData(t *testing.T) {
	ctx, api, srv := setupMocks(t)

	userID := uuid.NewString()
	now := time.Now()
	clientData := []*sync.ClientData{
		{
			DataId:    "data1",
			Type:      sync.DataType_LOGIN_PASSWORD,
			Data:      []byte("testdata"),
			Metadata:  []*sync.Metadata{{Key: "key1", Value: "value1"}},
			Operation: sync.Operation_UPDATE,
			UpdatedAt: now.Unix(),
		},
	}

	conflicts := []models.Conflict{
		{
			ID:         "conflict1",
			DataID:     "data1",
			ClientData: []byte("clientdata"),
			ServerData: []byte("serverdata"),
			Resolved:   false,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	tests := []struct {
		name      string
		req       *sync.SyncDataRequest
		setupMock func()
		want      *sync.SyncDataResponse
		wantErr   bool
		errCode   codes.Code
	}{
		{
			name: "successful sync",
			req: &sync.SyncDataRequest{
				UserId:     userID,
				ClientData: clientData,
			},
			setupMock: func() {
				srv.EXPECT().SyncData(ctx, userID, gomock.Any()).Return(conflicts, nil)
			},
			want: &sync.SyncDataResponse{
				Conflicts: []*sync.Conflict{
					{
						ConflictId: "conflict1",
						DataId:     "data1",
						ClientData: []byte("clientdata"),
						ServerData: []byte("serverdata"),
						Resolved:   false,
						CreatedAt:  now.Unix(),
						UpdatedAt:  now.Unix(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid request",
			req: &sync.SyncDataRequest{
				UserId: "",
			},
			setupMock: func() {},
			want:      nil,
			wantErr:   true,
			errCode:   codes.InvalidArgument,
		},
		{
			name: "internal error",
			req: &sync.SyncDataRequest{
				UserId:     userID,
				ClientData: clientData,
			},
			setupMock: func() {
				srv.EXPECT().SyncData(ctx, userID, gomock.Any()).Return(nil, errors.New("internal error"))
			},
			want:    nil,
			wantErr: true,
			errCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := api.SyncData(ctx, tt.req)

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

func TestResolveConflict(t *testing.T) {
	ctx, api, srv := setupMocks(t)

	conflictID := "a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6a7b8"
	userID := "a1b2c3d4-e5f6-7890-a1b2-c3d4e5f6a7b8"

	tests := []struct {
		name      string
		req       *sync.ResolveConflictRequest
		setupMock func()
		want      *sync.ResolveConflictResponse
		wantErr   bool
		errCode   codes.Code
	}{
		{
			name: "successful resolution",
			req: &sync.ResolveConflictRequest{
				ConflictId: conflictID,
				UserId:     userID,
				Strategy:   sync.ResolutionStrategy_USE_CLIENT_VERSION,
			},
			setupMock: func() {
				srv.EXPECT().ResolveConflict(
					ctx,
					conflictID,
					models.UseClientVersion,
				).Return(nil)
			},
			want: &sync.ResolveConflictResponse{
				Message: "Conflict resolved successfully",
			},
			wantErr: false,
		},
		{
			name: "conflict not found",
			req: &sync.ResolveConflictRequest{
				ConflictId: conflictID,
				UserId:     userID,
				Strategy:   sync.ResolutionStrategy_USE_CLIENT_VERSION,
			},
			setupMock: func() {
				srv.EXPECT().ResolveConflict(
					ctx,
					conflictID,
					models.UseClientVersion,
				).Return(service.ErrConflictNotFound)
			},
			want:    nil,
			wantErr: true,
			errCode: codes.NotFound,
		},
		{
			name: "internal error",
			req: &sync.ResolveConflictRequest{
				ConflictId: conflictID,
				UserId:     userID,
				Strategy:   sync.ResolutionStrategy_USE_CLIENT_VERSION,
			},
			setupMock: func() {
				srv.EXPECT().ResolveConflict(
					ctx,
					conflictID,
					models.UseClientVersion,
				).Return(errors.New("internal error"))
			},
			want:    nil,
			wantErr: true,
			errCode: codes.Internal,
		},
		{
			name: "missing required fields",
			req: &sync.ResolveConflictRequest{
				Strategy: sync.ResolutionStrategy_USE_CLIENT_VERSION,
			},
			setupMock: func() {
			},
			want:    nil,
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := api.ResolveConflict(ctx, tt.req)

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

func TestListConflicts(t *testing.T) {
	ctx, api, srv := setupMocks(t)

	userID := uuid.NewString()
	now := time.Now()
	conflicts := []models.Conflict{
		{
			ID:         "conflict1",
			DataID:     "data1",
			ClientData: []byte("clientdata"),
			ServerData: []byte("serverdata"),
			Resolved:   false,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	tests := []struct {
		name      string
		req       *sync.ListConflictsRequest
		setupMock func()
		want      *sync.ListConflictsResponse
		wantErr   bool
		errCode   codes.Code
	}{
		{
			name: "successful list",
			req: &sync.ListConflictsRequest{
				UserId: userID,
			},
			setupMock: func() {
				srv.EXPECT().ListConflicts(ctx, userID).Return(conflicts, nil)
			},
			want: &sync.ListConflictsResponse{
				Conflicts: []*sync.Conflict{
					{
						ConflictId: "conflict1",
						DataId:     "data1",
						ClientData: []byte("clientdata"),
						ServerData: []byte("serverdata"),
						Resolved:   false,
						CreatedAt:  now.Unix(),
						UpdatedAt:  now.Unix(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid request",
			req: &sync.ListConflictsRequest{
				UserId: "",
			},
			setupMock: func() {},
			want:      nil,
			wantErr:   true,
			errCode:   codes.InvalidArgument,
		},
		{
			name: "internal error",
			req: &sync.ListConflictsRequest{
				UserId: userID,
			},
			setupMock: func() {
				srv.EXPECT().ListConflicts(ctx, userID).Return(nil, errors.New("internal error"))
			},
			want:    nil,
			wantErr: true,
			errCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := api.ListConflicts(ctx, tt.req)

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

func TestToModelDataType(t *testing.T) {
	tests := []struct {
		name string
		in   sync.DataType
		want models.DataType
	}{
		{"login password", sync.DataType_LOGIN_PASSWORD, models.LoginPassword},
		{"text", sync.DataType_TEXT, models.Text},
		{"binary", sync.DataType_BINARY, models.Binary},
		{"card", sync.DataType_CARD, models.Card},
		{"default", sync.DataType(999), models.LoginPassword},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, toModelDataType(tt.in))
		})
	}
}

func TestToProtoDataType(t *testing.T) {
	tests := []struct {
		name string
		in   models.DataType
		want sync.DataType
	}{
		{"login password", models.LoginPassword, sync.DataType_LOGIN_PASSWORD},
		{"text", models.Text, sync.DataType_TEXT},
		{"binary", models.Binary, sync.DataType_BINARY},
		{"card", models.Card, sync.DataType_CARD},
		{"default", models.DataType(999), sync.DataType_LOGIN_PASSWORD},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, toProtoDataType(tt.in))
		})
	}
}

func TestToModelResolutionStrategy(t *testing.T) {
	tests := []struct {
		name string
		in   sync.ResolutionStrategy
		want models.ResolutionStrategy
	}{
		{"use client", sync.ResolutionStrategy_USE_CLIENT_VERSION, models.UseClientVersion},
		{"use server", sync.ResolutionStrategy_USE_SERVER_VERSION, models.UseServerVersion},
		{"merge", sync.ResolutionStrategy_MERGE_VERSIONS, models.MergeVersions},
		{"default", sync.ResolutionStrategy(999), models.UseClientVersion},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, toModelResolutionStrategy(tt.in))
		})
	}
}

func TestToModelOperation(t *testing.T) {
	tests := []struct {
		name string
		in   sync.Operation
		want models.Operation
	}{
		{"add", sync.Operation_ADD, models.Add},
		{"update", sync.Operation_UPDATE, models.Update},
		{"delete", sync.Operation_DELETE, models.Delete},
		{"default", sync.Operation(999), models.Add},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, toModelOperation(tt.in))
		})
	}
}
