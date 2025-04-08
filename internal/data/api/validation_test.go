package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/GlebRadaev/password-manager/internal/data/api"
	"github.com/GlebRadaev/password-manager/pkg/data"
)

func TestValidateAddDataRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  *data.AddDataRequest
		expected error
	}{
		{
			name: "valid request",
			request: &data.AddDataRequest{
				UserId: "550e8400-e29b-41d4-a716-446655440000",
				Data:   []byte("test data"),
				Type:   data.DataType_LOGIN_PASSWORD,
			},
			expected: nil,
		},
		{
			name: "invalid user_id",
			request: &data.AddDataRequest{
				UserId: "invalid-uuid",
				Data:   []byte("test data"),
				Type:   data.DataType_LOGIN_PASSWORD,
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidUserID.Error()),
		},
		{
			name: "empty data",
			request: &data.AddDataRequest{
				UserId: "550e8400-e29b-41d4-a716-446655440000",
				Data:   []byte{},
				Type:   data.DataType_LOGIN_PASSWORD,
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrEmptyData.Error()),
		},
		{
			name: "invalid data type",
			request: &data.AddDataRequest{
				UserId: "550e8400-e29b-41d4-a716-446655440000",
				Data:   []byte("test data"),
				Type:   555,
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidDataType.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := api.ValidateAddDataRequest(tt.request)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestValidateUpdateDataRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  *data.UpdateDataRequest
		expected error
	}{
		{
			name: "valid request",
			request: &data.UpdateDataRequest{
				UserId: "550e8400-e29b-41d4-a716-446655440000",
				DataId: "550e8400-e29b-41d4-a716-446655440001",
				Data:   []byte("updated data"),
			},
			expected: nil,
		},
		{
			name: "invalid user_id",
			request: &data.UpdateDataRequest{
				UserId: "invalid-uuid",
				DataId: "550e8400-e29b-41d4-a716-446655440001",
				Data:   []byte("updated data"),
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidUserID.Error()),
		},
		{
			name: "invalid data_id",
			request: &data.UpdateDataRequest{
				UserId: "550e8400-e29b-41d4-a716-446655440000",
				DataId: "invalid-uuid",
				Data:   []byte("updated data"),
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidDataID.Error()),
		},
		{
			name: "empty data",
			request: &data.UpdateDataRequest{
				UserId: "550e8400-e29b-41d4-a716-446655440000",
				DataId: "550e8400-e29b-41d4-a716-446655440001",
				Data:   []byte{},
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrEmptyData.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := api.ValidateUpdateDataRequest(tt.request)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestValidateDeleteDataRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  *data.DeleteDataRequest
		expected error
	}{
		{
			name: "valid request",
			request: &data.DeleteDataRequest{
				UserId: "550e8400-e29b-41d4-a716-446655440000",
				DataId: "550e8400-e29b-41d4-a716-446655440001",
			},
			expected: nil,
		},
		{
			name: "invalid user_id",
			request: &data.DeleteDataRequest{
				UserId: "invalid-uuid",
				DataId: "550e8400-e29b-41d4-a716-446655440001",
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidUserID.Error()),
		},
		{
			name: "invalid data_id",
			request: &data.DeleteDataRequest{
				UserId: "550e8400-e29b-41d4-a716-446655440000",
				DataId: "invalid-uuid",
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidDataID.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := api.ValidateDeleteDataRequest(tt.request)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestValidateListDataRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  *data.ListDataRequest
		expected error
	}{
		{
			name: "valid request",
			request: &data.ListDataRequest{
				UserId: "550e8400-e29b-41d4-a716-446655440000",
			},
			expected: nil,
		},
		{
			name: "invalid user_id",
			request: &data.ListDataRequest{
				UserId: "invalid-uuid",
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidUserID.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := api.ValidateListDataRequest(tt.request)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestValidateGetDataRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  *data.GetDataRequest
		expected error
	}{
		{
			name: "valid request",
			request: &data.GetDataRequest{
				UserId: "550e8400-e29b-41d4-a716-446655440000",
				DataId: "550e8400-e29b-41d4-a716-446655440001",
			},
			expected: nil,
		},
		{
			name: "invalid user_id",
			request: &data.GetDataRequest{
				UserId: "invalid-uuid",
				DataId: "550e8400-e29b-41d4-a716-446655440001",
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidUserID.Error()),
		},
		{
			name: "invalid data_id",
			request: &data.GetDataRequest{
				UserId: "550e8400-e29b-41d4-a716-446655440000",
				DataId: "invalid-uuid",
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidDataID.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := api.ValidateGetDataRequest(tt.request)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestValidateBatchProcessRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  *data.BatchProcessRequest
		expected error
	}{
		{
			name: "valid request",
			request: &data.BatchProcessRequest{
				UserId: "550e8400-e29b-41d4-a716-446655440000",
				Operations: []*data.DataOperation{
					{
						Operation: &data.DataOperation_Add{
							Add: &data.AddDataRequest{
								UserId: "550e8400-e29b-41d4-a716-446655440000",
								Data:   []byte("test"),
								Type:   data.DataType_LOGIN_PASSWORD,
							},
						},
					},
				},
			},
			expected: nil,
		},
		{
			name: "invalid user_id",
			request: &data.BatchProcessRequest{
				UserId: "invalid-uuid",
				Operations: []*data.DataOperation{
					{
						Operation: &data.DataOperation_Add{
							Add: &data.AddDataRequest{
								UserId: "invalid-uuid",
								Data:   []byte("test"),
								Type:   data.DataType_LOGIN_PASSWORD,
							},
						},
					},
				},
			},
			expected: status.Errorf(codes.InvalidArgument, api.ErrInvalidUserID.Error()),
		},
		{
			name: "empty operations",
			request: &data.BatchProcessRequest{
				UserId:     "550e8400-e29b-41d4-a716-446655440000",
				Operations: []*data.DataOperation{},
			},
			expected: status.Errorf(codes.InvalidArgument, "operations must contain between 1 and 100 items"),
		},
		{
			name: "too many operations",
			request: &data.BatchProcessRequest{
				UserId:     "550e8400-e29b-41d4-a716-446655440000",
				Operations: make([]*data.DataOperation, 101),
			},
			expected: status.Errorf(codes.InvalidArgument, "operations must contain between 1 and 100 items"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := api.ValidateBatchProcessRequest(tt.request)
			assert.Equal(t, tt.expected, err)
		})
	}
}
