package service

//go:generate mockgen -destination=service_mock.go -source=service.go -package=service
import (
	"context"
	"errors"
	"fmt"

	"github.com/GlebRadaev/password-manager/internal/common/pg"
	"github.com/GlebRadaev/password-manager/internal/data/models"
	"github.com/GlebRadaev/password-manager/internal/data/repo"
	"github.com/GlebRadaev/password-manager/pkg/data"
	"github.com/google/uuid"
)

var ErrDataNotFound = errors.New("data not found")

type Repo interface {
	AddList(ctx context.Context, entries []models.DataEntry) ([]string, error)
	UpdateData(ctx context.Context, entry models.DataEntry) error
	DeleteList(ctx context.Context, userID string, dataIDs []string) error
	ListData(ctx context.Context, userID string) ([]models.DataEntry, error)
}

type TxManager interface {
	Begin(ctx context.Context, fn pg.TransactionalFn) (err error)
}

type Service struct {
	repo      Repo
	txManager TxManager
}

func New(repo Repo, txManager TxManager) *Service {
	return &Service{
		repo:      repo,
		txManager: txManager,
	}
}

func (s *Service) AddData(ctx context.Context, entry models.DataEntry) (string, error) {
	ids, err := s.repo.AddList(ctx, []models.DataEntry{entry})
	if err != nil {
		return "", fmt.Errorf("failed to add data: %w", err)
	}
	if len(ids) == 0 {
		return "", fmt.Errorf("no data was added")
	}
	return ids[0], nil
}

func (s *Service) UpdateData(ctx context.Context, entry models.DataEntry) error {
	err := s.repo.UpdateData(ctx, entry)
	if errors.Is(err, repo.ErrDataNotFound) {
		return ErrDataNotFound
	}
	return err
}

func (s *Service) DeleteData(ctx context.Context, userID, dataID string) error {
	err := s.repo.DeleteList(ctx, userID, []string{dataID})
	if errors.Is(err, repo.ErrDataNotFound) {
		return ErrDataNotFound
	}
	return err
}

func (s *Service) ListData(ctx context.Context, userID string) ([]models.DataEntry, error) {
	return s.repo.ListData(ctx, userID)
}

func (s *Service) BatchProcess(ctx context.Context, userID string, operations []*data.DataOperation) ([]*data.DataOperationResult, error) {
	// TODO: to model structure
	var results []*data.DataOperationResult
	var addEntries []models.DataEntry
	var updateEntries []models.DataEntry
	var deleteIDs []string

	for i, op := range operations {
		switch op := op.Operation.(type) {
		case *data.DataOperation_Add:
			addEntries = append(addEntries, models.DataEntry{
				ID:       uuid.New().String(),
				UserID:   userID,
				Type:     toModelDataType(op.Add.Type),
				Data:     op.Add.Data,
				Metadata: toModelMetadata(op.Add.Metadata),
			})

		case *data.DataOperation_Update:
			updateEntries = append(updateEntries, models.DataEntry{
				UserID:   userID,
				ID:       op.Update.DataId,
				Data:     op.Update.Data,
				Metadata: toModelMetadata(op.Update.Metadata),
			})

		case *data.DataOperation_Delete:
			deleteIDs = append(deleteIDs, op.Delete.DataId)

		default:
			return nil, fmt.Errorf("unknown operation type at index %d", i)
		}
	}

	err := s.txManager.Begin(ctx, func(ctx context.Context) error {
		if len(addEntries) > 0 {
			dataIDs, err := s.repo.AddList(ctx, addEntries)
			if err != nil {
				return fmt.Errorf("failed to add data: %w", err)
			}
			if len(dataIDs) == 0 {
				return fmt.Errorf("no data was added")
			}
			for _, dataID := range dataIDs {
				results = append(results, &data.DataOperationResult{
					Result: &data.DataOperationResult_Add{
						Add: &data.AddDataResponse{DataId: dataID},
					},
				})
			}
		}

		if len(updateEntries) > 0 {
			for _, entry := range updateEntries {
				err := s.repo.UpdateData(ctx, entry)
				if err != nil {
					if errors.Is(err, repo.ErrDataNotFound) {
						return ErrDataNotFound
					}
					return fmt.Errorf("failed to update data: %w", err)
				}
				results = append(results, &data.DataOperationResult{
					Result: &data.DataOperationResult_Update{
						Update: &data.UpdateDataResponse{Message: "Data updated successfully"},
					},
				})
			}
		}

		if len(deleteIDs) > 0 {
			err := s.repo.DeleteList(ctx, userID, deleteIDs)
			if err != nil {
				if errors.Is(err, repo.ErrDataNotFound) {
					return ErrDataNotFound
				}
				return fmt.Errorf("failed to delete data: %w", err)
			}
			for range deleteIDs {
				results = append(results, &data.DataOperationResult{
					Result: &data.DataOperationResult_Delete{
						Delete: &data.DeleteDataResponse{Message: "Data deleted successfully"},
					},
				})
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to process batch operations: %w", err)
	}

	return results, nil
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
