package topic

import (
	"context"
	"gallery-service/internal/domain/models"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/zap"
)

type GetTopicByIDQueryHandler interface {
	Handle(ctx context.Context, command *GetTopicByIDQuery) (*models.Topic, error)
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
