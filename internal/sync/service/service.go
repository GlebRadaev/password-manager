// Package service implements synchronization logic between client and server data.
// Handles conflict detection, resolution and data synchronization operations.
package service

//go:generate mockgen -destination=service_mock.go -source=service.go -package=service
import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/GlebRadaev/password-manager/internal/sync/models"
	"github.com/GlebRadaev/password-manager/pkg/data"
)

// Common synchronization errors
var (
	ErrSyncFailed       = errors.New("sync failed")
	ErrConflictNotFound = errors.New("conflict not found")
)

// Repo defines interface for conflict persistence operations
type Repo interface {
	GetConflict(ctx context.Context, conflictID string) (models.Conflict, error)
	AddConflicts(ctx context.Context, conflicts []models.Conflict) error
	GetUnresolvedConflicts(ctx context.Context, userID string) ([]models.Conflict, error)
	ResolveConflict(ctx context.Context, conflictID string) error
	DeleteConflicts(ctx context.Context, conflictIDs []string) error
}

// DataClient defines interface for data operations
type DataClient interface {
	UpdateData(ctx context.Context, userID string, entry models.ClientData) error
	ListData(ctx context.Context, userID string) ([]models.DataEntry, error)
	BatchProcess(ctx context.Context, userID string, operations []*data.DataOperation) (string, error)
	CreateAddOperation(userID string, entry models.ClientData) *data.DataOperation
	CreateUpdateOperation(userID string, entry models.ClientData) *data.DataOperation
	CreateDeleteOperation(userID, dataID string) *data.DataOperation
}

// Service implements synchronization business logic
type Service struct {
	conflictRepo Repo
	dataClient   DataClient
}

// New creates a new Service instance with dependencies
func New(conflictRepo Repo, dataClient DataClient) *Service {
	return &Service{
		conflictRepo: conflictRepo,
		dataClient:   dataClient,
	}
}

// SyncData synchronizes client data with server and detects conflicts
// Returns list of detected conflicts or error if synchronization fails
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
				op := s.dataClient.CreateDeleteOperation(userID, clientEntry.DataID)
				operations = append(operations, op)
			}

		case models.Add:
			if exists {
				if serverEntry.UpdatedAt.Before(clientEntry.UpdatedAt) {
					if !bytes.Equal(serverEntry.Data, clientEntry.Data) {
						op := s.dataClient.CreateUpdateOperation(userID, clientEntry)
						operations = append(operations, op)
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
				op := s.dataClient.CreateAddOperation(userID, clientEntry)
				operations = append(operations, op)
			}

		case models.Update:
			if exists {
				if serverEntry.UpdatedAt.Before(clientEntry.UpdatedAt) {
					if !bytes.Equal(serverEntry.Data, clientEntry.Data) {
						op := s.dataClient.CreateUpdateOperation(userID, clientEntry)
						operations = append(operations, op)
					}
				} else if !serverEntry.UpdatedAt.Equal(clientEntry.UpdatedAt) {
					if !bytes.Equal(serverEntry.Data, clientEntry.Data) {
						conflict := models.Conflict{
							ID:         uuid.NewString(),
							UserID:     userID,
							DataID:     clientEntry.DataID,
							ClientData: clientEntry.Data,
							ServerData: serverEntry.Data,
						}
						conflicts = append(conflicts, conflict)
					}
				}
			} else {
				op := s.dataClient.CreateAddOperation(userID, clientEntry)
				operations = append(operations, op)
			}
		}
	}

	if len(conflicts) > 0 {
		if err := s.conflictRepo.AddConflicts(ctx, conflicts); err != nil {
			log.Error().Err(err).Msg("Failed to add conflicts")
			return nil, fmt.Errorf("failed to add conflicts: %w", err)
		}
	}

	if len(operations) > 0 {
		if _, err := s.dataClient.BatchProcess(ctx, userID, operations); err != nil {
			if len(conflicts) > 0 {
				if rollbackErr := s.RollbackInsertedConflicts(ctx, conflicts); rollbackErr != nil {
					log.Error().Err(rollbackErr).Msg("Failed to rollback conflicts")
					return nil, fmt.Errorf("failed to process batch operations and rollback conflicts: %w (rollback error: %v)", err, rollbackErr)
				}
			}
			return nil, fmt.Errorf("failed to process batch operations: %w", err)
		}
	}

	log.Info().Msg("SyncData completed successfully")
	return conflicts, nil
}

// ResolveConflict handles conflict resolution using specified strategy
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
		resolvedData, err = MergeData(conflict.ClientData, conflict.ServerData, models.Text)
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

// ListConflicts returns all unresolved conflicts for user
func (s *Service) ListConflicts(ctx context.Context, userID string) ([]models.Conflict, error) {
	conflicts, err := s.conflictRepo.GetUnresolvedConflicts(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unresolved conflicts: %w", err)
	}
	return conflicts, nil
}

// RollbackInsertedConflicts removes previously inserted conflicts (used for error recovery)
func (s *Service) RollbackInsertedConflicts(ctx context.Context, conflicts []models.Conflict) error {
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

// MergeData implements different merging strategies based on data type
func MergeData(clientData, serverData []byte, dataType models.DataType) ([]byte, error) {
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
