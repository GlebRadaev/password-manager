package service_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/GlebRadaev/password-manager/internal/auth/config"
	"github.com/GlebRadaev/password-manager/internal/auth/service"
)

func TestManager_Hash(t *testing.T) {
	m := service.NewManager(config.LocalConfig{})

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"success", "strongpassword123", false},
		{"empty password", "", false},
		{"max length password", string(make([]byte, 72)), false}, // 72 bytes is bcrypt max
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := m.Hash(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, hash)
		})
	}
}

func TestManager_Compare(t *testing.T) {
	m := service.NewManager(config.LocalConfig{})
	password := "testpassword"
	hash, _ := m.Hash(password)

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		wantErr        bool
	}{
		{"success", hash, password, false},
		{"wrong password", hash, "wrongpassword", true},
		{"empty password", hash, "", true},
		{"invalid hash", "invalidhash", password, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := m.Compare(tt.hashedPassword, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestManager_GenerateToken(t *testing.T) {
	cfg := config.LocalConfig{
		SecretKey:       "testsecret",
		TokenExpiration: time.Hour,
	}
	m := service.NewManager(cfg)

	t.Run("success", func(t *testing.T) {
		userID := "user123"
		token, expiresAt, err := m.GenerateToken(userID)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.WithinDuration(t, time.Now().Add(cfg.TokenExpiration), expiresAt, time.Second)

		parsedUserID, err := m.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, parsedUserID)
	})
}

func TestManager_ValidateToken(t *testing.T) {
	cfg := config.LocalConfig{
		SecretKey:       "testsecret",
		TokenExpiration: time.Hour,
	}
	m := service.NewManager(cfg)
	validToken, _, _ := m.GenerateToken("user123")

	tests := []struct {
		name       string
		token      string
		wantUserID string
		wantErr    bool
		errMsg     string
	}{
		{"valid token", validToken, "user123", false, ""},
		{"invalid token", "invalid.token", "", true, "failed to parse token"},
		{"empty token", "", "", true, "failed to parse token"},
		{"expired token", generateExpiredToken(cfg.SecretKey), "", true, "token is expired"},
		{"wrong secret", generateTokenWithWrongSecret(), "", true, "signature is invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := m.ValidateToken(tt.token)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantUserID, userID)
		})
	}
}

func TestManager_GenerateOTP(t *testing.T) {
	cfg := config.LocalConfig{
		OTPExpiration: time.Minute * 5,
	}
	m := service.NewManager(cfg)

	t.Run("success", func(t *testing.T) {
		otp, expiresAt, err := m.GenerateOTP()
		require.NoError(t, err)
		assert.NotEmpty(t, otp)
		assert.WithinDuration(t, time.Now().Add(cfg.OTPExpiration), expiresAt, time.Second)
	})
}

func TestManager_ValidateOTP(t *testing.T) {
	m := service.NewManager(config.LocalConfig{})

	tests := []struct {
		name           string
		storedOTP      string
		storedExpires  time.Time
		providedOTP    string
		wantValid      bool
		wantErr        bool
		wantErrMessage string
	}{
		{"valid", "123456", time.Now().Add(time.Hour), "123456", true, false, ""},
		{"expired", "123456", time.Now().Add(-time.Hour), "123456", false, true, "expired"},
		{"invalid code", "123456", time.Now().Add(time.Hour), "654321", false, true, "invalid"},
		{"empty stored", "", time.Now().Add(time.Hour), "123456", false, true, "invalid"},
		{"empty provided", "123456", time.Now().Add(time.Hour), "", false, true, "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := m.ValidateOTP(tt.storedOTP, tt.storedExpires, tt.providedOTP)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrMessage)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantValid, valid)
		})
	}
}

func TestManager_GenerateRefreshToken(t *testing.T) {
	cfg := config.LocalConfig{
		SecretKey:       "testsecret",
		TokenExpiration: time.Hour,
	}
	m := service.NewManager(cfg)

	t.Run("success", func(t *testing.T) {
		userID := "user123"
		token, expiresAt, err := m.GenerateRefreshToken(userID)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.WithinDuration(t, time.Now().Add(2*cfg.TokenExpiration), expiresAt, time.Second)

		parsedUserID, err := m.ValidateRefreshToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, parsedUserID)
	})
}

func TestManager_ValidateRefreshToken(t *testing.T) {
	cfg := config.LocalConfig{
		SecretKey:       "testsecret",
		TokenExpiration: time.Hour,
	}
	m := service.NewManager(cfg)
	validToken, _, _ := m.GenerateRefreshToken("user123")

	tests := []struct {
		name       string
		token      string
		wantUserID string
		wantErr    bool
		errMsg     string
	}{
		{"valid token", validToken, "user123", false, ""},
		{"invalid token", "invalid.token", "", true, "failed to parse refresh token"},
		{"empty token", "", "", true, "failed to parse refresh token"},
		{"expired token", generateExpiredRefreshToken(cfg.SecretKey), "", true, "token is expired"},
		{"wrong secret", generateRefreshTokenWithWrongSecret(), "", true, "signature is invalid"},
		{"access token", generateAccessToken(cfg.SecretKey), "", true, "invalid token type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := m.ValidateRefreshToken(tt.token)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantUserID, userID)
		})
	}
}

// Helper functions
func generateExpiredToken(secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "user123",
		"exp":     time.Now().Add(-time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func generateTokenWithWrongSecret() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "user123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("wrongsecret"))
	return tokenString
}

func generateExpiredRefreshToken(secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "user123",
		"exp":     time.Now().Add(-time.Hour).Unix(),
		"type":    "refresh",
	})
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func generateRefreshTokenWithWrongSecret() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "user123",
		"exp":     time.Now().Add(time.Hour).Unix(),
		"type":    "refresh",
	})
	tokenString, _ := token.SignedString([]byte("wrongsecret"))
	return tokenString
}

func generateAccessToken(secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "user123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}
