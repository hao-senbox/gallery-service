package service

import (
	topicCommands "gallery-service/internal/application/commands/v1/topic"
	"gallery-service/internal/application/queries/topic"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/kafka"
	"gallery-service/pkg/zap"
)

type TopicService struct {
	Commands *topicCommands.Commands
	Queries  *topic.Queries
}

var (
	topicService *TopicService
)

func NewTopicService(
	cfg kafka.Config,
	log zap.Logger,
	topicRepo repository.TopicRepository,
	folderRepo repository.FolderRepository,
) *TopicService {
	if topicService != nil {
		return topicService
	}

	createTopicHandler := topicCommands.NewCreateTopicHandler(cfg, log, topicRepo)
	updateTopicHandler := topicCommands.NewUpdateTopicHandler(log, topicRepo)
	deleteTopicHandler := topicCommands.NewDeleteTopicHandler(log, topicRepo)

	getAllTopicHandler := topic.NewGetAllTopicHandler(log, topicRepo)
	//getTopicFolder := topic.NewGetAllTopicFolderHandler(log, topicRepo)
	getTopicByIDHandler := topic.NewGetTopicByIDHandler(log, topicRepo)
	searchTopicsHandler := topic.NewSearchTopicsHandler(log, topicRepo)

	commands := topicCommands.NewTopicCommands(
		createTopicHandler,
		updateTopicHandler,
		deleteTopicHandler,
	)
	queries := topic.NewTopicQueries(
		getAllTopicHandler,
		getTopicByIDHandler,
		searchTopicsHandler,
		//getTopicFolder,
	)

	topicService = &TopicService{Commands: commands, Queries: queries}

	return topicService
}
