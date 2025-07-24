package cluster

import (
	"gallery-service/pkg/utils"
)

type Queries struct {
	GetAllCluster       GetAllClusterQueryHandler
	GetClusterByID      GetClusterByIDQueryHandler
	SearchClusters      SearchClustersQueryHandler
	GetAllClusterFolder GetAllClusterFolderQueryHandler
}

func NewClusterQueries(
	getAllCluster GetAllClusterQueryHandler,
	getClusterByID GetClusterByIDQueryHandler,
	searchClusters SearchClustersQueryHandler,
	getClusterFolder GetAllClusterFolderQueryHandler,
) *Queries {
	return &Queries{
		GetAllCluster:       getAllCluster,
		GetClusterByID:      getClusterByID,
		SearchClusters:      searchClusters,
		GetAllClusterFolder: getClusterFolder,
	}
}

type GetClusterByIDQuery struct {
	ID string `json:"id" validate:"required"`
}

func NewGetClusterByIDQuery(ID string) *GetClusterByIDQuery {
	return &GetClusterByIDQuery{ID: ID}
}

type GetFolderID struct {
	ID string `json:"folder_id" validate:"required"`
}

func NewGetFolderID(ID string) *GetFolderID {
	return &GetFolderID{ID: ID}
}

type SearchClustersQuery struct {
	Keyword string
	Pq      *utils.Pagination
}

func NewSearchClustersQuery(
	keyword string,
	pq *utils.Pagination,
) *SearchClustersQuery {
	return &SearchClustersQuery{
		Keyword: keyword,
		Pq:      pq,
	}
}
