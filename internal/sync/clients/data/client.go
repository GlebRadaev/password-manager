package client

import (
	"context"
	"fmt"
	"time"

	"github.com/GlebRadaev/password-manager/internal/common/app"
	"github.com/GlebRadaev/password-manager/internal/sync/models"
	"github.com/GlebRadaev/password-manager/pkg/data"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	errPrefix = "data service error"
)

type Client struct {
	client data.DataServiceClient
}

func NewClient(app *app.GRPCClient) (*Client, error) {
	conn, err := grpc.NewClient(app.Endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("%s failed to create gRPC connection: %w", errPrefix, err)
	}

	return &Client{client: data.NewDataServiceClient(conn)}, nil
}

func (c *Client) UpdateData(ctx context.Context, userID string, entry models.ClientData) error {
	_, err := c.client.UpdateData(ctx, &data.UpdateDataRequest{
		UserId:   userID,
		DataId:   entry.DataID,
		Data:     entry.Data,
		Metadata: toProtoMetadata(entry.Metadata),
	})
	if err != nil {
		return fmt.Errorf("%s failed to update data: %w", errPrefix, err)
	}

	return nil
}

func (c *Client) ListData(ctx context.Context, userID string) ([]models.DataEntry, error) {
	resp, err := c.client.ListData(ctx, &data.ListDataRequest{UserId: userID})
	if err != nil {
		return nil, fmt.Errorf("%s failed to list data: %w", errPrefix, err)
	}

	var entries []models.DataEntry
	for _, entry := range resp.Entries {
		entries = append(entries, models.DataEntry{
			DataID:    entry.DataId,
			Type:      toModelDataType(entry.Type),
			Data:      entry.Data,
			CreatedAt: time.Unix(entry.CreatedAt, 0),
			UpdatedAt: time.Unix(entry.UpdatedAt, 0),
			Metadata:  toModelMetadata(entry.Metadata),
		})
	}

	return entries, nil
}

func (c *Client) BatchProcess(ctx context.Context, userID string, operations []*data.DataOperation) (string, error) {
	_, err := c.client.BatchProcess(ctx, &data.BatchProcessRequest{
		UserId:     userID,
		Operations: operations,
	})
	if err != nil {
		return "", fmt.Errorf("%s failed to process batch operations: %w", errPrefix, err)
	}

	return "Batch operations processed successfully", nil
}

func (c *Client) CreateAddOperation(userID string, entry models.ClientData) *data.DataOperation {
	return &data.DataOperation{
		Operation: &data.DataOperation_Add{
			Add: &data.AddDataRequest{
				UserId:   userID,
				Type:     toProtoDataType(entry.Type),
				Data:     entry.Data,
				Metadata: toProtoMetadata(entry.Metadata),
			},
		},
	}
}

func (c *Client) CreateUpdateOperation(userID string, entry models.ClientData) *data.DataOperation {
	return &data.DataOperation{
		Operation: &data.DataOperation_Update{
			Update: &data.UpdateDataRequest{
				UserId:   userID,
				DataId:   entry.DataID,
				Data:     entry.Data,
				Metadata: toProtoMetadata(entry.Metadata),
			},
		},
	}
}

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

func toModelDataType(dt data.DataType) models.DataType {
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

func toModelMetadata(protoMetadata []*data.Metadata) []models.Metadata {
	var metadata []models.Metadata
	for _, m := range protoMetadata {
		metadata = append(metadata, models.Metadata{
			Key:   m.Key,
			Value: m.Value,
		})
	}
	return metadata
}

func toProtoDataType(dt models.DataType) data.DataType {
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

func toProtoMetadata(metadata []models.Metadata) []*data.Metadata {
	var protoMetadata []*data.Metadata
	for _, m := range metadata {
		protoMetadata = append(protoMetadata, &data.Metadata{
			Key:   m.Key,
			Value: m.Value,
		})
	}
	return protoMetadata
}
