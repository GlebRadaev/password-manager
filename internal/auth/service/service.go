package service

//go:generate mockgen -destination=service_mock.go -source=service.go -package=service
import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/GlebRadaev/password-manager/internal/auth/models"
	"github.com/GlebRadaev/password-manager/internal/common/pg"
	"github.com/google/uuid"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidOTP         = errors.New("invalid OTP")
)

type Repo interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	CheckExists(ctx context.Context, username, email string) (bool, error)
	CreateOTP(ctx context.Context, otp models.OTP) error
	GetOTP(ctx context.Context, userID string, otpCode string) (models.OTP, error)
}

type TxManager interface {
	Begin(ctx context.Context, fn pg.TransactionalFn) (err error)
}

type AuthManager interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
	GenerateToken(userID string) (string, time.Time, error)
	ValidateToken(token string) (string, error)
	GenerateOTP() (string, time.Time, error)
	ValidateOTP(storedOTPCode string, storedExpiresAt time.Time, providedOTPCode string) (bool, error)
}

type Service struct {
	repo        Repo
	txManager   TxManager
	authManager AuthManager
}

func New(repo Repo, txManager TxManager, authManager AuthManager) *Service {
	return &Service{
		repo:        repo,
		txManager:   txManager,
		authManager: authManager,
	}
}

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

func (s *Service) Login(ctx context.Context, username, password string) (string, int64, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get user: %w", err)
	}

	if err := s.authManager.Compare(user.PasswordHash, password); err != nil {
		return "", 0, ErrInvalidCredentials
	}

	token, expiresAt, err := s.authManager.GenerateToken(user.ID)
	if err != nil {
		return "", 0, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, expiresAt.Unix(), nil
}

func (s *Service) ValidateToken(ctx context.Context, token string) (bool, string, error) {
	userID, err := s.authManager.ValidateToken(token)
	if err != nil {
		return false, "", fmt.Errorf("failed to validate token: %w", err)
	}

	return true, userID, nil
}

func (s *Service) GenerateOTP(ctx context.Context, userID string) (string, error) {
	otpCode, expiresAt, err := s.authManager.GenerateOTP()
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	otp := models.OTP{
		UserID:    userID,
		OTPCode:   otpCode,
		ExpiresAt: expiresAt,
	}

	if err := s.repo.CreateOTP(ctx, otp); err != nil {
		return "", fmt.Errorf("failed to create OTP: %w", err)
	}

	return otpCode, nil
}

func (s *Service) ValidateOTP(ctx context.Context, userID, otpCode string) (bool, error) {
	otp, err := s.repo.GetOTP(ctx, userID, otpCode)
	if err != nil {
		return false, fmt.Errorf("failed to get OTP: %w", err)
	}

	valid, err := s.authManager.ValidateOTP(otp.OTPCode, otp.ExpiresAt, otpCode)
	if err != nil {
		return false, ErrInvalidOTP
	}

	return valid, nil
}
