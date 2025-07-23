package folder

import (
	"context"
	"gallery-service/internal/application/dto/responses/folder"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/zap"
)

type SearchFoldersQueryHandler interface {
	Handle(ctx context.Context, command *SearchFoldersQuery) (*folder.GetAllFolderResponseDto, error)
}

type searchFoldersHandler struct {
	log            zap.Logger
	taskRepository repository.FolderRepository
}

func NewSearchFoldersHandler(log zap.Logger, taskRepository repository.FolderRepository) *searchFoldersHandler {
	return &searchFoldersHandler{log: log, taskRepository: taskRepository}
}

func (s *searchFoldersHandler) Handle(ctx context.Context, command *SearchFoldersQuery) (*folder.GetAllFolderResponseDto, error) {
	query := make(map[string]interface{})
	query["keyword"] = command.Keyword

	return s.taskRepository.Search(ctx, query, command.Pq)
}
