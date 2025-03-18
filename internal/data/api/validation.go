package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/GlebRadaev/password-manager/pkg/data"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidUserID   = errors.New("user_id must be a valid UUID")
	ErrInvalidDataID   = errors.New("data_id must be a valid UUID")
	ErrEmptyData       = errors.New("data cannot be empty")
	ErrInvalidDataType = errors.New("invalid data type")
)

func ValidateCreateDataRequest(req *data.CreateDataRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "UserId":
			return status.Errorf(codes.InvalidArgument, ErrInvalidUserID.Error())
		case "Data":
			return status.Errorf(codes.InvalidArgument, ErrEmptyData.Error())
		default:
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err.Error())
		}
	}
	return nil
}

func ValidateGetDataRequest(req *data.GetDataRequest) error {
	if err := req.Validate(); err != nil {
		fmt.Print(err.Error())
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "UserId":
			return status.Errorf(codes.InvalidArgument, ErrInvalidUserID.Error())
		case "DataId":
			return status.Errorf(codes.InvalidArgument, ErrInvalidDataID.Error())
		default:
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err.Error())
		}
	}
	return nil
}

func ValidateUpdateDataRequest(req *data.UpdateDataRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "UserId":
			return status.Errorf(codes.InvalidArgument, ErrInvalidUserID.Error())
		case "DataId":
			return status.Errorf(codes.InvalidArgument, ErrInvalidDataID.Error())
		case "Data":
			return status.Errorf(codes.InvalidArgument, ErrEmptyData.Error())
		default:
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err.Error())
		}
	}
	return nil
}

func ValidateDeleteDataRequest(req *data.DeleteDataRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "UserId":
			return status.Errorf(codes.InvalidArgument, ErrInvalidUserID.Error())
		case "DataId":
			return status.Errorf(codes.InvalidArgument, ErrInvalidDataID.Error())
		default:
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err.Error())
		}
	}
	return nil
}

func ValidateListDataRequest(req *data.ListDataRequest) error {
	if err := req.Validate(); err != nil {
		fieldName := extractFieldFromError(err.Error())
		switch fieldName {
		case "UserId":
			return status.Errorf(codes.InvalidArgument, ErrInvalidUserID.Error())
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
