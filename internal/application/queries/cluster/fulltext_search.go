package cluster

import (
	"context"
	"gallery-service/internal/application/dto/responses/cluster"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/zap"
)

type SearchClustersQueryHandler interface {
	Handle(ctx context.Context, command *SearchClustersQuery) (*cluster.GetAllClusterResponseDto, error)
}

type searchClustersHandler struct {
	log            zap.Logger
	taskRepository repository.ClusterRepository
}

func NewSearchClustersHandler(log zap.Logger, taskRepository repository.ClusterRepository) *searchClustersHandler {
	return &searchClustersHandler{log: log, taskRepository: taskRepository}
}

func (s *searchClustersHandler) Handle(ctx context.Context, command *SearchClustersQuery) (*cluster.GetAllClusterResponseDto, error) {
	query := make(map[string]interface{})
	query["keyword"] = command.Keyword

	return s.taskRepository.Search(ctx, query, command.Pq)
}
