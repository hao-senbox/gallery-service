package folder

import (
	"context"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/zap"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteFolderCommandHandler interface {
	Handle(ctx context.Context, command *DeleteFolderCommand) error
}

type deleteFolderHandler struct {
	log        zap.Logger
	folderRepo repository.FolderRepository
}

func NewDeleteFolderHandler(
	log zap.Logger,
	folderRepo repository.FolderRepository,
) *deleteFolderHandler {
	return &deleteFolderHandler{
		log:        log,
		folderRepo: folderRepo,
	}
}

func (u *deleteFolderHandler) Handle(ctx context.Context, command *DeleteFolderCommand) error {
	id, _ := primitive.ObjectIDFromHex(command.ID)
	exist, err := u.folderRepo.Exists(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if !exist {
		return errors.New("folder not found")
	}

	if ok, err := u.folderRepo.Delete(ctx, command.ID); !ok || err != nil {
		u.log.Errorf("(DeleteFolderCommandHandler.Handle) err: {%v}", err)
		return errors.New("failed to delete folder")
	}

	return nil
}
