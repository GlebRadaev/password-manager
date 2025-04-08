package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/GlebRadaev/password-manager/internal/auth/models"
)

func setupMocks(t *testing.T) (context.Context, *Service, *MockRepo, *MockAuthManager) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	repo := NewMockRepo(ctrl)
	authManager := NewMockAuthManager(ctrl)
	txManager := NewMockTxManager(ctrl)
	srv := New(repo, txManager, authManager)
	return context.Background(), srv, repo, authManager
}

func TestRegister(t *testing.T) {
	ctx, srv, repo, authManager := setupMocks(t)

	username := "testuser"
	password := "testpass"
	email := "test@example.com"
	passwordHash := "hashedpassword"

	tests := []struct {
		name      string
		username  string
		password  string
		email     string
		setupMock func()
		wantErr   error
	}{
		{
			name:     "successful registration",
			username: username,
			password: password,
			email:    email,
			setupMock: func() {
				repo.EXPECT().CheckExists(ctx, username, email).Return(false, nil)
				authManager.EXPECT().Hash(password).Return(passwordHash, nil)
				repo.EXPECT().CreateUser(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, u models.User) error {
					assert.NotEmpty(t, u.ID)
					assert.Equal(t, username, u.Username)
					assert.Equal(t, passwordHash, u.PasswordHash)
					assert.Equal(t, email, u.Email)
					return nil
				})
			},
			wantErr: nil,
		},
		{
			name:     "user already exists",
			username: username,
			password: password,
			email:    email,
			setupMock: func() {
				repo.EXPECT().CheckExists(ctx, username, email).Return(true, nil)
			},
			wantErr: ErrUserExists,
		},
		{
			name:     "check exists error",
			username: username,
			password: password,
			email:    email,
			setupMock: func() {
				repo.EXPECT().CheckExists(ctx, username, email).Return(false, errors.New("db error"))
			},
			wantErr: errors.New("failed to check user existence: db error"),
		},
		{
			name:     "hash password error",
			username: username,
			password: password,
			email:    email,
			setupMock: func() {
				repo.EXPECT().CheckExists(ctx, username, email).Return(false, nil)
				authManager.EXPECT().Hash(password).Return("", errors.New("hash error"))
			},
			wantErr: errors.New("failed to hash password: hash error"),
		},
		{
			name:     "create user error",
			username: username,
			password: password,
			email:    email,
			setupMock: func() {
				repo.EXPECT().CheckExists(ctx, username, email).Return(false, nil)
				authManager.EXPECT().Hash(password).Return(passwordHash, nil)
				repo.EXPECT().CreateUser(ctx, gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: errors.New("failed to create user: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			_, err := srv.Register(ctx, tt.username, tt.password, tt.email)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestLogin(t *testing.T) {
	ctx, srv, repo, authManager := setupMocks(t)

	username := "testuser"
	password := "testpass"
	userID := uuid.NewString()
	hashedPassword := "hashedpassword"
	accessToken := "access_token"
	refreshToken := "refresh_token"
	expiresAt := time.Now().Add(time.Hour)
	user := models.User{
		ID:           userID,
		Username:     username,
		PasswordHash: hashedPassword,
	}

	tests := []struct {
		name      string
		username  string
		password  string
		setupMock func()
		wantToken string
		wantRT    string
		wantExp   int64
		wantErr   error
	}{
		{
			name:     "successful login",
			username: username,
			password: password,
			setupMock: func() {
				repo.EXPECT().GetUserByUsername(ctx, username).Return(user, nil)
				authManager.EXPECT().Compare(hashedPassword, password).Return(nil)
				authManager.EXPECT().GenerateToken(userID).Return(accessToken, expiresAt, nil)
				authManager.EXPECT().GenerateRefreshToken(userID).Return(refreshToken, expiresAt, nil)
				repo.EXPECT().CreateSession(ctx, gomock.Any()).Return(nil)
			},
			wantToken: accessToken,
			wantRT:    refreshToken,
			wantExp:   expiresAt.Unix(),
			wantErr:   nil,
		},
		{
			name:     "user not found",
			username: username,
			password: password,
			setupMock: func() {
				repo.EXPECT().GetUserByUsername(ctx, username).Return(models.User{}, errors.New("not found"))
			},
			wantToken: "",
			wantRT:    "",
			wantExp:   0,
			wantErr:   errors.New("failed to get user: not found"),
		},
		{
			name:     "invalid credentials",
			username: username,
			password: "wrongpass",
			setupMock: func() {
				repo.EXPECT().GetUserByUsername(ctx, username).Return(user, nil)
				authManager.EXPECT().Compare(hashedPassword, "wrongpass").Return(errors.New("invalid"))
			},
			wantToken: "",
			wantRT:    "",
			wantExp:   0,
			wantErr:   ErrInvalidCredentials,
		},
		{
			name:     "generate token error",
			username: username,
			password: password,
			setupMock: func() {
				repo.EXPECT().GetUserByUsername(ctx, username).Return(user, nil)
				authManager.EXPECT().Compare(hashedPassword, password).Return(nil)
				authManager.EXPECT().GenerateToken(userID).Return("", time.Time{}, errors.New("token error"))
			},
			wantToken: "",
			wantRT:    "",
			wantExp:   0,
			wantErr:   errors.New("failed to generate access token: token error"),
		},
		{
			name:     "create session error",
			username: username,
			password: password,
			setupMock: func() {
				repo.EXPECT().GetUserByUsername(ctx, username).Return(user, nil)
				authManager.EXPECT().Compare(hashedPassword, password).Return(nil)
				authManager.EXPECT().GenerateToken(userID).Return(accessToken, expiresAt, nil)
				authManager.EXPECT().GenerateRefreshToken(userID).Return(refreshToken, expiresAt, nil)
				repo.EXPECT().CreateSession(ctx, gomock.Any()).Return(errors.New("session error"))
			},
			wantToken: "",
			wantRT:    "",
			wantExp:   0,
			wantErr:   errors.New("failed to create session: session error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			token, rt, exp, err := srv.Login(ctx, tt.username, tt.password)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantToken, token)
				assert.Equal(t, tt.wantRT, rt)
				assert.Equal(t, tt.wantExp, exp)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	ctx, srv, _, authManager := setupMocks(t)

	token := "test_token"
	userID := uuid.NewString()

	tests := []struct {
		name      string
		token     string
		setupMock func()
		wantValid bool
		wantID    string
		wantErr   error
	}{
		{
			name:  "valid token",
			token: token,
			setupMock: func() {
				authManager.EXPECT().ValidateToken(token).Return(userID, nil)
			},
			wantValid: true,
			wantID:    userID,
			wantErr:   nil,
		},
		{
			name:  "invalid token",
			token: token,
			setupMock: func() {
				authManager.EXPECT().ValidateToken(token).Return("", errors.New("invalid"))
			},
			wantValid: false,
			wantID:    "",
			wantErr:   errors.New("failed to validate token: invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			valid, id, err := srv.ValidateToken(ctx, tt.token)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValid, valid)
				assert.Equal(t, tt.wantID, id)
			}
		})
	}
}

func TestGenerateOTP(t *testing.T) {
	ctx, srv, repo, authManager := setupMocks(t)

	userID := uuid.NewString()
	deviceID := "device1"
	otpCode := "123456"
	expiresAt := time.Now().Add(time.Minute * 5)

	tests := []struct {
		name      string
		userID    string
		deviceID  string
		setupMock func()
		want      string
		wantErr   error
	}{
		{
			name:     "successful OTP generation",
			userID:   userID,
			deviceID: deviceID,
			setupMock: func() {
				authManager.EXPECT().GenerateOTP().Return(otpCode, expiresAt, nil)
				repo.EXPECT().CreateOTP(ctx, models.OTP{
					UserID:    userID,
					OTPCode:   otpCode,
					ExpiresAt: expiresAt,
					DeviceID:  deviceID,
				}).Return(nil)
			},
			want:    otpCode,
			wantErr: nil,
		},
		{
			name:     "generate OTP error",
			userID:   userID,
			deviceID: deviceID,
			setupMock: func() {
				authManager.EXPECT().GenerateOTP().Return("", time.Time{}, errors.New("otp error"))
			},
			want:    "",
			wantErr: errors.New("failed to generate OTP: otp error"),
		},
		{
			name:     "create OTP error",
			userID:   userID,
			deviceID: deviceID,
			setupMock: func() {
				authManager.EXPECT().GenerateOTP().Return(otpCode, expiresAt, nil)
				repo.EXPECT().CreateOTP(ctx, gomock.Any()).Return(errors.New("db error"))
			},
			want:    "",
			wantErr: errors.New("failed to create OTP: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			got, err := srv.GenerateOTP(ctx, tt.userID, tt.deviceID)

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

func TestValidateOTP(t *testing.T) {
	ctx, srv, repo, authManager := setupMocks(t)

	userID := uuid.NewString()
	otpCode := "123456"
	deviceID := "device1"
	storedOTP := models.OTP{
		UserID:    userID,
		OTPCode:   "654321",
		ExpiresAt: time.Now().Add(time.Minute * 5),
		DeviceID:  deviceID,
	}

	tests := []struct {
		name      string
		userID    string
		otpCode   string
		deviceID  string
		setupMock func()
		want      bool
		wantErr   string
	}{
		{
			name:     "valid OTP",
			userID:   userID,
			otpCode:  otpCode,
			deviceID: deviceID,
			setupMock: func() {
				repo.EXPECT().GetOTP(ctx, userID, otpCode, deviceID).Return(storedOTP, nil)
				authManager.EXPECT().ValidateOTP(storedOTP.OTPCode, storedOTP.ExpiresAt, otpCode).Return(true, nil)
			},
			want:    true,
			wantErr: "",
		},
		{
			name:     "invalid OTP",
			userID:   userID,
			otpCode:  "wrong",
			deviceID: deviceID,
			setupMock: func() {
				repo.EXPECT().GetOTP(ctx, userID, "wrong", deviceID).Return(storedOTP, nil)
				authManager.EXPECT().ValidateOTP(storedOTP.OTPCode, storedOTP.ExpiresAt, "wrong").Return(false, nil)
			},
			want:    false,
			wantErr: ErrInvalidOTP.Error(),
		},
		{
			name:     "get OTP error",
			userID:   userID,
			otpCode:  otpCode,
			deviceID: deviceID,
			setupMock: func() {
				repo.EXPECT().GetOTP(ctx, userID, otpCode, deviceID).Return(models.OTP{}, errors.New("db error"))
			},
			want:    false,
			wantErr: "failed to get OTP: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			valid, err := srv.ValidateOTP(ctx, tt.userID, tt.otpCode, tt.deviceID)

			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, valid)
		})
	}
}

func TestListSessions(t *testing.T) {
	ctx, srv, repo, _ := setupMocks(t)

	userID := uuid.NewString()
	sessions := []models.Session{
		{SessionID: "sess1", UserID: userID, DeviceInfo: "device1"},
		{SessionID: "sess2", UserID: userID, DeviceInfo: "device2"},
	}

	tests := []struct {
		name      string
		userID    string
		setupMock func()
		want      []models.Session
		wantErr   error
	}{
		{
			name:   "successful list sessions",
			userID: userID,
			setupMock: func() {
				repo.EXPECT().ListSessions(ctx, userID).Return(sessions, nil)
			},
			want:    sessions,
			wantErr: nil,
		},
		{
			name:   "list sessions error",
			userID: userID,
			setupMock: func() {
				repo.EXPECT().ListSessions(ctx, userID).Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: errors.New("failed to list sessions: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			got, err := srv.ListSessions(ctx, tt.userID)

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

func TestTerminateSession(t *testing.T) {
	ctx, srv, repo, _ := setupMocks(t)

	sessionID := "sess1"

	tests := []struct {
		name       string
		sessionID  string
		setupMock  func()
		wantErr    error
		wantErrMsg string
	}{
		{
			name:      "successful session termination",
			sessionID: sessionID,
			setupMock: func() {
				repo.EXPECT().TerminateSession(ctx, sessionID).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:      "session not found",
			sessionID: sessionID,
			setupMock: func() {
				repo.EXPECT().TerminateSession(ctx, sessionID).Return(ErrSessionNotFound)
			},
			wantErr:    ErrSessionNotFound,
			wantErrMsg: "failed to terminate session: session not found",
		},
		{
			name:      "terminate session error",
			sessionID: sessionID,
			setupMock: func() {
				repo.EXPECT().TerminateSession(ctx, sessionID).Return(errors.New("db error"))
			},
			wantErr:    errors.New("db error"),
			wantErrMsg: "failed to terminate session: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := srv.TerminateSession(ctx, tt.sessionID)

			if tt.wantErr != nil {
				assert.Error(t, err)
				if tt.wantErrMsg != "" {
					assert.EqualError(t, err, tt.wantErrMsg)
				} else {
					assert.Equal(t, tt.wantErr, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	ctx, srv, _, authManager := setupMocks(t)

	userID := uuid.NewString()
	refreshToken := "refresh_token"
	newAccessToken := "new_access_token"
	newRefreshToken := "new_refresh_token"
	expiresAt := time.Now().Add(time.Hour)

	tests := []struct {
		name         string
		refreshToken string
		setupMock    func()
		wantToken    string
		wantRT       string
		wantExp      int64
		wantErr      error
		wantErrMsg   string
	}{
		{
			name:         "successful token refresh",
			refreshToken: refreshToken,
			setupMock: func() {
				authManager.EXPECT().ValidateRefreshToken(refreshToken).Return(userID, nil)
				authManager.EXPECT().GenerateToken(userID).Return(newAccessToken, expiresAt, nil)
				authManager.EXPECT().GenerateRefreshToken(userID).Return(newRefreshToken, expiresAt, nil)
			},
			wantToken: newAccessToken,
			wantRT:    newRefreshToken,
			wantExp:   expiresAt.Unix(),
			wantErr:   nil,
		},
		{
			name:         "invalid refresh token",
			refreshToken: "invalid",
			setupMock: func() {
				authManager.EXPECT().ValidateRefreshToken("invalid").Return("", errors.New("invalid token"))
			},
			wantToken: "",
			wantRT:    "",
			wantExp:   0,
			wantErr:   errors.New("invalid refresh token: invalid token"),
		},
		{
			name:         "generate access token error",
			refreshToken: refreshToken,
			setupMock: func() {
				authManager.EXPECT().ValidateRefreshToken(refreshToken).Return(userID, nil)
				authManager.EXPECT().GenerateToken(userID).Return("", time.Time{}, errors.New("token error"))
			},
			wantToken: "",
			wantRT:    "",
			wantExp:   0,
			wantErr:   errors.New("failed to generate access token: token error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			token, rt, exp, err := srv.RefreshToken(ctx, tt.refreshToken)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantToken, token)
				assert.Equal(t, tt.wantRT, rt)
				assert.Equal(t, tt.wantExp, exp)
			}
		})
	}
}
