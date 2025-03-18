package api

import (
	"context"
	"errors"

	"github.com/GlebRadaev/password-manager/internal/data/models"
	"github.com/GlebRadaev/password-manager/internal/data/service"
	"github.com/GlebRadaev/password-manager/pkg/data"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	CreateData(ctx context.Context, userID string, dataType models.DataType, data []byte, metadata map[string]string) (string, error)
	GetData(ctx context.Context, userID, dataID string) (models.Data, error)
	UpdateData(ctx context.Context, userID, dataID string, data []byte, metadata map[string]string) error
	DeleteData(ctx context.Context, userID, dataID string) error
	ListData(ctx context.Context, userID string) ([]models.Data, error)
}

type Api struct {
	data.UnimplementedDataServiceServer
	srv Service
}

func New(srv Service) *Api {
	return &Api{srv: srv}
}

func (a *Api) CreateData(ctx context.Context, req *data.CreateDataRequest) (*data.CreateDataResponse, error) {
	if err := ValidateCreateDataRequest(req); err != nil {
		return nil, err
	}

	dataType := toModelDataType(req.Type)
	dataID, err := a.srv.CreateData(ctx, req.UserId, dataType, req.Data, req.Metadata)
	if err != nil {
		if errors.Is(err, service.ErrDataNotFound) {
			return nil, status.Errorf(codes.NotFound, "data not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to create data: %v", err)
	}

	return &data.CreateDataResponse{DataId: dataID}, nil
}

func (a *Api) GetData(ctx context.Context, req *data.GetDataRequest) (*data.GetDataResponse, error) {
	if err := ValidateGetDataRequest(req); err != nil {
		return nil, err
	}

	dataModel, err := a.srv.GetData(ctx, req.UserId, req.DataId)
	if err != nil {
		if errors.Is(err, service.ErrDataNotFound) {
			return nil, status.Errorf(codes.NotFound, "data not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get data: %v", err)
	}

	return &data.GetDataResponse{
		Type:     data.DataType(dataModel.Type),
		Data:     dataModel.Data,
		Metadata: dataModel.Metadata,
	}, nil
}

func (a *Api) UpdateData(ctx context.Context, req *data.UpdateDataRequest) (*data.UpdateDataResponse, error) {
	if err := ValidateUpdateDataRequest(req); err != nil {
		return nil, err
	}

	if err := a.srv.UpdateData(ctx, req.UserId, req.DataId, req.Data, req.Metadata); err != nil {
		if errors.Is(err, service.ErrDataNotFound) {
			return nil, status.Errorf(codes.NotFound, "data not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update data: %v", err)
	}

	return &data.UpdateDataResponse{Message: "Data updated successfully"}, nil
}

func (a *Api) DeleteData(ctx context.Context, req *data.DeleteDataRequest) (*data.DeleteDataResponse, error) {
	if err := ValidateDeleteDataRequest(req); err != nil {
		return nil, err
	}

	if err := a.srv.DeleteData(ctx, req.UserId, req.DataId); err != nil {
		if errors.Is(err, service.ErrDataNotFound) {
			return nil, status.Errorf(codes.NotFound, "data not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete data: %v", err)
	}

	return &data.DeleteDataResponse{Message: "Data deleted successfully"}, nil
}

func (a *Api) ListData(ctx context.Context, req *data.ListDataRequest) (*data.ListDataResponse, error) {
	if err := ValidateListDataRequest(req); err != nil {
		return nil, err
	}

	dataList, err := a.srv.ListData(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list data: %v", err)
	}

	response := &data.ListDataResponse{}
	for _, dataModel := range dataList {
		response.Data = append(response.Data, &data.GetDataResponse{
			Type:     data.DataType(dataModel.Type),
			Data:     dataModel.Data,
			Metadata: dataModel.Metadata,
		})
	}

	return response, nil
}

func toModelDataType(dt data.DataType) models.DataType {
	switch dt {
	case data.DataType_LOGIN_PASSWORD:
		return models.LoginPassword
	case data.DataType_TEXT:
		return models.Text
	case data.DataType_BINARY:
		return models.Binary
	case data.DataType_BANK_CARD:
		return models.BankCard
	default:
		return models.LoginPassword
	}
}
