package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/GlebRadaev/password-manager/internal/auth/models"
	"github.com/GlebRadaev/password-manager/internal/common/pg"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db pg.Database
}

func New(db pg.Database) *Repository {
	return &Repository{db: db}
}

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

func (r *Repository) CreateOTP(ctx context.Context, otp models.OTP) error {
	query := `
		INSERT INTO auth.otp_codes (user_id, otp_code, expires_at)
		VALUES ($1, $2, $3)`

	_, err := r.db.Exec(ctx, query, otp.UserID, otp.OTPCode, otp.ExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to create OTP: %w", err)
	}

	return nil
}

func (r *Repository) GetOTP(ctx context.Context, userID string, otpCode string) (models.OTP, error) {
	query := `
		SELECT user_id, otp_code, expires_at
		FROM auth.otp_codes
		WHERE user_id = $1 AND otp_code = $2`

	var otp models.OTP
	err := r.db.QueryRow(ctx, query, userID, otpCode).Scan(
		&otp.UserID,
		&otp.OTPCode,
		&otp.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.OTP{}, fmt.Errorf("OTP not found: %w", err)
		}
		return models.OTP{}, fmt.Errorf("failed to get OTP: %w", err)
	}

	return otp, nil
}
