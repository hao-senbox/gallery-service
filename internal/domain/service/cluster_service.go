package service

import (
	clusterCommands "gallery-service/internal/application/commands/v1/cluster"
	"gallery-service/internal/application/queries/cluster"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/kafka"
	"gallery-service/pkg/zap"
)

type ClusterService struct {
	Commands *clusterCommands.Commands
	Queries  *cluster.Queries
}

var (
	clusterService *ClusterService
)

func NewClusterService(
	cfg kafka.Config,
	log zap.Logger,
	clusterRepo repository.ClusterRepository,
	folderRepo repository.FolderRepository,
) *ClusterService {
	if clusterService != nil {
		return clusterService
	}

	createClusterHandler := clusterCommands.NewCreateClusterHandler(cfg, log, clusterRepo, folderRepo)
	updateClusterHandler := clusterCommands.NewUpdateClusterHandler(log, clusterRepo, folderRepo)
	deleteClusterHandler := clusterCommands.NewDeleteClusterHandler(log, clusterRepo)

	getAllClusterHandler := cluster.NewGetAllClusterHandler(log, clusterRepo)
	getClusterFolder := cluster.NewGetAllClusterFolderHandler(log, clusterRepo)
	getClusterByIDHandler := cluster.NewGetClusterByIDHandler(log, clusterRepo)
	searchClustersHandler := cluster.NewSearchClustersHandler(log, clusterRepo)

	commands := clusterCommands.NewClusterCommands(
		createClusterHandler,
		updateClusterHandler,
		deleteClusterHandler,
	)
	queries := cluster.NewClusterQueries(
		getAllClusterHandler,
		getClusterByIDHandler,
		searchClustersHandler,
		getClusterFolder,
	)

	clusterService = &ClusterService{Commands: commands, Queries: queries}

	return clusterService
}
