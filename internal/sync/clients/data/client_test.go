package client_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/require"

	client "github.com/GlebRadaev/password-manager/internal/sync/clients/data"
	"github.com/GlebRadaev/password-manager/internal/sync/models"
	"github.com/GlebRadaev/password-manager/pkg/data"
)

func TestUpdateData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := data.NewMockDataServiceClient(ctrl)
	c := &client.Client{Client: mockClient}

	ctx := context.Background()
	userID := "user1"
	entry := models.ClientData{
		DataID: "data1",
		Data:   []byte("password123"),
		Metadata: []models.Metadata{
			{Key: "service", Value: "github"},
		},
	}

	mockClient.EXPECT().UpdateData(
		ctx,
		&data.UpdateDataRequest{
			UserId:   userID,
			DataId:   entry.DataID,
			Data:     entry.Data,
			Metadata: client.ToProtoMetadata(entry.Metadata),
		},
	).Return(nil, nil)

	err := c.UpdateData(ctx, userID, entry)
	require.NoError(t, err)
}

func TestUpdateData_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := data.NewMockDataServiceClient(ctrl)
	c := &client.Client{Client: mockClient}

	ctx := context.Background()
	userID := "user1"
	entry := models.ClientData{DataID: "data1"}

	expectedErr := errors.New("gRPC error")
	mockClient.EXPECT().UpdateData(
		ctx,
		gomock.Any(),
	).Return(nil, expectedErr)

	err := c.UpdateData(ctx, userID, entry)
	require.Error(t, err)
	require.Contains(t, err.Error(), "data service error")
}

func TestListData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := data.NewMockDataServiceClient(ctrl)
	c := &client.Client{Client: mockClient}

	ctx := context.Background()
	userID := "user1"

	now := time.Now()
	mockClient.EXPECT().ListData(
		ctx,
		&data.ListDataRequest{UserId: userID},
	).Return(&data.ListDataResponse{
		Entries: []*data.DataEntry{
			{
				DataId:    "data1",
				Type:      data.DataType_LOGIN_PASSWORD,
				Data:      []byte("login:pass"),
				CreatedAt: now.Unix(),
				UpdatedAt: now.Unix(),
				Metadata: []*data.Metadata{
					{Key: "service", Value: "github"},
				},
			},
		},
	}, nil)

	entries, err := c.ListData(ctx, userID)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, "data1", entries[0].DataID)
	require.Equal(t, models.LoginPassword, entries[0].Type)
	require.Equal(t, now.Unix(), entries[0].CreatedAt.Unix())
}

func TestListData_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := data.NewMockDataServiceClient(ctrl)
	c := &client.Client{Client: mockClient}

	ctx := context.Background()
	userID := "user1"

	expectedErr := errors.New("gRPC list error")
	mockClient.EXPECT().ListData(
		ctx,
		&data.ListDataRequest{UserId: userID},
	).Return(nil, expectedErr)

	entries, err := c.ListData(ctx, userID)
	require.Nil(t, entries)
	require.Error(t, err)
	require.Contains(t, err.Error(), "data service error")
	require.Contains(t, err.Error(), "failed to list data")
}

func TestBatchProcess_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := data.NewMockDataServiceClient(ctrl)
	c := &client.Client{Client: mockClient}

	ctx := context.Background()
	userID := "user1"
	operations := []*data.DataOperation{
		c.CreateAddOperation(userID, models.ClientData{DataID: "data1"}),
	}

	mockClient.EXPECT().BatchProcess(
		ctx,
		&data.BatchProcessRequest{
			UserId:     userID,
			Operations: operations,
		},
	).Return(nil, nil)

	msg, err := c.BatchProcess(ctx, userID, operations)
	require.NoError(t, err)
	require.Equal(t, "Batch operations processed successfully", msg)
}

