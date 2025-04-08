// Package repo provides data access layer for authentication service.
// Handles all database operations for users, OTP codes and sessions.
package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/GlebRadaev/password-manager/internal/auth/models"
	"github.com/GlebRadaev/password-manager/internal/common/pg"
)

// Repository manages persistence operations for auth data.
type Repository struct {
	db pg.Database
}

// New creates a new Repository instance with given database connection.
func New(db pg.Database) *Repository {
	return &Repository{db: db}
}

// CreateUser persists new user in database.
func (r *Repository) CreateUser(ctx context.Context, user models.User) error {
	query := `
		INSERT INTO auth.users (user_id, username, password_hash, email, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(ctx, query, user.ID, user.Username, user.PasswordHash, user.Email, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByUsername retrieves user by username.
func (r *Repository) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	query := `
		SELECT user_id, username, password_hash, email, created_at
		FROM auth.users
		WHERE username = $1`

	var user models.User
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("user not found: %w", err)
		}
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// CheckExists verifies if user with given username or email exists.
func (r *Repository) CheckExists(ctx context.Context, username, email string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM auth.users
			WHERE username = $1 OR email = $2
		)`

	var exists bool
	err := r.db.QueryRow(ctx, query, username, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

// CreateOTP stores new OTP code in database.
func (r *Repository) CreateOTP(ctx context.Context, otp models.OTP) error {
	query := `
		INSERT INTO auth.otp_codes (user_id, otp_code, expires_at, device_id)
		VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(ctx, query, otp.UserID, otp.OTPCode, otp.ExpiresAt, otp.DeviceID)
	if err != nil {
		return fmt.Errorf("failed to create OTP: %w", err)
	}

	return nil
}

// GetOTP retrieves OTP code for given user and device.
func (r *Repository) GetOTP(ctx context.Context, userID, otpCode, deviceID string) (models.OTP, error) {
	query := `
		SELECT user_id, otp_code, expires_at, device_id
		FROM auth.otp_codes
		WHERE user_id = $1 AND otp_code = $2 AND device_id = $3`

	var otp models.OTP
	err := r.db.QueryRow(ctx, query, userID, otpCode, deviceID).Scan(
		&otp.UserID,
		&otp.OTPCode,
		&otp.ExpiresAt,
		&otp.DeviceID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.OTP{}, fmt.Errorf("OTP not found: %w", err)
		}
		return models.OTP{}, fmt.Errorf("failed to get OTP: %w", err)
	}

	return otp, nil
}

// CreateSession stores new user session in database.
func (r *Repository) CreateSession(ctx context.Context, session models.Session) error {
	query := `
		INSERT INTO auth.sessions (session_id, user_id, device_info, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(ctx, query, session.SessionID, session.UserID, session.DeviceInfo, session.CreatedAt, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// ListSessions returns all active sessions for given user.
func (r *Repository) ListSessions(ctx context.Context, userID string) ([]models.Session, error) {
	query := `
		SELECT session_id, device_info, created_at, expires_at
		FROM auth.sessions
		WHERE user_id = $1`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var session models.Session
		if err := rows.Scan(&session.SessionID, &session.DeviceInfo, &session.CreatedAt, &session.ExpiresAt); err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// TerminateSession removes session by its ID.
func (r *Repository) TerminateSession(ctx context.Context, sessionID string) error {
	query := `
		DELETE FROM auth.sessions
		WHERE session_id = $1`

	_, err := r.db.Exec(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to terminate session: %w", err)
	}

	return nil
}
