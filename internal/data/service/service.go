package service

//go:generate mockgen -destination=service_mock.go -source=service.go -package=service
import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/GlebRadaev/password-manager/internal/common/pg"
	"github.com/GlebRadaev/password-manager/internal/data/models"
	"github.com/GlebRadaev/password-manager/internal/data/repo"
)

var (
	ErrDataNotFound = errors.New("data not found")
)

type Repo interface {
	CreateData(ctx context.Context, data models.Data) error
	GetData(ctx context.Context, userID, dataID string) (models.Data, error)
	UpdateData(ctx context.Context, data models.Data) error
	DeleteData(ctx context.Context, userID, dataID string) error
	ListData(ctx context.Context, userID string) ([]models.Data, error)
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

func (s *Service) CreateData(ctx context.Context, userID string, dataType models.DataType, data []byte, metadata map[string]string) (string, error) {
	dataID := uuid.NewString()
	dataModel := models.Data{
		ID:        dataID,
		UserID:    userID,
		Type:      dataType,
		Data:      data,
		Metadata:  metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateData(ctx, dataModel); err != nil {
		return "", fmt.Errorf("failed to create data: %w", err)
	}

	return dataID, nil
}

func (s *Service) GetData(ctx context.Context, userID, dataID string) (models.Data, error) {
	data, err := s.repo.GetData(ctx, userID, dataID)
	if err != nil {
		if errors.Is(err, repo.ErrDataNotFound) {
			return models.Data{}, ErrDataNotFound
		}
		return models.Data{}, fmt.Errorf("failed to get data: %w", err)
	}

	return data, nil
}

func (s *Service) UpdateData(ctx context.Context, userID, dataID string, data []byte, metadata map[string]string) error {
	existingData, err := s.repo.GetData(ctx, userID, dataID)
	if err != nil {
		if errors.Is(err, repo.ErrDataNotFound) {
			return ErrDataNotFound
		}
		return fmt.Errorf("failed to get data for update: %w", err)
	}

	existingData.Data = data
	existingData.Metadata = metadata
	existingData.UpdatedAt = time.Now()

	if err := s.repo.UpdateData(ctx, existingData); err != nil {
		return fmt.Errorf("failed to update data: %w", err)
	}

	return nil
}

func (s *Service) DeleteData(ctx context.Context, userID, dataID string) error {
	if err := s.repo.DeleteData(ctx, userID, dataID); err != nil {
		if errors.Is(err, repo.ErrDataNotFound) {
			return ErrDataNotFound
		}
		return fmt.Errorf("failed to delete data: %w", err)
	}

	return nil
}

func (s *Service) ListData(ctx context.Context, userID string) ([]models.Data, error) {
	dataList, err := s.repo.ListData(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list data: %w", err)
	}

	return dataList, nil
}
