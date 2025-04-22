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
	"go.uber.org/mock/gomock"

	"github.com/GlebRadaev/password-manager/internal/auth/models"
)

func NewMock(t *testing.T) (*Repository, pgxmock.PgxPoolIface) {
	ctrl := gomock.NewController(t)

	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	repo := New(mockDB)

	t.Cleanup(func() {
		mockDB.Close()
		ctrl.Finish()
	})

	return repo, mockDB
}

func TestRepository_CreateUser(t *testing.T) {
	repo, mock := NewMock(t)

	user := models.User{
		ID:           "user1",
		Username:     "testuser",
		PasswordHash: "hash",
		Email:        "test@example.com",
		CreatedAt:    time.Now(),
	}

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
	}{
		{
			name: "successful user creation",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`
					INSERT INTO auth.users (user_id, username, password_hash, email, created_at)
					VALUES ($1, $2, $3, $4, $5)`)).
					WithArgs(user.ID, user.Username, user.PasswordHash, user.Email, user.CreatedAt).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectErr: false,
		},
		{
			name: "database error",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO auth.users`)).
					WithArgs(user.ID, user.Username, user.PasswordHash, user.Email, user.CreatedAt).
					WillReturnError(errors.New("db error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.CreateUser(context.Background(), user)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to create user")
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetUserByUsername(t *testing.T) {
	repo, mock := NewMock(t)

	username := "testuser"
	expectedUser := models.User{
		ID:           "user1",
		Username:     username,
		PasswordHash: "hash",
		Email:        "test@example.com",
		CreatedAt:    time.Now(),
	}

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
		expected  models.User
	}{
		{
			name: "successful user retrieval",
			mockSetup: func() {
				rows := pgxmock.NewRows([]string{"user_id", "username", "password_hash", "email", "created_at"}).
					AddRow(expectedUser.ID, expectedUser.Username, expectedUser.PasswordHash, expectedUser.Email, expectedUser.CreatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT user_id, username, password_hash, email, created_at
					FROM auth.users
					WHERE username = $1`)).
					WithArgs(username).
					WillReturnRows(rows)
			},
			expectErr: false,
			expected:  expectedUser,
		},
		{
			name: "user not found",
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
					WithArgs(username).
					WillReturnError(pgx.ErrNoRows)
			},
			expectErr: true,
		},
		{
			name: "database error",
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
					WithArgs(username).
					WillReturnError(errors.New("db error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, err := repo.GetUserByUsername(context.Background(), username)

			if tt.expectErr {
				assert.Error(t, err)
				if errors.Is(err, pgx.ErrNoRows) {
					assert.Contains(t, err.Error(), "user not found")
				} else {
					assert.Contains(t, err.Error(), "failed to get user")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, user)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_CheckExists(t *testing.T) {
	repo, mock := NewMock(t)

	username := "testuser"
	email := "test@example.com"

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
		expected  bool
	}{
		{
			name: "user exists",
			mockSetup: func() {
				rows := pgxmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT EXISTS (
						SELECT 1
						FROM auth.users
						WHERE username = $1 OR email = $2
					)`)).
					WithArgs(username, email).
					WillReturnRows(rows)
			},
			expectErr: false,
			expected:  true,
		},
		{
			name: "user does not exist",
			mockSetup: func() {
				rows := pgxmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS`)).
					WithArgs(username, email).
					WillReturnRows(rows)
			},
			expectErr: false,
			expected:  false,
		},
		{
			name: "database error",
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS`)).
					WithArgs(username, email).
					WillReturnError(errors.New("db error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			exists, err := repo.CheckExists(context.Background(), username, email)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to check user existence")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, exists)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_CreateOTP(t *testing.T) {
	repo, mock := NewMock(t)

	otp := models.OTP{
		UserID:    "user1",
		OTPCode:   "123456",
		ExpiresAt: time.Now().Add(5 * time.Minute),
		DeviceID:  "device1",
	}

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
	}{
		{
			name: "successful OTP creation",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`
					INSERT INTO auth.otp_codes (user_id, otp_code, expires_at, device_id)
					VALUES ($1, $2, $3, $4)`)).
					WithArgs(otp.UserID, otp.OTPCode, otp.ExpiresAt, otp.DeviceID).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectErr: false,
		},
		{
			name: "database error",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO auth.otp_codes`)).
					WithArgs(otp.UserID, otp.OTPCode, otp.ExpiresAt, otp.DeviceID).
					WillReturnError(errors.New("db error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.CreateOTP(context.Background(), otp)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to create OTP")
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetOTP(t *testing.T) {
	repo, mock := NewMock(t)

	otp := models.OTP{
		UserID:    "user1",
		OTPCode:   "123456",
		ExpiresAt: time.Now().Add(5 * time.Minute),
		DeviceID:  "device1",
	}

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
		expected  models.OTP
	}{
		{
			name: "successful OTP retrieval",
			mockSetup: func() {
				rows := pgxmock.NewRows([]string{"user_id", "otp_code", "expires_at", "device_id"}).
					AddRow(otp.UserID, otp.OTPCode, otp.ExpiresAt, otp.DeviceID)
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT user_id, otp_code, expires_at, device_id
					FROM auth.otp_codes
					WHERE user_id = $1 AND otp_code = $2 AND device_id = $3`)).
					WithArgs(otp.UserID, otp.OTPCode, otp.DeviceID).
					WillReturnRows(rows)
			},
			expectErr: false,
			expected:  otp,
		},
		{
			name: "OTP not found",
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
					WithArgs(otp.UserID, otp.OTPCode, otp.DeviceID).
					WillReturnError(pgx.ErrNoRows)
			},
			expectErr: true,
		},
		{
			name: "database error",
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
					WithArgs(otp.UserID, otp.OTPCode, otp.DeviceID).
					WillReturnError(errors.New("db error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := repo.GetOTP(context.Background(), otp.UserID, otp.OTPCode, otp.DeviceID)

			if tt.expectErr {
				assert.Error(t, err)
				if errors.Is(err, pgx.ErrNoRows) {
					assert.Contains(t, err.Error(), "OTP not found")
				} else {
					assert.Contains(t, err.Error(), "failed to get OTP")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_CreateSession(t *testing.T) {
	repo, mock := NewMock(t)

	session := models.Session{
		SessionID:  "session1",
		UserID:     "user1",
		DeviceInfo: "device1",
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
	}{
		{
			name: "successful session creation",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`
					INSERT INTO auth.sessions (session_id, user_id, device_info, created_at, expires_at)
					VALUES ($1, $2, $3, $4, $5)`)).
					WithArgs(session.SessionID, session.UserID, session.DeviceInfo, session.CreatedAt, session.ExpiresAt).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectErr: false,
		},
		{
			name: "database error",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO auth.sessions`)).
					WithArgs(session.SessionID, session.UserID, session.DeviceInfo, session.CreatedAt, session.ExpiresAt).
					WillReturnError(errors.New("db error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.CreateSession(context.Background(), session)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to create session")
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_ListSessions(t *testing.T) {
	repo, mock := NewMock(t)

	userID := "user1"
	sessions := []models.Session{
		{
			SessionID:  "session1",
			DeviceInfo: "device1",
			CreatedAt:  time.Now(),
			ExpiresAt:  time.Now().Add(24 * time.Hour),
		},
		{
			SessionID:  "session2",
			DeviceInfo: "device2",
			CreatedAt:  time.Now(),
			ExpiresAt:  time.Now().Add(24 * time.Hour),
		},
	}

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
		expected  []models.Session
	}{
		{
			name: "successful session list",
			mockSetup: func() {
				rows := pgxmock.NewRows([]string{"session_id", "device_info", "created_at", "expires_at"}).
					AddRow(sessions[0].SessionID, sessions[0].DeviceInfo, sessions[0].CreatedAt, sessions[0].ExpiresAt).
					AddRow(sessions[1].SessionID, sessions[1].DeviceInfo, sessions[1].CreatedAt, sessions[1].ExpiresAt)
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT session_id, device_info, created_at, expires_at
					FROM auth.sessions
					WHERE user_id = $1`)).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectErr: false,
			expected:  sessions,
		},
		{
			name: "no sessions found",
			mockSetup: func() {
				rows := pgxmock.NewRows([]string{"session_id", "device_info", "created_at", "expires_at"})
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

			result, err := repo.ListSessions(context.Background(), userID)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to list sessions")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_TerminateSession(t *testing.T) {
	repo, mock := NewMock(t)

	sessionID := "session1"

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
	}{
		{
			name: "successful session termination",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`
					DELETE FROM auth.sessions
					WHERE session_id = $1`)).
					WithArgs(sessionID).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			expectErr: false,
		},
		{
			name: "database error",
			mockSetup: func() {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM auth.sessions`)).
					WithArgs(sessionID).
					WillReturnError(errors.New("db error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.TerminateSession(context.Background(), sessionID)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to terminate session")
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
