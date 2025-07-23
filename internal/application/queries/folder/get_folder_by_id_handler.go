package folder

import (
	"context"
	"gallery-service/internal/domain/models"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/zap"
)

type GetFolderByIDQueryHandler interface {
	Handle(ctx context.Context, command *GetFolderByIDQuery) (*models.Folder, error)
}

type getFolderByIDHandler struct {
	log      zap.Logger
	taskRepo repository.FolderRepository
}

func NewGetFolderByIDHandler(log zap.Logger, taskRepo repository.FolderRepository) *getFolderByIDHandler {
	return &getFolderByIDHandler{log: log, taskRepo: taskRepo}
}

func (q *getFolderByIDHandler) Handle(ctx context.Context, query *GetFolderByIDQuery) (*models.Folder, error) {
	return q.taskRepo.GetByID(ctx, query.ID)
}