func TestBatchProcess_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := data.NewMockDataServiceClient(ctrl)
	c := &client.Client{Client: mockClient}

	ctx := context.Background()
	userID := "user1"
	operations := []*data.DataOperation{
		c.CreateAddOperation(userID, models.ClientData{DataID: "data1"}),
	}

	expectedErr := errors.New("gRPC batch error")
	mockClient.EXPECT().BatchProcess(
		ctx,
		&data.BatchProcessRequest{
			UserId:     userID,
			Operations: operations,
		},
	).Return(nil, expectedErr)

	msg, err := c.BatchProcess(ctx, userID, operations)
	require.Empty(t, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "data service error")
	require.Contains(t, err.Error(), "failed to process batch operations")
}

func TestCreateAddOperation(t *testing.T) {
	c := &client.Client{}
	userID := "user1"
	entry := models.ClientData{
		DataID: "data1",
		Type:   models.LoginPassword,
		Data:   []byte("secret"),
	}

	op := c.CreateAddOperation(userID, entry)
	require.NotNil(t, op.GetAdd())
	require.Equal(t, userID, op.GetAdd().UserId)
	require.Equal(t, data.DataType_LOGIN_PASSWORD, op.GetAdd().Type)
}

func TestCreateUpdateOperation(t *testing.T) {
	c := &client.Client{}
	userID := "user1"
	entry := models.ClientData{
		DataID: "data1",
		Type:   models.LoginPassword,
		Data:   []byte("updated-secret"),
		Metadata: []models.Metadata{
			{Key: "service", Value: "github"},
		},
	}

	op := c.CreateUpdateOperation(userID, entry)
	require.NotNil(t, op.GetUpdate())
	require.Equal(t, userID, op.GetUpdate().UserId)
	require.Equal(t, entry.DataID, op.GetUpdate().DataId)
	require.Equal(t, entry.Data, op.GetUpdate().Data)
	require.Len(t, op.GetUpdate().Metadata, 1)
}

func TestCreateDeleteOperation(t *testing.T) {
	c := &client.Client{}
	userID := "user1"
	dataID := "data1"

	op := c.CreateDeleteOperation(userID, dataID)
	require.NotNil(t, op.GetDelete())
	require.Equal(t, dataID, op.GetDelete().DataId)
}
func TestToProtoDataType(t *testing.T) {
	require.Equal(t, data.DataType_LOGIN_PASSWORD, client.ToProtoDataType(models.LoginPassword))
	require.Equal(t, data.DataType_TEXT, client.ToProtoDataType(models.Text))
	require.Equal(t, data.DataType_BINARY, client.ToProtoDataType(models.Binary))
	require.Equal(t, data.DataType_CARD, client.ToProtoDataType(models.Card))
}

func TestToModelMetadata(t *testing.T) {
	protoMetadata := []*data.Metadata{
		{Key: "key1", Value: "value1"},
	}
	modelMetadata := client.ToModelMetadata(protoMetadata)
	require.Len(t, modelMetadata, 1)
	require.Equal(t, "key1", modelMetadata[0].Key)
}

func TestToModelDataType(t *testing.T) {
	require.Equal(t, models.LoginPassword, client.ToModelDataType(data.DataType_LOGIN_PASSWORD))
	require.Equal(t, models.Text, client.ToModelDataType(data.DataType_TEXT))
	require.Equal(t, models.Binary, client.ToModelDataType(data.DataType_BINARY))
	require.Equal(t, models.Card, client.ToModelDataType(data.DataType_CARD))
	require.Equal(t, models.LoginPassword, client.ToModelDataType(data.DataType(999)))
}

func TestToProtoMetadata(t *testing.T) {
	modelMetadata := []models.Metadata{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}

	protoMetadata := client.ToProtoMetadata(modelMetadata)
	require.Len(t, protoMetadata, 2)
	require.Equal(t, "key1", protoMetadata[0].Key)
	require.Equal(t, "value2", protoMetadata[1].Value)
}
