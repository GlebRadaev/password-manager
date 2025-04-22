// Package api provides request validation utilities for data service
package api

import (
	"errors"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/GlebRadaev/password-manager/pkg/data"
)

// Common validation errors
var (
	ErrInvalidUserID    = errors.New("user_id must be a valid UUID")
	ErrInvalidDataID    = errors.New("data_id must be a valid UUID")
	ErrEmptyData        = errors.New("data cannot be empty")
	ErrInvalidDataType  = errors.New("invalid data type")
	ErrDataNotFound     = errors.New("data not found")
	ErrValidationFailed = "validation failed: %v"
)

// Operation limits
const (
	MinOperations = 1
	MaxOperations = 100
)

// Error field parsing constants
const (
	errorSplitParts   = 2
	errorFieldParts   = 2
	errorFieldDivider = ":"
	errorFieldPrefix  = "."
)

// ValidateAddDataRequest validates AddDataRequest fields
func ValidateAddDataRequest(req *data.AddDataRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "UserId":
			return status.Error(codes.InvalidArgument, ErrInvalidUserID.Error())
		case "Data":
			return status.Error(codes.InvalidArgument, ErrEmptyData.Error())
		case "Type":
			return status.Error(codes.InvalidArgument, ErrInvalidDataType.Error())
		default:
			return status.Errorf(codes.InvalidArgument, ErrValidationFailed, err)
		}
	}
	return nil
}

// ValidateUpdateDataRequest validates UpdateDataRequest fields
func ValidateUpdateDataRequest(req *data.UpdateDataRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "UserId":
			return status.Error(codes.InvalidArgument, ErrInvalidUserID.Error())
		case "DataId":
			return status.Error(codes.InvalidArgument, ErrInvalidDataID.Error())
		case "Data":
			return status.Error(codes.InvalidArgument, ErrEmptyData.Error())
		default:
			return status.Errorf(codes.InvalidArgument, ErrValidationFailed, err.Error())
		}
	}
	return nil
}

// ValidateDeleteDataRequest validates DeleteDataRequest fields
func ValidateDeleteDataRequest(req *data.DeleteDataRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "UserId":
			return status.Error(codes.InvalidArgument, ErrInvalidUserID.Error())
		case "DataId":
			return status.Error(codes.InvalidArgument, ErrInvalidDataID.Error())
		default:
			return status.Errorf(codes.InvalidArgument, ErrValidationFailed, err)
		}
	}
	return nil
}

// ValidateListDataRequest validates ListDataRequest fields
func ValidateListDataRequest(req *data.ListDataRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "UserId":
			return status.Error(codes.InvalidArgument, ErrInvalidUserID.Error())
		default:
			return status.Errorf(codes.InvalidArgument, ErrValidationFailed, err)
		}
	}
	return nil
}

// ValidateGetDataRequest validates GetDataRequest fields
func ValidateGetDataRequest(req *data.GetDataRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "UserId":
			return status.Error(codes.InvalidArgument, ErrInvalidUserID.Error())
		case "DataId":
			return status.Error(codes.InvalidArgument, ErrInvalidDataID.Error())
		default:
			return status.Errorf(codes.InvalidArgument, ErrValidationFailed, err)
		}
	}
	return nil
}

// ValidateBatchProcessRequest validates BatchProcessRequest fields
func ValidateBatchProcessRequest(req *data.BatchProcessRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "UserId":
			return status.Error(codes.InvalidArgument, ErrInvalidUserID.Error())
		case "Operations":
			return status.Errorf(
				codes.InvalidArgument,
				"operations must contain between %d and %d items",
				MinOperations,
				MaxOperations,
			)
		default:
			return status.Errorf(codes.InvalidArgument, ErrValidationFailed, err)
		}
	}
	return nil
}

// extractFieldFromError extracts field name from validation error string
func extractFieldFromError(errStr string) string {
	parts := strings.SplitN(errStr, errorFieldPrefix, errorSplitParts)
	if len(parts) < errorSplitParts {
		return ""
	}

	fieldParts := strings.SplitN(parts[1], errorFieldDivider, errorFieldParts)
	if len(fieldParts) < errorFieldParts {
		return strings.TrimSpace(parts[1])
	}

	return strings.TrimSpace(fieldParts[0])
}
