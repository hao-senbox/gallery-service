package topic

import (
	"context"
	"gallery-service/internal/application/dto/responses/topic"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/utils"
	"gallery-service/pkg/zap"
)

type GetAllTopicQueryHandler interface {
	Handle(ctx context.Context, pq *utils.Pagination) (*topic.GetAllTopicResponseDto, error)
	Handle4App(ctx context.Context) (*topic.GetAllTopicForAppResponseDto, error)
}

type getTopicHandler struct {
	log       zap.Logger
	topicRepo repository.TopicRepository
}

func NewGetAllTopicHandler(log zap.Logger, topicRepo repository.TopicRepository) *getTopicHandler {
	return &getTopicHandler{log: log, topicRepo: topicRepo}
}

func (q *getTopicHandler) Handle(ctx context.Context, pq *utils.Pagination) (*topic.GetAllTopicResponseDto, error) {
	return q.topicRepo.GetAll(ctx, pq)
}

func (q *getTopicHandler) Handle4App(ctx context.Context) (*topic.GetAllTopicForAppResponseDto, error) {
	return q.topicRepo.GetAll4App(ctx)
}
