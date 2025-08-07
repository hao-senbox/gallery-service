package topic

import (
	"context"
	"gallery-service/internal/domain/models"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/kafka"
	"gallery-service/pkg/zap"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateTopicCommandHandler interface {
	Handle(ctx context.Context, command *CreateTopicCommand) (*string, error)
}

type createTopicHandler struct {
	cfg        kafka.Config
	log        zap.Logger
	topicRepo  repository.TopicRepository
	folderRepo repository.FolderRepository
}

func NewCreateTopicHandler(
	cfg kafka.Config,
	log zap.Logger,
	topicRepo repository.TopicRepository,
	folderRepo repository.FolderRepository,
) *createTopicHandler {
	return &createTopicHandler{
		cfg:        cfg,
		log:        log,
		topicRepo:  topicRepo,
		folderRepo: folderRepo,
	}
}

func (c *createTopicHandler) Handle(ctx context.Context, command *CreateTopicCommand) (*string, error) {
	id := primitive.NewObjectID()

	topic := models.Topic{
		ID:             id,
		TopicName:      command.TopicName,
		Title:          command.Title,
		Note:           command.Note,
		Images:         command.Image,
		LanguageConfig: []models.TopicLanguageConfig{command.LanguageConfig},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Save to database
	topicID, err := c.topicRepo.Insert(ctx, &topic)
	if err != nil {
		return nil, err
	}

	return &topicID, nil
}
