package topic

import (
	"context"
	"gallery-service/internal/domain/models"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/zap"
	"time"
)

type UpdateTopicCommandHandler interface {
	Handle(ctx context.Context, command *UpdateTopicCommand) error
}

type updateTopicHandler struct {
	log        zap.Logger
	topicRepo  repository.TopicRepository
	folderRepo repository.FolderRepository
}

func NewUpdateTopicHandler(
	log zap.Logger,
	topicRepo repository.TopicRepository,
	folderRepo repository.FolderRepository,
) *updateTopicHandler {
	return &updateTopicHandler{
		log:        log,
		topicRepo:  topicRepo,
		folderRepo: folderRepo,
	}
}

func (u *updateTopicHandler) Handle(ctx context.Context, command *UpdateTopicCommand) error {
	topic, err := u.topicRepo.GetByID(ctx, command.ID)
	if err != nil {
		return err
	}

	t := models.Topic{
		ID:             topic.ID,
		TopicName:      command.TopicName,
		Title:          command.Title,
		Note:           command.Note,
		Images:         command.Image,
		LanguageConfig: command.LanguageConfig,
		CreatedAt:      topic.CreatedAt,
		UpdatedAt:      time.Now(),
	}

	// Save to database
	return u.topicRepo.Update(ctx, &t)
}
