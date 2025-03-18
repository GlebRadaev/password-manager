package api

//go:generate mockgen -destination=api_mock.go -source=api.go -package=api
import (
	"context"
	"time"

	"github.com/GlebRadaev/password-manager/internal/sync/models"
	"github.com/GlebRadaev/password-manager/pkg/sync"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	GetChanges(ctx context.Context, userID string, lastSyncTime int64) ([]models.DataChange, error)
	PushChanges(ctx context.Context, userID string, changes []models.DataChange) error
}
type Api struct {
	sync.UnimplementedSyncServiceServer
	srv Service
}

func New(srv Service) *Api {
	return &Api{srv: srv}
}

func (a *Api) GetChanges(ctx context.Context, req *sync.GetChangesRequest) (*sync.GetChangesResponse, error) {
	changes, err := a.srv.GetChanges(ctx, req.UserId, req.LastSyncTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get changes: %v", err)
	}

	var protoChanges []*sync.DataChange
	for _, change := range changes {
		protoChanges = append(protoChanges, &sync.DataChange{
			DataId:    change.DataID,
			Type:      change.Type,
			Data:      change.Data,
			Metadata:  change.Metadata,
			Timestamp: change.Timestamp.Unix(),
		})
	}

	return &sync.GetChangesResponse{Changes: protoChanges}, nil
}

func (a *Api) PushChanges(ctx context.Context, req *sync.PushChangesRequest) (*sync.PushChangesResponse, error) {
	var changes []models.DataChange
	for _, protoChange := range req.Changes {
		changes = append(changes, models.DataChange{
			UserID:    req.UserId,
			DataID:    protoChange.DataId,
			Type:      protoChange.Type,
			Data:      protoChange.Data,
			Metadata:  protoChange.Metadata,
			Timestamp: time.Unix(protoChange.Timestamp, 0),
		})
	}

	if err := a.srv.PushChanges(ctx, req.UserId, changes); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to push changes: %v", err)
	}

	return &sync.PushChangesResponse{Success: true, Message: "Changes pushed successfully"}, nil
}
