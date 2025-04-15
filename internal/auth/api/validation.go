// Package api provides validation functions for auth service requests.
package api

import (
	"errors"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/GlebRadaev/password-manager/pkg/auth"
)

// Common validation errors
var (
	ErrInvalidUsername = errors.New("username must be between 3 and 20 characters and contain only letters, numbers, and underscores")
	ErrInvalidPassword = errors.New("password must be between 8 and 50 characters")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidToken    = errors.New("token must be at least 10 characters long")
	ErrInvalidUserID   = errors.New("user_id must be a valid UUID")
	ErrInvalidOTPCode  = errors.New("OTP code must be exactly 6 characters long")
)

// ValidateRegisterRequest validates RegisterRequest fields.
func ValidateRegisterRequest(req *auth.RegisterRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "Username":
			return status.Error(codes.InvalidArgument, ErrInvalidUsername.Error())
		case "Password":
			return status.Error(codes.InvalidArgument, ErrInvalidPassword.Error())
		case "Email":
			return status.Error(codes.InvalidArgument, ErrInvalidEmail.Error())
		default:
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
		}
	}
	return nil
}

// ValidateLoginRequest validates LoginRequest fields.
func ValidateLoginRequest(req *auth.LoginRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "Username":
			return status.Error(codes.InvalidArgument, ErrInvalidUsername.Error())
		case "Password":
			return status.Error(codes.InvalidArgument, ErrInvalidPassword.Error())
		default:
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err.Error())
		}
	}
	return nil
}

// ValidateTokenRequest validates ValidateTokenRequest fields.
func ValidateTokenRequest(req *auth.ValidateTokenRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "Token":
			return status.Error(codes.InvalidArgument, ErrInvalidToken.Error())
		default:
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err.Error())
		}
	}
	return nil
}

// ValidateGenerateOTPRequest validates GenerateOTPRequest fields.
func ValidateGenerateOTPRequest(req *auth.GenerateOTPRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "UserID":
			return status.Error(codes.InvalidArgument, ErrInvalidUserID.Error())
		default:
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err.Error())
		}
	}
	return nil
}

// ValidateOTPRequest validates ValidateOTPRequest fields.
func ValidateOTPRequest(req *auth.ValidateOTPRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "UserID":
			return status.Error(codes.InvalidArgument, ErrInvalidUserID.Error())
		case "OtpCode":
			return status.Error(codes.InvalidArgument, ErrInvalidOTPCode.Error())
		default:
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err.Error())
		}
	}
	return nil
}

// ValidateRefreshTokenRequest validates RefreshTokenRequest fields.
func ValidateRefreshTokenRequest(req *auth.RefreshTokenRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "RefreshToken":
			return status.Error(codes.InvalidArgument, ErrInvalidToken.Error())
		default:
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err.Error())
		}
	}
	return nil
}

// extractFieldFromError extracts field name from validation error string.
func extractFieldFromError(errStr string) string {
	parts := strings.Split(errStr, ".")
	if len(parts) < 2 {
		return ""
	}
	fieldPart := parts[1]
	fieldName := strings.Split(fieldPart, ":")[0]
	return strings.TrimSpace(fieldName)
}
