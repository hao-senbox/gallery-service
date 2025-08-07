package topic

import (
	"context"
	"gallery-service/internal/application/dto/responses/topic"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/zap"
)

type SearchTopicsQueryHandler interface {
	Handle(ctx context.Context, command *SearchTopicsQuery) (*topic.GetAllTopicResponseDto, error)
}

type searchTopicsHandler struct {
	log            zap.Logger
	taskRepository repository.TopicRepository
}

func NewSearchTopicsHandler(log zap.Logger, taskRepository repository.TopicRepository) *searchTopicsHandler {
	return &searchTopicsHandler{log: log, taskRepository: taskRepository}
}

func (s *searchTopicsHandler) Handle(ctx context.Context, command *SearchTopicsQuery) (*topic.GetAllTopicResponseDto, error) {
	query := make(map[string]interface{})
	query["keyword"] = command.Keyword

	return s.taskRepository.Search(ctx, query, command.Pq)
}
