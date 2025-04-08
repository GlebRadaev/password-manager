package interceptor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/GlebRadaev/password-manager/pkg/auth"
)

func TestAuthInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := auth.NewMockAuthServiceClient(ctrl)
	interceptor := AuthInterceptor(mockClient)

	tests := []struct {
		name        string
		method      string
		setupCtx    func() context.Context
		setupMock   func()
		expectedErr error
	}{
		{
			name:   "public method - no auth required",
			method: "/api.auth.AuthService/Login",
			setupCtx: func() context.Context {
				return context.Background()
			},
			setupMock:   func() {},
			expectedErr: nil,
		},
		{
			name:   "missing token - unauthenticated",
			method: "/api.sync.SyncService/SyncData",
			setupCtx: func() context.Context {
				return context.Background()
			},
			setupMock:   func() {},
			expectedErr: status.Error(codes.Unauthenticated, "authorization token missing"),
		},
		{
			name:   "invalid token format - missing Bearer prefix",
			method: "/api.sync.SyncService/SyncData",
			setupCtx: func() context.Context {
				md := metadata.New(map[string]string{
					"authorization": "InvalidTokenFormat",
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			setupMock:   func() {},
			expectedErr: status.Error(codes.Unauthenticated, "authorization token missing"),
		},
		{
			name:   "valid token in incoming context",
			method: "/api.sync.SyncService/SyncData",
			setupCtx: func() context.Context {
				md := metadata.New(map[string]string{
					"authorization": "Bearer valid_token",
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			setupMock: func() {
				mockClient.EXPECT().ValidateToken(
					gomock.Any(),
					&auth.ValidateTokenRequest{Token: "valid_token"},
				).Return(&auth.ValidateTokenResponse{
					Valid:  true,
					UserId: "user123",
				}, nil)
			},
			expectedErr: nil,
		},
		{
			name:   "valid token in outgoing context",
			method: "/api.sync.SyncService/SyncData",
			setupCtx: func() context.Context {
				md := metadata.New(map[string]string{
					"authorization": "Bearer valid_token",
				})
				return metadata.NewOutgoingContext(context.Background(), md)
			},
			setupMock: func() {
				mockClient.EXPECT().ValidateToken(
					gomock.Any(),
					&auth.ValidateTokenRequest{Token: "valid_token"},
				).Return(&auth.ValidateTokenResponse{
					Valid:  true,
					UserId: "user123",
				}, nil)
			},
			expectedErr: nil,
		},
		{
			name:   "invalid token - validation failed",
			method: "/api.sync.SyncService/SyncData",
			setupCtx: func() context.Context {
				md := metadata.New(map[string]string{
					"authorization": "Bearer invalid_token",
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			setupMock: func() {
				mockClient.EXPECT().ValidateToken(
					gomock.Any(),
					&auth.ValidateTokenRequest{Token: "invalid_token"},
				).Return(&auth.ValidateTokenResponse{
					Valid: false,
				}, nil)
			},
			expectedErr: status.Error(codes.Unauthenticated, "invalid token"),
		},
		{
			name:   "token validation error",
			method: "/api.sync.SyncService/SyncData",
			setupCtx: func() context.Context {
				md := metadata.New(map[string]string{
					"authorization": "Bearer error_token",
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			setupMock: func() {
				mockClient.EXPECT().ValidateToken(
					gomock.Any(),
					&auth.ValidateTokenRequest{Token: "error_token"},
				).Return(nil, status.Error(codes.Internal, "validation error"))
			},
			expectedErr: status.Error(codes.Unauthenticated, "invalid token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx()
			tt.setupMock()

			var invokedCtx context.Context
			invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				invokedCtx = ctx
				return nil
			}

			err := interceptor(ctx, tt.method, nil, nil, nil, invoker)
			assert.Equal(t, tt.expectedErr, err)

			if tt.expectedErr == nil && err == nil && !isPublicMethod(tt.method) {
				userID := invokedCtx.Value(UserIDKey)
				assert.NotNil(t, userID)
				assert.Equal(t, "user123", userID)
			}
		})
	}
}

func TestIsPublicMethod(t *testing.T) {
	tests := []struct {
		method   string
		expected bool
	}{
		{method: "/api.auth.AuthService/Register", expected: true},
		{method: "/api.auth.AuthService/Login", expected: true},
		{method: "/api.auth.AuthService/ValidateToken", expected: true},
		{method: "/api.auth.AuthService/RefreshToken", expected: true},
		{method: "/api.sync.SyncService/SyncData", expected: false},
		{method: "/api.other.Service/Method", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			assert.Equal(t, tt.expected, isPublicMethod(tt.method))
		})
	}
}
