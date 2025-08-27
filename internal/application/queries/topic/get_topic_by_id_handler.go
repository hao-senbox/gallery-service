package topic

import (
	"context"
	"gallery-service/internal/application/dto/responses/topic"
	"gallery-service/internal/application/mappers"
	"gallery-service/internal/domain/models"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/zap"
)

type GetTopicByIDQueryHandler interface {
	Handle(ctx context.Context, command *GetTopicByIDQuery) (*models.Topic, error)
	Handle4Gateway(ctx context.Context, command *GetTopicByIDQuery) (*topic.Topic4GatwayResponseDto, error)
}

type getTopicByIDHandler struct {
	log      zap.Logger
	taskRepo repository.TopicRepository
}

func NewGetTopicByIDHandler(log zap.Logger, taskRepo repository.TopicRepository) *getTopicByIDHandler {
	return &getTopicByIDHandler{log: log, taskRepo: taskRepo}
}

func (q *getTopicByIDHandler) Handle(ctx context.Context, query *GetTopicByIDQuery) (*models.Topic, error) {
	return q.taskRepo.GetByID(ctx, query.ID)
}

func (q *getTopicByIDHandler) Handle4Gateway(ctx context.Context, query *GetTopicByIDQuery) (*topic.Topic4GatwayResponseDto, error) {

	topic, err := q.taskRepo.GetByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	return mappers.GetTopic4GatewayFromModel(topic), nil
}
