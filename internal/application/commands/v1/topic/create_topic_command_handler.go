package topic

import (
	"context"
	"gallery-service/internal/domain/models"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/kafka"
	"gallery-service/pkg/zap"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
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
	folderID, err := primitive.ObjectIDFromHex(command.FolderID)
	if err != nil {
		return nil, errors.New("invalid folder id")
	}

	// check if folder exist
	exist, err := c.folderRepo.Exists(ctx, bson.M{"_id": folderID})
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("folder not found")
	}

	topic := models.Topic{
		ID:             id,
		TopicName:      command.TopicName,
		Title:          command.Title,
		Note:           command.Note,
		Image:          command.Image,
		LanguageConfig: command.LanguageConfig,
		FolderID:       folderID,
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
