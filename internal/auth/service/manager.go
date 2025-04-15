// Package service provides core authentication business logic.
// Handles password hashing, JWT tokens, and OTP generation/validation.
package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/GlebRadaev/password-manager/internal/auth/config"
)

// Manager handles security operations for authentication.
type Manager struct {
	secretKey       string
	tokenExpiration time.Duration
	otpExpiration   time.Duration
}

// NewManager creates a new Manager with security configurations.
func NewManager(cfg config.LocalConfig) *Manager {
	return &Manager{
		secretKey:       cfg.SecretKey,
		tokenExpiration: cfg.TokenExpiration,
		otpExpiration:   cfg.OTPExpiration,
	}
}

// Hash generates bcrypt hash of the password.
func (m *Manager) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// Compare verifies password against its hash.
func (m *Manager) Compare(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("password comparison failed: %w", err)
	}
	return nil
}

// GenerateToken creates JWT access token for user.
func (m *Manager) GenerateToken(userID string) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.tokenExpiration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
	})

	tokenString, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// ValidateToken verifies JWT token and extracts userID.
func (m *Manager) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("invalid user_id in token")
		}

		return userIDStr, nil
	}

	return "", errors.New("invalid token")
}

// GenerateOTP creates random 6-character OTP code.
func (m *Manager) GenerateOTP() (string, time.Time, error) {
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate OTP: %w", err)
	}

	otpCode := hex.EncodeToString(bytes)
	expiresAt := time.Now().Add(m.otpExpiration)

	return otpCode, expiresAt, nil
}

// ValidateOTP checks if provided OTP matches stored one and isn't expired.
func (m *Manager) ValidateOTP(storedOTPCode string, storedExpiresAt time.Time, providedOTPCode string) (bool, error) {
	if time.Now().After(storedExpiresAt) {
		return false, errors.New("OTP has expired")
	}

	if providedOTPCode != storedOTPCode {
		return false, errors.New("invalid OTP code")
	}

	return true, nil
}

// GenerateRefreshToken creates long-lived refresh token.
func (m *Manager) GenerateRefreshToken(userID string) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.tokenExpiration * 2)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
		"type":    "refresh",
	})

	tokenString, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// ValidateRefreshToken verifies refresh token and extracts userID.
func (m *Manager) ValidateRefreshToken(refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse refresh token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"] != "refresh" {
			return "", errors.New("invalid token type")
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("invalid user_id in token")
		}

		return userIDStr, nil
	}

	return "", errors.New("invalid refresh token")
}
