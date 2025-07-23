package cluster

import (
	"context"
	"gallery-service/internal/application/dto/responses/cluster"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/utils"
	"gallery-service/pkg/zap"
)

type GetAllClusterQueryHandler interface {
	Handle(ctx context.Context, pq *utils.Pagination) (*cluster.GetAllClusterResponseDto, error)
}

type getClusterHandler struct {
	log         zap.Logger
	clusterRepo repository.ClusterRepository
}

func NewGetAllClusterHandler(log zap.Logger, clusterRepo repository.ClusterRepository) *getClusterHandler {
	return &getClusterHandler{log: log, clusterRepo: clusterRepo}
}

func (q *getClusterHandler) Handle(ctx context.Context, pq *utils.Pagination) (*cluster.GetAllClusterResponseDto, error) {
	return q.clusterRepo.GetAll(ctx, pq)
}
