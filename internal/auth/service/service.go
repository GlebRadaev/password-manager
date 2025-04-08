// Package service implements core authentication business logic.
// Handles user registration, login, sessions, and token management.
package service

//go:generate mockgen -destination=service_mock.go -source=service.go -package=service
import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/GlebRadaev/password-manager/internal/auth/models"
	"github.com/GlebRadaev/password-manager/internal/common/pg"
)

// Common service errors
var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidOTP         = errors.New("invalid OTP")
	ErrSessionNotFound    = errors.New("session not found")
)

// Repo defines repository interface for persistence operations
type Repo interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	CheckExists(ctx context.Context, username, email string) (bool, error)
	CreateOTP(ctx context.Context, otp models.OTP) error
	GetOTP(ctx context.Context, userID, otpCode, deviceID string) (models.OTP, error)
	CreateSession(ctx context.Context, session models.Session) error
	ListSessions(ctx context.Context, userID string) ([]models.Session, error)
	TerminateSession(ctx context.Context, sessionID string) error
}

// TxManager handles database transactions
type TxManager interface {
	Begin(ctx context.Context, fn pg.TransactionalFn) (err error)
}

// AuthManager handles security operations
type AuthManager interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
	GenerateToken(userID string) (string, time.Time, error)
	ValidateToken(token string) (string, error)
	GenerateOTP() (string, time.Time, error)
	ValidateOTP(storedOTPCode string, storedExpiresAt time.Time, providedOTPCode string) (bool, error)
	GenerateRefreshToken(userID string) (string, time.Time, error)
	ValidateRefreshToken(refreshToken string) (string, error)
}

// Service implements authentication business logic
type Service struct {
	repo        Repo
	txManager   TxManager
	authManager AuthManager
}

// New creates a new Service instance with dependencies
func New(repo Repo, txManager TxManager, authManager AuthManager) *Service {
	return &Service{
		repo:        repo,
		txManager:   txManager,
		authManager: authManager,
	}
}

// Register creates new user with hashed password
func (s *Service) Register(ctx context.Context, username, password, email string) (string, error) {
	exists, err := s.repo.CheckExists(ctx, username, email)
	if err != nil {
		return "", fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return "", ErrUserExists
	}

	userID := uuid.NewString()
	passwordHash, err := s.authManager.Hash(password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	user := models.User{
		ID:           userID,
		Username:     username,
		PasswordHash: passwordHash,
		Email:        email,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}

// Login authenticates user and returns tokens
func (s *Service) Login(ctx context.Context, username, password string) (string, string, int64, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to get user: %w", err)
	}

	if err := s.authManager.Compare(user.PasswordHash, password); err != nil {
		return "", "", 0, ErrInvalidCredentials
	}

	accessToken, accessExpiresAt, err := s.authManager.GenerateToken(user.ID)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, refreshExpiresAt, err := s.authManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	deviceInfo := "CLI"
	session := models.Session{
		SessionID:  uuid.NewString(),
		UserID:     user.ID,
		DeviceInfo: deviceInfo,
		CreatedAt:  time.Now(),
		ExpiresAt:  refreshExpiresAt,
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return "", "", 0, fmt.Errorf("failed to create session: %w", err)
	}

	return accessToken, refreshToken, accessExpiresAt.Unix(), nil
}

// ValidateToken verifies JWT token and returns user ID
func (s *Service) ValidateToken(ctx context.Context, token string) (bool, string, error) {
	userID, err := s.authManager.ValidateToken(token)
	if err != nil {
		return false, "", fmt.Errorf("failed to validate token: %w", err)
	}

	return true, userID, nil
}

// GenerateOTP creates and stores one-time password
func (s *Service) GenerateOTP(ctx context.Context, userID, deviceID string) (string, error) {
	otpCode, expiresAt, err := s.authManager.GenerateOTP()
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	otpRecord := models.OTP{
		UserID:    userID,
		OTPCode:   otpCode,
		ExpiresAt: expiresAt,
		DeviceID:  deviceID,
	}

	if err := s.repo.CreateOTP(ctx, otpRecord); err != nil {
		return "", fmt.Errorf("failed to create OTP: %w", err)
	}

	return otpCode, nil
}

// ValidateOTP verifies one-time password
func (s *Service) ValidateOTP(ctx context.Context, userID, otpCode, deviceID string) (bool, error) {
	otp, err := s.repo.GetOTP(ctx, userID, otpCode, deviceID)
	if err != nil {
		return false, fmt.Errorf("failed to get OTP: %w", err)
	}

	valid, err := s.authManager.ValidateOTP(otp.OTPCode, otp.ExpiresAt, otpCode)
	if !valid || err != nil {
		return false, ErrInvalidOTP
	}

	return true, nil
}

// ListSessions returns all active sessions for user
func (s *Service) ListSessions(ctx context.Context, userID string) ([]models.Session, error) {
	sessions, err := s.repo.ListSessions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	return sessions, nil
}

// TerminateSession ends specified session
func (s *Service) TerminateSession(ctx context.Context, sessionID string) error {
	err := s.repo.TerminateSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to terminate session: %w", err)
	}
	return nil
}

// RefreshToken generates new access token using refresh token
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (string, string, int64, error) {
	userID, err := s.authManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid refresh token: %w", err)
	}

	accessToken, accessExpiresAt, err := s.authManager.GenerateToken(userID)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, refreshExpiresAt, err := s.authManager.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	_ = refreshExpiresAt

	return accessToken, refreshToken, accessExpiresAt.Unix(), nil
}
