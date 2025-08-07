package topic

import (
	"context"
	"gallery-service/internal/domain/models"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/zap"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	cluster, err := u.topicRepo.GetByID(ctx, command.ID)
	if err != nil {
		return err
	}

	folderID, err := primitive.ObjectIDFromHex(command.FolderID)
	if err != nil {
		return errors.New("invalid folder id")
	}

	// check if cluster exist
	exist, err := u.folderRepo.Exists(ctx, bson.M{"_id": folderID})
	if err != nil {
		return err
	}

	if !exist {
		return errors.New("folder not found")
	}

	t := models.Topic{
		ID:             cluster.ID,
		TopicName:      command.TopicName,
		Title:          command.Title,
		Note:           command.Note,
		Image:          command.Image,
		LanguageConfig: command.LanguageConfig,
		FolderID:       folderID,
		CreatedAt:      cluster.CreatedAt,
		UpdatedAt:      time.Now(),
	}

	// Save to database
	return u.topicRepo.Update(ctx, &t)
}
