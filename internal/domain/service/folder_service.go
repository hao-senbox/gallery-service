package service

import (
	folderCommands "gallery-service/internal/application/commands/v1/folder"
	"gallery-service/internal/application/queries/folder"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/kafka"
	"gallery-service/pkg/zap"
)

type FolderService struct {
	Commands *folderCommands.Commands
	Queries  *folder.Queries
}

var (
	folderService *FolderService
)

func NewFolderService(
	cfg kafka.Config,
	log zap.Logger,
	folderRepo repository.FolderRepository,
) *FolderService {
	if folderService != nil {
		return folderService
	}

	createFolderHandler := folderCommands.NewCreateFolderHandler(cfg, log, folderRepo)
	updateFolderHandler := folderCommands.NewUpdateFolderHandler(log, folderRepo)
	deleteFolderHandler := folderCommands.NewDeleteFolderHandler(log, folderRepo)

	getAllFolderHandler := folder.NewGetAllFolderHandler(log, folderRepo)
	getFolderByIDHandler := folder.NewGetFolderByIDHandler(log, folderRepo)
	searchFoldersHandler := folder.NewSearchFoldersHandler(log, folderRepo)

	commands := folderCommands.NewFolderCommands(
		createFolderHandler,
		updateFolderHandler,
		deleteFolderHandler,
	)
	queries := folder.NewFolderQueries(
		getAllFolderHandler,
		getFolderByIDHandler,
		searchFoldersHandler,
	)

	folderService = &FolderService{Commands: commands, Queries: queries}

	return folderService
}
