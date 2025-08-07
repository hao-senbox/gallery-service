package topic

import (
	"context"
	"gallery-service/internal/application/dto/responses/topic"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/utils"
	"gallery-service/pkg/zap"
)

type GetAllTopicFolderQueryHandler interface {
	Handle(ctx context.Context, command *GetFolderID, pq *utils.Pagination) (*topic.GetAllTopicResponseDto, error)
}

type getTopicFolderHandler struct {
	log       zap.Logger
	topicRepo repository.TopicRepository
}

func NewGetAllClusterFolderHandler(log zap.Logger, topicRepo repository.TopicRepository) *getTopicFolderHandler {
	return &getTopicFolderHandler{log: log, topicRepo: topicRepo}
}

// func (q *getTopicFolderHandler) Handle(ctx context.Context, command *GetFolderID, pq *utils.Pagination) (*topic.GetAllTopicResponseDto, error) {
// 	return q.topicRepo.GetAllByFolderID(ctx, command.ID, pq)
// }
