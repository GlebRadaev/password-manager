package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/GlebRadaev/password-manager/internal/auth/api"
	"github.com/GlebRadaev/password-manager/pkg/auth"
)

func TestValidateRegisterRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  *auth.RegisterRequest
		expected error
	}{
		{
			name: "valid request",
			request: &auth.RegisterRequest{
				Username: "valid_user",
				Password: "valid_password123",
				Email:    "valid@example.com",
			},
			expected: nil,
		},
		{
			name: "invalid username - too short",
			request: &auth.RegisterRequest{
				Username: "a",
				Password: "valid_password123",
				Email:    "valid@example.com",
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidUsername.Error()),
		},
		{
			name: "invalid password - too short",
			request: &auth.RegisterRequest{
				Username: "valid_user",
				Password: "short",
				Email:    "valid@example.com",
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidPassword.Error()),
		},
		{
			name: "invalid email",
			request: &auth.RegisterRequest{
				Username: "valid_user",
				Password: "valid_password123",
				Email:    "invalid-email",
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidEmail.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := api.ValidateRegisterRequest(tt.request)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestValidateLoginRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  *auth.LoginRequest
		expected error
	}{
		{
			name: "valid request",
			request: &auth.LoginRequest{
				Username: "valid_user",
				Password: "valid_password123",
			},
			expected: nil,
		},
		{
			name: "invalid username",
			request: &auth.LoginRequest{
				Username: "a",
				Password: "valid_password123",
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidUsername.Error()),
		},
		{
			name: "invalid password",
			request: &auth.LoginRequest{
				Username: "valid_user",
				Password: "short",
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidPassword.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := api.ValidateLoginRequest(tt.request)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestValidateTokenRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  *auth.ValidateTokenRequest
		expected error
	}{
		{
			name: "valid request",
			request: &auth.ValidateTokenRequest{
				Token: "valid_token_long_enough",
			},
			expected: nil,
		},
		{
			name: "invalid token - too short",
			request: &auth.ValidateTokenRequest{
				Token: "short",
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidToken.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := api.ValidateTokenRequest(tt.request)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestValidateGenerateOTPRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  *auth.GenerateOTPRequest
		expected error
	}{
		{
			name: "valid request",
			request: &auth.GenerateOTPRequest{
				UserId: "550e8400-e29b-41d4-a716-446655440000",
			},
			expected: nil,
		},
		{
			name: "invalid user_id",
			request: &auth.GenerateOTPRequest{
				UserId: "not-a-uuid",
			},
			expected: status.Errorf(
				codes.InvalidArgument,
				"validation failed: invalid GenerateOTPRequest.UserId: value must be a valid UUID | caused by: invalid uuid format",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := api.ValidateGenerateOTPRequest(tt.request)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestValidateOTPRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  *auth.ValidateOTPRequest
		expected error
	}{
		{
			name: "valid request",
			request: &auth.ValidateOTPRequest{
				UserId:  "550e8400-e29b-41d4-a716-446655440000",
				OtpCode: "123456",
			},
			expected: nil,
		},
		{
			name: "invalid user_id",
			request: &auth.ValidateOTPRequest{
				UserId:  "not-a-uuid",
				OtpCode: "123456",
			},
			expected: status.Errorf(
				codes.InvalidArgument,
				"validation failed: invalid ValidateOTPRequest.UserId: value must be a valid UUID | caused by: invalid uuid format",
			),
		},
		{
			name: "invalid otp code - too short",
			request: &auth.ValidateOTPRequest{
				UserId:  "550e8400-e29b-41d4-a716-446655440000",
				OtpCode: "123",
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidOTPCode.Error()),
		},
		{
			name: "invalid otp code - too long",
			request: &auth.ValidateOTPRequest{
				UserId:  "550e8400-e29b-41d4-a716-446655440000",
				OtpCode: "1234567",
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidOTPCode.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := api.ValidateOTPRequest(tt.request)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestValidateRefreshTokenRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  *auth.RefreshTokenRequest
		expected error
	}{
		{
			name: "valid request",
			request: &auth.RefreshTokenRequest{
				RefreshToken: "valid_refresh_token_long_enough",
			},
			expected: nil,
		},
		{
			name: "invalid refresh token - too short",
			request: &auth.RefreshTokenRequest{
				RefreshToken: "short",
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidToken.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := api.ValidateRefreshTokenRequest(tt.request)
			assert.Equal(t, tt.expected, err)
		})
	}
}
