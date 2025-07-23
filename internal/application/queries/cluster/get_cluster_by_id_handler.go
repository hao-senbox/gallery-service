package cluster

import (
	"context"
	"gallery-service/internal/domain/models"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/zap"
)

type GetClusterByIDQueryHandler interface {
	Handle(ctx context.Context, command *GetClusterByIDQuery) (*models.Cluster, error)
}

type getClusterByIDHandler struct {
	log      zap.Logger
	taskRepo repository.ClusterRepository
}

func NewGetClusterByIDHandler(log zap.Logger, taskRepo repository.ClusterRepository) *getClusterByIDHandler {
	return &getClusterByIDHandler{log: log, taskRepo: taskRepo}
}

func (q *getClusterByIDHandler) Handle(ctx context.Context, query *GetClusterByIDQuery) (*models.Cluster, error) {
	return q.taskRepo.GetByID(ctx, query.ID)
}
