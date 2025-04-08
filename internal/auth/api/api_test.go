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

	"github.com/GlebRadaev/password-manager/internal/auth/models"
	"github.com/GlebRadaev/password-manager/internal/auth/service"
	"github.com/GlebRadaev/password-manager/pkg/auth"
)

func setupMocks(t *testing.T) (context.Context, *API, *MockService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	srv := NewMockService(ctrl)
	api := New(srv)

	return context.Background(), api, srv
}

func TestRegister(t *testing.T) {
	ctx, api, srv := setupMocks(t)

	username := "testuser"
	password := "testpass"
	email := "test@example.com"
	userID := uuid.NewString()

	tests := []struct {
		name      string
		req       *auth.RegisterRequest
		setupMock func()
		want      *auth.RegisterResponse
		wantErr   bool
		errCode   codes.Code
	}{
		{
			name: "successful registration",
			req: &auth.RegisterRequest{
				Username: username,
				Password: password,
				Email:    email,
			},
			setupMock: func() {
				srv.EXPECT().Register(ctx, username, password, email).Return(userID, nil)
			},
			want: &auth.RegisterResponse{
				UserId:  userID,
				Message: "User registered successfully",
			},
			wantErr: false,
		},
		{
			name: "user already exists",
			req: &auth.RegisterRequest{
				Username: username,
				Password: password,
				Email:    email,
			},
			setupMock: func() {
				srv.EXPECT().Register(ctx, username, password, email).Return("", service.ErrUserExists)
			},
			want:    nil,
			wantErr: true,
			errCode: codes.AlreadyExists,
		},
		{
			name: "internal error",
			req: &auth.RegisterRequest{
				Username: username,
				Password: password,
				Email:    email,
			},
			setupMock: func() {
				srv.EXPECT().Register(ctx, username, password, email).Return("", errors.New("internal error"))
			},
			want:    nil,
			wantErr: true,
			errCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := api.Register(ctx, tt.req)

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

func TestLogin(t *testing.T) {
	ctx, api, srv := setupMocks(t)

	username := "testuser"
	password := "testpass"
	accessToken := "access_token"
	refreshToken := "refresh_token"
	expiresIn := int64(3600)

	tests := []struct {
		name      string
		req       *auth.LoginRequest
		setupMock func()
		want      *auth.LoginResponse
		wantErr   bool
		errCode   codes.Code
	}{
		{
			name: "successful login",
			req: &auth.LoginRequest{
				Username: username,
				Password: password,
			},
			setupMock: func() {
				srv.EXPECT().Login(ctx, username, password).Return(accessToken, refreshToken, expiresIn, nil)
			},
			want: &auth.LoginResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
				ExpiresIn:    expiresIn,
			},
			wantErr: false,
		},
		{
			name: "invalid credentials",
			req: &auth.LoginRequest{
				Username: username,
				Password: password,
			},
			setupMock: func() {
				srv.EXPECT().Login(ctx, username, password).Return("", "", int64(0), service.ErrInvalidCredentials)
			},
			want:    nil,
			wantErr: true,
			errCode: codes.Unauthenticated,
		},
		{
			name: "internal error",
			req: &auth.LoginRequest{
				Username: username,
				Password: password,
			},
			setupMock: func() {
				srv.EXPECT().Login(ctx, username, password).Return("", "", int64(0), errors.New("internal error"))
			},
			want:    nil,
			wantErr: true,
			errCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := api.Login(ctx, tt.req)

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

func TestValidateToken(t *testing.T) {
	ctx, api, srv := setupMocks(t)

	token := "test_token"
	userID := uuid.NewString()

	tests := []struct {
		name      string
		req       *auth.ValidateTokenRequest
		setupMock func()
		want      *auth.ValidateTokenResponse
		wantErr   bool
		errCode   codes.Code
	}{
		{
			name: "valid token",
			req: &auth.ValidateTokenRequest{
				Token: token,
			},
			setupMock: func() {
				srv.EXPECT().ValidateToken(ctx, token).Return(true, userID, nil)
			},
			want: &auth.ValidateTokenResponse{
				Valid:  true,
				UserId: userID,
			},
			wantErr: false,
		},
		{
			name: "invalid token",
			req: &auth.ValidateTokenRequest{
				Token: token,
			},
			setupMock: func() {
				srv.EXPECT().ValidateToken(ctx, token).Return(false, "", errors.New("invalid token"))
			},
			want:    nil,
			wantErr: true,
			errCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := api.ValidateToken(ctx, tt.req)

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

func TestGenerateOTP(t *testing.T) {
	ctx, api, srv := setupMocks(t)

	userID := uuid.NewString()
	deviceID := "device1"
	otpCode := "123456"

	tests := []struct {
		name      string
		req       *auth.GenerateOTPRequest
		setupMock func()
		want      *auth.GenerateOTPResponse
		wantErr   bool
		errCode   codes.Code
	}{
		{
			name: "successful OTP generation",
			req: &auth.GenerateOTPRequest{
				UserId:   userID,
				DeviceId: deviceID,
			},
			setupMock: func() {
				srv.EXPECT().GenerateOTP(ctx, userID, deviceID).Return(otpCode, nil)
			},
			want: &auth.GenerateOTPResponse{
				OtpCode: otpCode,
			},
			wantErr: false,
		},
		{
			name: "internal error",
			req: &auth.GenerateOTPRequest{
				UserId:   userID,
				DeviceId: deviceID,
			},
			setupMock: func() {
				srv.EXPECT().GenerateOTP(ctx, userID, deviceID).Return("", errors.New("internal error"))
			},
			want:    nil,
			wantErr: true,
			errCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := api.GenerateOTP(ctx, tt.req)

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

func TestValidateOTP(t *testing.T) {
	ctx, api, srv := setupMocks(t)

	userID := uuid.NewString()
	otpCode := "123456"
	deviceID := "device1"

	tests := []struct {
		name      string
		req       *auth.ValidateOTPRequest
		setupMock func()
		want      *auth.ValidateOTPResponse
		wantErr   bool
		errCode   codes.Code
	}{
		{
			name: "valid OTP",
			req: &auth.ValidateOTPRequest{
				UserId:   userID,
				OtpCode:  otpCode,
				DeviceId: deviceID,
			},
			setupMock: func() {
				srv.EXPECT().ValidateOTP(ctx, userID, otpCode, deviceID).Return(true, nil)
			},
			want: &auth.ValidateOTPResponse{
				Valid: true,
			},
			wantErr: false,
		},
		{
			name: "invalid OTP",
			req: &auth.ValidateOTPRequest{
				UserId:   userID,
				OtpCode:  otpCode,
				DeviceId: deviceID,
			},
			setupMock: func() {
				srv.EXPECT().ValidateOTP(ctx, userID, otpCode, deviceID).Return(false, service.ErrInvalidOTP)
			},
			want:    nil,
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "internal error",
			req: &auth.ValidateOTPRequest{
				UserId:   userID,
				OtpCode:  otpCode,
				DeviceId: deviceID,
			},
			setupMock: func() {
				srv.EXPECT().ValidateOTP(ctx, userID, otpCode, deviceID).Return(false, errors.New("internal error"))
			},
			want:    nil,
			wantErr: true,
			errCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := api.ValidateOTP(ctx, tt.req)

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

func TestListSessions(t *testing.T) {
	ctx, api, srv := setupMocks(t)

	userID := uuid.NewString()
	now := time.Now()
	sessions := []models.Session{
		{
			SessionID:  "session1",
			UserID:     userID,
			DeviceInfo: "device1",
			CreatedAt:  now,
			ExpiresAt:  now.Add(time.Hour),
		},
	}

	tests := []struct {
		name      string
		req       *auth.ListSessionsRequest
		setupMock func()
		want      *auth.ListSessionsResponse
		wantErr   bool
		errCode   codes.Code
	}{
		{
			name: "successful session list",
			req: &auth.ListSessionsRequest{
				UserId: userID,
			},
			setupMock: func() {
				srv.EXPECT().ListSessions(ctx, userID).Return(sessions, nil)
			},
			want: &auth.ListSessionsResponse{
				Sessions: []*auth.Session{
					{
						SessionId:  "session1",
						DeviceInfo: "device1",
						CreatedAt:  now.Unix(),
						ExpiresAt:  now.Add(time.Hour).Unix(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "internal error",
			req: &auth.ListSessionsRequest{
				UserId: userID,
			},
			setupMock: func() {
				srv.EXPECT().ListSessions(ctx, userID).Return(nil, errors.New("internal error"))
			},
			want:    nil,
			wantErr: true,
			errCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := api.ListSessions(ctx, tt.req)

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

func TestTerminateSession(t *testing.T) {
	ctx, api, srv := setupMocks(t)

	sessionID := "session1"

	tests := []struct {
		name      string
		req       *auth.TerminateSessionRequest
		setupMock func()
		want      *auth.TerminateSessionResponse
		wantErr   bool
		errCode   codes.Code
	}{
		{
			name: "successful session termination",
			req: &auth.TerminateSessionRequest{
				SessionId: sessionID,
			},
			setupMock: func() {
				srv.EXPECT().TerminateSession(ctx, sessionID).Return(nil)
			},
			want: &auth.TerminateSessionResponse{
				Message: "Session terminated successfully",
			},
			wantErr: false,
		},
		{
			name: "session not found",
			req: &auth.TerminateSessionRequest{
				SessionId: sessionID,
			},
			setupMock: func() {
				srv.EXPECT().TerminateSession(ctx, sessionID).Return(service.ErrSessionNotFound)
			},
			want:    nil,
			wantErr: true,
			errCode: codes.NotFound,
		},
		{
			name: "internal error",
			req: &auth.TerminateSessionRequest{
				SessionId: sessionID,
			},
			setupMock: func() {
				srv.EXPECT().TerminateSession(ctx, sessionID).Return(errors.New("internal error"))
			},
			want:    nil,
			wantErr: true,
			errCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := api.TerminateSession(ctx, tt.req)

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

func TestRefreshToken(t *testing.T) {
	ctx, api, srv := setupMocks(t)

	refreshToken := "refresh_token"
	newAccessToken := "new_access_token"
	newRefreshToken := "new_refresh_token"
	expiresIn := int64(3600)

	tests := []struct {
		name      string
		req       *auth.RefreshTokenRequest
		setupMock func()
		want      *auth.RefreshTokenResponse
		wantErr   bool
		errCode   codes.Code
	}{
		{
			name: "successful token refresh",
			req: &auth.RefreshTokenRequest{
				RefreshToken: refreshToken,
			},
			setupMock: func() {
				srv.EXPECT().RefreshToken(ctx, refreshToken).Return(newAccessToken, newRefreshToken, expiresIn, nil)
			},
			want: &auth.RefreshTokenResponse{
				AccessToken:  newAccessToken,
				RefreshToken: newRefreshToken,
				ExpiresIn:    expiresIn,
			},
			wantErr: false,
		},
		{
			name: "invalid refresh token",
			req: &auth.RefreshTokenRequest{
				RefreshToken: refreshToken,
			},
			setupMock: func() {
				srv.EXPECT().RefreshToken(ctx, refreshToken).Return("", "", int64(0), errors.New("invalid token"))
			},
			want:    nil,
			wantErr: true,
			errCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := api.RefreshToken(ctx, tt.req)

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
