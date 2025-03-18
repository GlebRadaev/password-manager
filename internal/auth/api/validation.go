package api

import (
	"errors"
	"strings"

	"github.com/GlebRadaev/password-manager/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidUsername = errors.New("username must be between 3 and 20 characters and contain only letters, numbers, and underscores")
	ErrInvalidPassword = errors.New("password must be between 8 and 50 characters")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidToken    = errors.New("token must be at least 10 characters long")
	ErrInvalidUserID   = errors.New("user_id must be a valid UUID")
	ErrInvalidOTPCode  = errors.New("OTP code must be exactly 6 characters long")
)

func ValidateRegisterRequest(req *auth.RegisterRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "Username":
			return status.Errorf(codes.InvalidArgument, ErrInvalidUsername.Error())
		case "Password":
			return status.Errorf(codes.InvalidArgument, ErrInvalidPassword.Error())
		case "Email":
			return status.Errorf(codes.InvalidArgument, ErrInvalidEmail.Error())
		default:
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err.Error())
		}
	}
	return nil
}

func ValidateLoginRequest(req *auth.LoginRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "Username":
			return status.Errorf(codes.InvalidArgument, ErrInvalidUsername.Error())
		case "Password":
			return status.Errorf(codes.InvalidArgument, ErrInvalidPassword.Error())
		default:
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err.Error())
		}
	}
	return nil
}

func ValidateTokenRequest(req *auth.ValidateTokenRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "Token":
			return status.Errorf(codes.InvalidArgument, ErrInvalidToken.Error())
		default:
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err.Error())
		}
	}
	return nil
}

func ValidateGenerateOTPRequest(req *auth.GenerateOTPRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "UserID":
			return status.Errorf(codes.InvalidArgument, ErrInvalidUserID.Error())
		default:
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err.Error())
		}
	}
	return nil
}

func ValidateOTPRequest(req *auth.ValidateOTPRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "UserID":
			return status.Errorf(codes.InvalidArgument, ErrInvalidUserID.Error())
		case "OtpCode":
			return status.Errorf(codes.InvalidArgument, ErrInvalidOTPCode.Error())
		default:
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err.Error())
		}
	}
	return nil
}

func extractFieldFromError(errStr string) string {
	parts := strings.Split(errStr, ".")
	if len(parts) < 2 {
		return ""
	}
	fieldPart := parts[1]
	fieldName := strings.Split(fieldPart, ":")[0]
	return strings.TrimSpace(fieldName)
}
