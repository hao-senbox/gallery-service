package cluster

type Commands struct {
	CreateCluster CreateClusterCommandHandler
	UpdateCluster UpdateClusterCommandHandler
	DeleteCluster DeleteClusterCommandHandler
}

func NewClusterCommands(
	createCluster CreateClusterCommandHandler,
	updateCluster UpdateClusterCommandHandler,
	deleteCluster DeleteClusterCommandHandler,
) *Commands {
	return &Commands{
		CreateCluster: createCluster,
		UpdateCluster: updateCluster,
		DeleteCluster: deleteCluster,
	}
}
