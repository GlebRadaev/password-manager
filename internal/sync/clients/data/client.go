// Package client provides a gRPC client implementation for the data service
package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/GlebRadaev/password-manager/internal/common/app"
	"github.com/GlebRadaev/password-manager/internal/sync/models"
	"github.com/GlebRadaev/password-manager/pkg/data"
)

const (
	errPrefix = "data service error" // Prefix for client error messages
)

// Client wraps the gRPC data service client with helper methods
type Client struct {
	Client data.DataServiceClient
}

// NewClient creates a new gRPC client connection to the data service
// Takes GRPCClient config and returns initialized Client or error
func NewClient(app *app.GRPCClient) (*Client, error) {
	conn, err := grpc.NewClient(app.Endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("%s failed to create gRPC connection: %w", errPrefix, err)
	}

	return &Client{Client: data.NewDataServiceClient(conn)}, nil
}

// UpdateData sends an update request for specific data entry
// Returns error if the operation fails
func (c *Client) UpdateData(ctx context.Context, userID string, entry models.ClientData) error {
	_, err := c.Client.UpdateData(ctx, &data.UpdateDataRequest{
		UserId:   userID,
		DataId:   entry.DataID,
		Data:     entry.Data,
		Metadata: ToProtoMetadata(entry.Metadata),
	})
	if err != nil {
		return fmt.Errorf("%s failed to update data: %w", errPrefix, err)
	}

	return nil
}

// ListData retrieves all data entries for a user
// Returns slice of DataEntry or error
func (c *Client) ListData(ctx context.Context, userID string) ([]models.DataEntry, error) {
	resp, err := c.Client.ListData(ctx, &data.ListDataRequest{UserId: userID})
	if err != nil {
		return nil, fmt.Errorf("%s failed to list data: %w", errPrefix, err)
	}

	var entries []models.DataEntry
	for _, entry := range resp.Entries {
		entries = append(entries, models.DataEntry{
			DataID:    entry.DataId,
			Type:      ToModelDataType(entry.Type),
			Data:      entry.Data,
			CreatedAt: time.Unix(entry.CreatedAt, 0),
			UpdatedAt: time.Unix(entry.UpdatedAt, 0),
			Metadata:  ToModelMetadata(entry.Metadata),
		})
	}

	return entries, nil
}

// BatchProcess executes multiple operations in a single request
// Returns success message or error
func (c *Client) BatchProcess(ctx context.Context, userID string, operations []*data.DataOperation) (string, error) {
	_, err := c.Client.BatchProcess(ctx, &data.BatchProcessRequest{
		UserId:     userID,
		Operations: operations,
	})
	if err != nil {
		return "", fmt.Errorf("%s failed to process batch operations: %w", errPrefix, err)
	}

	return "Batch operations processed successfully", nil
}

// CreateAddOperation constructs an add operation for batch processing
func (c *Client) CreateAddOperation(userID string, entry models.ClientData) *data.DataOperation {
	return &data.DataOperation{
		Operation: &data.DataOperation_Add{
			Add: &data.AddDataRequest{
				UserId:   userID,
				Type:     ToProtoDataType(entry.Type),
				Data:     entry.Data,
				Metadata: ToProtoMetadata(entry.Metadata),
			},
		},
	}
}

// CreateUpdateOperation constructs an update operation for batch processing
func (c *Client) CreateUpdateOperation(userID string, entry models.ClientData) *data.DataOperation {
	return &data.DataOperation{
		Operation: &data.DataOperation_Update{
			Update: &data.UpdateDataRequest{
				UserId:   userID,
				DataId:   entry.DataID,
				Data:     entry.Data,
				Metadata: ToProtoMetadata(entry.Metadata),
			},
		},
	}
}

// CreateDeleteOperation constructs a delete operation for batch processing
func (c *Client) CreateDeleteOperation(userID, dataID string) *data.DataOperation {
	return &data.DataOperation{
		Operation: &data.DataOperation_Delete{
			Delete: &data.DeleteDataRequest{
				UserId: userID,
				DataId: dataID,
			},
		},
	}
}

// ToModelDataType converts protobuf DataType to domain model
func ToModelDataType(dt data.DataType) models.DataType {
	switch dt {
	case data.DataType_LOGIN_PASSWORD:
		return models.LoginPassword
	case data.DataType_TEXT:
		return models.Text
	case data.DataType_BINARY:
		return models.Binary
	case data.DataType_CARD:
		return models.Card
	default:
		return models.LoginPassword
	}
}

// ToModelMetadata converts protobuf Metadata to domain model
func ToModelMetadata(protoMetadata []*data.Metadata) []models.Metadata {
	var metadata []models.Metadata
	for _, m := range protoMetadata {
		metadata = append(metadata, models.Metadata{
			Key:   m.Key,
			Value: m.Value,
		})
	}
	return metadata
}

// ToProtoDataType converts domain DataType to protobuf
func ToProtoDataType(dt models.DataType) data.DataType {
	switch dt {
	case models.LoginPassword:
		return data.DataType_LOGIN_PASSWORD
	case models.Text:
		return data.DataType_TEXT
	case models.Binary:
		return data.DataType_BINARY
	case models.Card:
		return data.DataType_CARD
	default:
		return data.DataType_LOGIN_PASSWORD
	}
}

// ToProtoMetadata converts domain Metadata to protobuf
func ToProtoMetadata(metadata []models.Metadata) []*data.Metadata {
	var protoMetadata []*data.Metadata
	for _, m := range metadata {
		protoMetadata = append(protoMetadata, &data.Metadata{
			Key:   m.Key,
			Value: m.Value,
		})
	}
	return protoMetadata
}
