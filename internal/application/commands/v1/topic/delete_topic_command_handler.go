package topic

import (
	"context"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/zap"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteTopicCommandHandler interface {
	Handle(ctx context.Context, command *DeleteTopicCommand) error
}

type deleteTopicHandler struct {
	log       zap.Logger
	topicRepo repository.TopicRepository
}

func NewDeleteTopicHandler(
	log zap.Logger,
	topicRepo repository.TopicRepository,
) *deleteTopicHandler {
	return &deleteTopicHandler{
		log:       log,
		topicRepo: topicRepo,
	}
}

func (u *deleteTopicHandler) Handle(ctx context.Context, command *DeleteTopicCommand) error {
	id, _ := primitive.ObjectIDFromHex(command.ID)
	exist, err := u.topicRepo.Exists(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if !exist {
		return errors.New("topic not found")
	}

	if ok, err := u.topicRepo.Delete(ctx, command.ID); !ok || err != nil {
		u.log.Errorf("(DeleteTopicCommandHandler.Handle) err: {%v}", err)
		return errors.New("failed to delete topic")
	}

	return nil
}
