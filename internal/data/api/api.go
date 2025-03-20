package api

//go:generate mockgen -destination=api_mock.go -source=api.go -package=api
import (
	"context"
	"errors"
	"time"

	"github.com/GlebRadaev/password-manager/internal/data/models"
	"github.com/GlebRadaev/password-manager/internal/data/service"
	"github.com/GlebRadaev/password-manager/pkg/data"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	AddData(ctx context.Context, entry models.DataEntry) (string, error)
	UpdateData(ctx context.Context, entry models.DataEntry) error
	DeleteData(ctx context.Context, userID, dataID string) error
	ListData(ctx context.Context, userID string) ([]models.DataEntry, error)
	BatchProcess(ctx context.Context, userID string, operations []*data.DataOperation) ([]*data.DataOperationResult, error)
}

type Api struct {
	data.UnimplementedDataServiceServer
	srv Service
}

func New(srv Service) *Api {
	return &Api{srv: srv}
}

func (s *Api) AddData(ctx context.Context, req *data.AddDataRequest) (*data.AddDataResponse, error) {
	if err := ValidateAddDataRequest(req); err != nil {
		return nil, err
	}
	dataID, err := s.srv.AddData(ctx, models.DataEntry{
		UserID:    req.UserId,
		Type:      models.DataType(models.Add),
		Data:      req.Data,
		Metadata:  toModelMetadata(req.Metadata),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, FromError(err, "add data")
	}
	return &data.AddDataResponse{DataId: dataID}, nil
}

func (s *Api) UpdateData(ctx context.Context, req *data.UpdateDataRequest) (*data.UpdateDataResponse, error) {
	if err := ValidateUpdateDataRequest(req); err != nil {
		return nil, err
	}
	err := s.srv.UpdateData(ctx, models.DataEntry{
		UserID:    req.UserId,
		ID:        req.DataId,
		Type:      models.DataType(models.Update),
		Data:      req.Data,
		Metadata:  toModelMetadata(req.Metadata),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, FromError(err, "update data")
	}
	return &data.UpdateDataResponse{Message: "Data updated successfully"}, nil
}

func (s *Api) DeleteData(ctx context.Context, req *data.DeleteDataRequest) (*data.DeleteDataResponse, error) {
	if err := ValidateDeleteDataRequest(req); err != nil {
		return nil, err
	}
	err := s.srv.DeleteData(ctx, req.UserId, req.DataId)
	if err != nil {
		return nil, FromError(err, "delete data")
	}
	return &data.DeleteDataResponse{Message: "Data deleted successfully"}, nil
}

func (s *Api) ListData(ctx context.Context, req *data.ListDataRequest) (*data.ListDataResponse, error) {
	if err := ValidateListDataRequest(req); err != nil {
		return nil, err
	}
	entries, err := s.srv.ListData(ctx, req.UserId)
	if err != nil {
		return nil, FromError(err, "list data")
	}

	var protoEntries []*data.DataEntry
	for _, entry := range entries {
		protoEntries = append(protoEntries, &data.DataEntry{
			DataId:    entry.ID,
			Type:      toProtoDataType(entry.Type),
			Data:      entry.Data,
			Metadata:  toProtoMetadata(entry.Metadata),
			CreatedAt: entry.CreatedAt.Unix(),
			UpdatedAt: entry.UpdatedAt.Unix(),
		})
	}
	return &data.ListDataResponse{Entries: protoEntries}, nil
}

func (s *Api) BatchProcess(ctx context.Context, req *data.BatchProcessRequest) (*data.BatchProcessResponse, error) {
	if err := ValidateBatchProcessRequest(req); err != nil {
		return nil, err
	}
	results, err := s.srv.BatchProcess(ctx, req.UserId, req.Operations)
	if err != nil {
		return nil, FromError(err, "batch process")
	}
	return &data.BatchProcessResponse{Results: results}, nil
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

func FromError(err error, operation string) error {
	if err == nil {
		return nil
	}

	var code codes.Code
	switch {
	case errors.Is(err, service.ErrDataNotFound):
		code = codes.NotFound
	default:
		code = codes.Internal
	}

	return status.Errorf(code, "failed to %s: %v", operation, err)
}
