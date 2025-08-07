package topic

import (
	"gallery-service/pkg/utils"
)

type Queries struct {
	GetAllTopic       GetAllTopicQueryHandler
	GetTopicByID      GetTopicByIDQueryHandler
	SearchTopics      SearchTopicsQueryHandler
	GetAllTopicFolder GetAllTopicFolderQueryHandler
}

func NewTopicQueries(
	getAllTopic GetAllTopicQueryHandler,
	getTopicByID GetTopicByIDQueryHandler,
	searchTopics SearchTopicsQueryHandler,
	//getTopicFolder GetAllTopicFolderQueryHandler,
) *Queries {
	return &Queries{
		GetAllTopic:  getAllTopic,
		GetTopicByID: getTopicByID,
		SearchTopics: searchTopics,
		//GetAllTopicFolder: getTopicFolder,
	}
}

type GetTopicByIDQuery struct {
	ID string `json:"id" validate:"required"`
}

func NewGetTopicByIDQuery(ID string) *GetTopicByIDQuery {
	return &GetTopicByIDQuery{ID: ID}
}

type GetFolderID struct {
	ID string `json:"folder_id" validate:"required"`
}

func NewGetFolderID(ID string) *GetFolderID {
	return &GetFolderID{ID: ID}
}

type SearchTopicsQuery struct {
	Keyword string
	Pq      *utils.Pagination
}

func NewSearchTopicsQuery(
	keyword string,
	pq *utils.Pagination,
) *SearchTopicsQuery {
	return &SearchTopicsQuery{
		Keyword: keyword,
		Pq:      pq,
	}
}
