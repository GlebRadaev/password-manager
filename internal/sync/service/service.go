package service

//go:generate mockgen -destination=service_mock.go -source=service.go -package=service
import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/GlebRadaev/password-manager/internal/sync/models"
	"github.com/GlebRadaev/password-manager/pkg/data"
	"github.com/google/uuid"
)

var (
	ErrSyncFailed       = errors.New("sync failed")
	ErrConflictNotFound = errors.New("conflict not found")
)

type Repo interface {
	GetConflict(ctx context.Context, conflictID string) (models.Conflict, error)
	AddConflicts(ctx context.Context, conflicts []models.Conflict) error
	GetUnresolvedConflicts(ctx context.Context, userID string) ([]models.Conflict, error)
	ResolveConflict(ctx context.Context, conflictID string) error
	DeleteConflicts(ctx context.Context, conflictIDs []string) error
}

type DataClient interface {
	UpdateData(ctx context.Context, userID string, entry models.ClientData) error
	ListData(ctx context.Context, userID string) ([]models.DataEntry, error)
	BatchProcess(ctx context.Context, userID string, operations []*data.DataOperation) (string, error)
	CreateAddOperation(userID string, entry models.ClientData) *data.DataOperation
	CreateUpdateOperation(userID string, entry models.ClientData) *data.DataOperation
	CreateDeleteOperation(userID, dataID string) *data.DataOperation
}

type Service struct {
	conflictRepo Repo
	dataClient   DataClient
}

func New(conflictRepo Repo, dataClient DataClient) *Service {
	return &Service{
		conflictRepo: conflictRepo,
		dataClient:   dataClient,
	}
}

func (s *Service) SyncData(ctx context.Context, userID string, clientData []models.ClientData) ([]models.Conflict, error) {
	log.Info().Msg("Starting SyncData")

	var conflicts []models.Conflict
	var operations []*data.DataOperation

	serverEntries, err := s.dataClient.ListData(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list server data")
		return nil, fmt.Errorf("failed to list server data: %w", err)
	}

	serverDataMap := make(map[string]models.DataEntry)
	for _, entry := range serverEntries {
		serverDataMap[entry.DataID] = entry
	}

	for _, clientEntry := range clientData {
		serverEntry, exists := serverDataMap[clientEntry.DataID]

		switch clientEntry.Operation {
		case models.Delete:
			if exists {
				operations = append(operations, s.dataClient.CreateDeleteOperation(userID, clientEntry.DataID))
			}

		case models.Add:
			if exists {
				if serverEntry.UpdatedAt.Before(clientEntry.UpdatedAt) {
					if !bytes.Equal(serverEntry.Data, clientEntry.Data) {
						operations = append(operations, s.dataClient.CreateUpdateOperation(userID, clientEntry))
					}
				} else if !bytes.Equal(serverEntry.Data, clientEntry.Data) {
					conflict := models.Conflict{
						ID:         uuid.NewString(),
						UserID:     userID,
						DataID:     clientEntry.DataID,
						ClientData: clientEntry.Data,
						ServerData: serverEntry.Data,
					}
					conflicts = append(conflicts, conflict)
				}
			} else {
				operations = append(operations, s.dataClient.CreateAddOperation(userID, clientEntry))
			}

		case models.Update:
			if exists {
				if serverEntry.UpdatedAt.Before(clientEntry.UpdatedAt) {
					if !bytes.Equal(serverEntry.Data, clientEntry.Data) {
						operations = append(operations, s.dataClient.CreateUpdateOperation(userID, clientEntry))
					}
				} else if !bytes.Equal(serverEntry.Data, clientEntry.Data) {
					conflict := models.Conflict{
						ID:         uuid.NewString(),
						UserID:     userID,
						DataID:     clientEntry.DataID,
						ClientData: clientEntry.Data,
						ServerData: serverEntry.Data,
					}
					conflicts = append(conflicts, conflict)
				}
			} else {
				operations = append(operations, s.dataClient.CreateAddOperation(userID, clientEntry))
			}
		}
	}

	if len(conflicts) > 0 {
		if err := s.conflictRepo.AddConflicts(ctx, conflicts); err != nil {
			log.Error().Err(err).Msg("Failed to add conflicts")
			return nil, fmt.Errorf("failed to add conflicts: %w", err)
		}
	} else {
		log.Info().Msg("No conflicts detected")
	}

	if len(operations) > 0 {
		_, err := s.dataClient.BatchProcess(ctx, userID, operations)
		if err != nil {
			if len(conflicts) > 0 {
				rollbackErr := s.rollbackInsertedConflicts(ctx, conflicts)
				if rollbackErr != nil {
					log.Error().Err(rollbackErr).Msg("Failed to rollback conflicts")
					return nil, fmt.Errorf("failed to process batch operations and rollback conflicts: %w (rollback error: %v)", err, rollbackErr)
				}
			}
			return nil, fmt.Errorf("failed to process batch operations: %w", err)
		}
	} else {
		log.Info().Msg("No operations to process")
	}

	log.Info().Msg("SyncData completed successfully")
	return conflicts, nil
}

func (s *Service) ResolveConflict(ctx context.Context, conflictID string, strategy models.ResolutionStrategy) error {
	conflict, err := s.conflictRepo.GetConflict(ctx, conflictID)
	if err != nil {
		if errors.Is(err, ErrConflictNotFound) {
			return ErrConflictNotFound
		}
		return fmt.Errorf("failed to get conflict: %w", err)
	}

	var resolvedData []byte
	switch strategy {
	case models.UseClientVersion:
		resolvedData = conflict.ClientData
	case models.UseServerVersion:
		resolvedData = conflict.ServerData
	case models.MergeVersions:
		resolvedData, err = mergeData(conflict.ClientData, conflict.ServerData, models.Text)
		if err != nil {
			return fmt.Errorf("failed to merge data: %w", err)
		}
	default:
		return fmt.Errorf("unknown resolution strategy")
	}

	if err := s.dataClient.UpdateData(ctx, conflict.UserID, models.ClientData{
		DataID: conflict.DataID,
		Data:   resolvedData,
	}); err != nil {
		return fmt.Errorf("failed to update data: %w", err)
	}

	if err := s.conflictRepo.ResolveConflict(ctx, conflictID); err != nil {
		return fmt.Errorf("failed to resolve conflict: %w", err)
	}

	return nil
}

func (s *Service) ListConflicts(ctx context.Context, userID string) ([]models.Conflict, error) {
	return s.conflictRepo.GetUnresolvedConflicts(ctx, userID)
}

func (s *Service) rollbackInsertedConflicts(ctx context.Context, conflicts []models.Conflict) error {
	if len(conflicts) == 0 {
		return nil
	}

	conflictIDs := make([]string, 0, len(conflicts))
	for _, conflict := range conflicts {
		conflictIDs = append(conflictIDs, conflict.ID)
	}

	err := s.conflictRepo.DeleteConflicts(ctx, conflictIDs)
	if err != nil {
		return fmt.Errorf("failed to rollback inserted conflicts: %w", err)
	}

	return nil
}

func mergeData(clientData, serverData []byte, dataType models.DataType) ([]byte, error) {
	switch dataType {
	case models.Text:
		return []byte(string(clientData) + "\n" + string(serverData)), nil
	case models.Binary:
		return clientData, nil
	case models.LoginPassword:
		return clientData, nil
	case models.Card:
		return clientData, nil
	default:
		return nil, fmt.Errorf("unsupported data type for merging: %v", dataType)
	}
}
