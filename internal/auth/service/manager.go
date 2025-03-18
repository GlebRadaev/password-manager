package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/GlebRadaev/password-manager/internal/auth/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type Manager struct {
	secretKey       string
	tokenExpiration time.Duration
	otpExpiration   time.Duration
}

func NewManager(cfg config.LocalConfig) *Manager {
	return &Manager{
		secretKey:       cfg.SecretKey,
		tokenExpiration: cfg.TokenExpiration,
		otpExpiration:   cfg.OTPExpiration,
	}
}

func (m *Manager) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

func (m *Manager) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

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

func (m *Manager) GenerateOTP() (string, time.Time, error) {
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate OTP: %w", err)
	}

	otpCode := hex.EncodeToString(bytes)
	expiresAt := time.Now().Add(m.otpExpiration)

	return otpCode, expiresAt, nil
}

func (m *Manager) ValidateOTP(storedOTPCode string, storedExpiresAt time.Time, providedOTPCode string) (bool, error) {
	if time.Now().After(storedExpiresAt) {
		return false, errors.New("OTP has expired")
	}

	if providedOTPCode != storedOTPCode {
		return false, errors.New("invalid OTP code")
	}

	return true, nil
}

func (m *Manager) GenerateTOTPSecret(userID string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "PasswordManager",
		AccountName: userID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP secret: %w", err)
	}
	return key.Secret(), nil
}

func (m *Manager) ValidateTOTP(secret, code string) (bool, error) {
	valid := totp.Validate(code, secret)
	if !valid {
		return false, errors.New("invalid TOTP code")
	}
	return true, nil
}
