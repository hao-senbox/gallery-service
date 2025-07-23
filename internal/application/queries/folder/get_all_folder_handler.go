package folder

import (
	"context"
	"gallery-service/internal/application/dto/responses/folder"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/utils"
	"gallery-service/pkg/zap"
)

type GetAllFolderQueryHandler interface {
	Handle(ctx context.Context, pq *utils.Pagination) (*folder.GetAllFolderResponseDto, error)
}

type getFolderHandler struct {
	log         zap.Logger
	clusterRepo repository.FolderRepository
}

func NewGetAllFolderHandler(log zap.Logger, clusterRepo repository.FolderRepository) *getFolderHandler {
	return &getFolderHandler{log: log, clusterRepo: clusterRepo}
}

func (q *getFolderHandler) Handle(ctx context.Context, pq *utils.Pagination) (*folder.GetAllFolderResponseDto, error) {
	return q.clusterRepo.GetAll(ctx, pq)
}
