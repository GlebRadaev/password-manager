// Package api implements gRPC server interface for sync service.
package api

//go:generate mockgen -destination=api_mock.go -source=api.go -package=api
import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/GlebRadaev/password-manager/internal/sync/models"
	"github.com/GlebRadaev/password-manager/internal/sync/service"
	"github.com/GlebRadaev/password-manager/pkg/sync"
)

// Service defines the business logic interface for sync operations.
type Service interface {
	SyncData(ctx context.Context, userID string, clientData []models.ClientData) ([]models.Conflict, error)
	ResolveConflict(ctx context.Context, conflictID string, strategy models.ResolutionStrategy) error
	ListConflicts(ctx context.Context, userID string) ([]models.Conflict, error)
}

// API implements sync.SyncServiceServer gRPC interface.
type API struct {
	sync.UnimplementedSyncServiceServer
	srv Service
}

// New creates new Api instance with given service.
func New(srv Service) *API {
	return &API{srv: srv}
}

// SyncData handles data synchronization between client and server.
// Returns list of conflicts if any occur during sync.
func (s *API) SyncData(ctx context.Context, req *sync.SyncDataRequest) (*sync.SyncDataResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	var clientData []models.ClientData
	for _, entry := range req.ClientData {
		var metadata []models.Metadata
		for _, m := range entry.Metadata {
			metadata = append(metadata, models.Metadata{
				Key:   m.Key,
				Value: m.Value,
			})
		}

		clientData = append(clientData, models.ClientData{
			DataID:    entry.DataId,
			Type:      toModelDataType(entry.Type),
			Data:      entry.Data,
			Metadata:  metadata,
			Operation: toModelOperation(entry.Operation),
			UpdatedAt: time.Unix(entry.UpdatedAt, 0),
		})
	}

	conflicts, err := s.srv.SyncData(ctx, req.UserId, clientData)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	var protoConflicts []*sync.Conflict
	for _, conflict := range conflicts {
		protoConflicts = append(protoConflicts, &sync.Conflict{
			ConflictId: conflict.ID,
			DataId:     conflict.DataID,
			ClientData: conflict.ClientData,
			ServerData: conflict.ServerData,
			Resolved:   conflict.Resolved,
			CreatedAt:  conflict.CreatedAt.Unix(),
			UpdatedAt:  conflict.UpdatedAt.Unix(),
		})
	}

	return &sync.SyncDataResponse{Conflicts: protoConflicts}, nil
}

// ResolveConflict handles conflict resolution with specified strategy.
func (s *API) ResolveConflict(ctx context.Context, req *sync.ResolveConflictRequest) (*sync.ResolveConflictResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	strategy := toModelResolutionStrategy(req.Strategy)
	if err := s.srv.ResolveConflict(ctx, req.ConflictId, strategy); err != nil {
		if errors.Is(err, service.ErrConflictNotFound) {
			return nil, status.Errorf(codes.NotFound, "%v", err)
		}
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &sync.ResolveConflictResponse{Message: "Conflict resolved successfully"}, nil
}

// ListConflicts returns all unresolved conflicts for given user.
func (s *API) ListConflicts(ctx context.Context, req *sync.ListConflictsRequest) (*sync.ListConflictsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	conflicts, err := s.srv.ListConflicts(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &sync.ListConflictsResponse{
		Conflicts: convertConflictsToProto(conflicts),
	}, nil
}

// convertConflictsToProto converts domain conflicts to protobuf format.
func convertConflictsToProto(conflicts []models.Conflict) []*sync.Conflict {
	var protoConflicts []*sync.Conflict
	for _, c := range conflicts {
		protoConflicts = append(protoConflicts, &sync.Conflict{
			ConflictId: c.ID,
			DataId:     c.DataID,
			ClientData: c.ClientData,
			ServerData: c.ServerData,
			Resolved:   c.Resolved,
			CreatedAt:  c.CreatedAt.Unix(),
			UpdatedAt:  c.UpdatedAt.Unix(),
		})
	}
	return protoConflicts
}

// toModelDataType converts protobuf DataType to domain model.
func toModelDataType(dt sync.DataType) models.DataType {
	switch dt {
	case sync.DataType_LOGIN_PASSWORD:
		return models.LoginPassword
	case sync.DataType_TEXT:
		return models.Text
	case sync.DataType_BINARY:
		return models.Binary
	case sync.DataType_CARD:
		return models.Card
	default:
		return models.LoginPassword
	}
}

// toProtoDataType converts domain DataType to protobuf (unused in current implementation).
func toProtoDataType(dt models.DataType) sync.DataType {
	switch dt {
	case models.LoginPassword:
		return sync.DataType_LOGIN_PASSWORD
	case models.Text:
		return sync.DataType_TEXT
	case models.Binary:
		return sync.DataType_BINARY
	case models.Card:
		return sync.DataType_CARD
	default:
		return sync.DataType_LOGIN_PASSWORD
	}
}

// toModelResolutionStrategy converts protobuf ResolutionStrategy to domain model.
func toModelResolutionStrategy(strategy sync.ResolutionStrategy) models.ResolutionStrategy {
	switch strategy {
	case sync.ResolutionStrategy_USE_CLIENT_VERSION:
		return models.UseClientVersion
	case sync.ResolutionStrategy_USE_SERVER_VERSION:
		return models.UseServerVersion
	case sync.ResolutionStrategy_MERGE_VERSIONS:
		return models.MergeVersions
	default:
		return models.UseClientVersion
	}
}

// toModelOperation converts protobuf Operation to domain model.
func toModelOperation(op sync.Operation) models.Operation {
	switch op {
	case sync.Operation_ADD:
		return models.Add
	case sync.Operation_UPDATE:
		return models.Update
	case sync.Operation_DELETE:
		return models.Delete
	default:
		return models.Add
	}
}
