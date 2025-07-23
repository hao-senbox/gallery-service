package folder

import (
	"context"
	"gallery-service/internal/domain/models"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/kafka"
	"gallery-service/pkg/zap"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateFolderCommandHandler interface {
	Handle(ctx context.Context, command *CreateFolderCommand) (*string, error)
}

type createFolderHandler struct {
	cfg        kafka.Config
	log        zap.Logger
	folderRepo repository.FolderRepository
}

func NewCreateFolderHandler(
	cfg kafka.Config,
	log zap.Logger,
	folderRepo repository.FolderRepository,
) *createFolderHandler {
	return &createFolderHandler{
		cfg:        cfg,
		log:        log,
		folderRepo: folderRepo,
	}
}

func (c *createFolderHandler) Handle(ctx context.Context, command *CreateFolderCommand) (*string, error) {
	id := primitive.NewObjectID()
	var parentID *primitive.ObjectID
	if command.ParentID != nil {
		pID, err := primitive.ObjectIDFromHex(*command.ParentID)
		if err != nil {
			return nil, err
		}

		parentID = &pID
	}
	folder := models.Folder{
		ID:                 id,
		FolderName:         command.FolderName,
		FolderThumbnailKey: command.FolderThumbnailKey,
		FolderThumbnailURL: command.FolderThumbnailURL,
		ParentID:           parentID,
	}

	// Save to database
	folderID, err := c.folderRepo.Insert(ctx, &folder)
	if err != nil {
		return nil, err
	}

	return &folderID, nil
}
