package folder

import (
	"gallery-service/pkg/utils"
)

type Queries struct {
	GetAllFolder  GetAllFolderQueryHandler
	GetFolderByID GetFolderByIDQueryHandler
	SearchFolders SearchFoldersQueryHandler
}

func NewFolderQueries(
	getAllFolder GetAllFolderQueryHandler,
	getFolderByID GetFolderByIDQueryHandler,
	searchFolders SearchFoldersQueryHandler,
) *Queries {
	return &Queries{
		GetAllFolder:  getAllFolder,
		GetFolderByID: getFolderByID,
		SearchFolders: searchFolders,
	}
}

type GetFolderByIDQuery struct {
	ID string `json:"id" validate:"required"`
}

func NewGetFolderByIDQuery(ID string) *GetFolderByIDQuery {
	return &GetFolderByIDQuery{ID: ID}
}

type SearchFoldersQuery struct {
	Keyword string
	Pq      *utils.Pagination
}

func NewSearchFoldersQuery(
	keyword string,
	pq *utils.Pagination,
) *SearchFoldersQuery {
	return &SearchFoldersQuery{
		Keyword: keyword,
		Pq:      pq,
	}
}
