package cluster

import (
	"context"
	"gallery-service/internal/application/dto/responses/cluster"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/utils"
	"gallery-service/pkg/zap"
)

type GetAllClusterFolderQueryHandler interface {
	Handle(ctx context.Context, command *GetFolderID, pq *utils.Pagination) (*cluster.GetAllClusterResponseDto, error)
}

type getClusterFolderHandler struct {
	log         zap.Logger
	clusterRepo repository.ClusterRepository
}

func NewGetAllClusterFolderHandler(log zap.Logger, clusterRepo repository.ClusterRepository) *getClusterFolderHandler {
	return &getClusterFolderHandler{log: log, clusterRepo: clusterRepo}
}

func (q *getClusterFolderHandler) Handle(ctx context.Context, command *GetFolderID, pq *utils.Pagination) (*cluster.GetAllClusterResponseDto, error) {
	return q.clusterRepo.GetAllByFolderID(ctx, command.ID, pq)
}
