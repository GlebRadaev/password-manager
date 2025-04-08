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
	ErrInvalidUserID   = errors.New("user_id must be a valid UUID")
	ErrInvalidDataID   = errors.New("data_id must be a valid UUID")
	ErrEmptyData       = errors.New("data cannot be empty")
	ErrInvalidDataType = errors.New("invalid data type")
	ErrDataNotFound    = errors.New("data not found")
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
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
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
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err.Error())
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
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
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
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
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
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
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
			return status.Error(codes.InvalidArgument, "operations must contain between 1 and 100 items")
		default:
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
		}
	}
	return nil
}

// extractFieldFromError extracts field name from validation error string
func extractFieldFromError(errStr string) string {
	parts := strings.Split(errStr, ".")
	if len(parts) < 2 {
		return ""
	}
	fieldPart := parts[1]
	fieldName := strings.Split(fieldPart, ":")[0]
	return strings.TrimSpace(fieldName)
}
