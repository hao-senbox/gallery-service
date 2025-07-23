package cluster

import (
	"context"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/zap"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteClusterCommandHandler interface {
	Handle(ctx context.Context, command *DeleteClusterCommand) error
}

type deleteClusterHandler struct {
	log         zap.Logger
	clusterRepo repository.ClusterRepository
}

func NewDeleteClusterHandler(
	log zap.Logger,
	clusterRepo repository.ClusterRepository,
) *deleteClusterHandler {
	return &deleteClusterHandler{
		log:         log,
		clusterRepo: clusterRepo,
	}
}

func (u *deleteClusterHandler) Handle(ctx context.Context, command *DeleteClusterCommand) error {
	id, _ := primitive.ObjectIDFromHex(command.ID)
	exist, err := u.clusterRepo.Exists(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if !exist {
		return errors.New("cluster not found")
	}

	if ok, err := u.clusterRepo.Delete(ctx, command.ID); !ok || err != nil {
		u.log.Errorf("(DeleteClusterCommandHandler.Handle) err: {%v}", err)
		return errors.New("failed to delete cluster")
	}

	return nil
}
