package service

//go:generate mockgen -destination=service_mock.go -source=service.go -package=service
import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/GlebRadaev/password-manager/internal/sync/models"
)

type Repo interface {
	SaveChanges(ctx context.Context, changes []models.DataChange) error
	GetChanges(ctx context.Context, userID string, lastSyncTime int64) ([]models.DataChange, error)
	GetChangesByDataID(ctx context.Context, userID, dataID string) ([]models.DataChange, error)
}

type Service struct {
	repo Repo
}

func New(repo Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetChanges(ctx context.Context, userID string, lastSyncTime int64) ([]models.DataChange, error) {
	return s.repo.GetChanges(ctx, userID, lastSyncTime)
}

func (s *Service) PushChanges(ctx context.Context, userID string, changes []models.DataChange) error {
	for i := range changes {
		changes[i].ID = uuid.NewString()
		changes[i].Timestamp = time.Now()
	}

	var resolvedChanges []models.DataChange
	for _, change := range changes {
		existingChanges, err := s.repo.GetChangesByDataID(ctx, userID, change.DataID)
		if err != nil {
			return fmt.Errorf("failed to get existing changes: %w", err)
		}

		resolvedChange := s.ResolveConflicts(existingChanges, change)
		resolvedChanges = append(resolvedChanges, resolvedChange)
	}

	return s.repo.SaveChanges(ctx, resolvedChanges)
}

func (s *Service) ResolveConflicts(existingChanges []models.DataChange, newChange models.DataChange) models.DataChange {
	if len(existingChanges) == 0 {
		return newChange
	}

	latestChange := existingChanges[0]
	for _, change := range existingChanges {
		if change.Timestamp.After(latestChange.Timestamp) {
			latestChange = change
		}
	}

	if newChange.Timestamp.After(latestChange.Timestamp) {
		return newChange
	}

	return latestChange
}
