package folder

import (
	"context"
	"gallery-service/internal/domain/models"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/zap"
)

type UpdateFolderCommandHandler interface {
	Handle(ctx context.Context, command *UpdateFolderCommand) error
}

type updateFolderHandler struct {
	log        zap.Logger
	folderRepo repository.FolderRepository
}

func NewUpdateFolderHandler(
	log zap.Logger,
	folderRepo repository.FolderRepository,
) *updateFolderHandler {
	return &updateFolderHandler{
		log:        log,
		folderRepo: folderRepo,
	}
}

func (u *updateFolderHandler) Handle(ctx context.Context, command *UpdateFolderCommand) error {
	folder, err := u.folderRepo.GetByID(ctx, command.ID)
	if err != nil {
		return err
	}

	t := models.Folder{
		ID:                 folder.ID,
		FolderName:         command.FolderName,
		FolderThumbnailKey: command.FolderThumbnailKey,
		FolderThumbnailURL: command.FolderThumbnailURL,
		ParentID:           folder.ParentID,
	}

	// Save to database
	return u.folderRepo.Update(ctx, &t)
}
