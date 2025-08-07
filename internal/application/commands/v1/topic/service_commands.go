package topic

type Commands struct {
	CreateTopic CreateTopicCommandHandler
	UpdateTopic UpdateTopicCommandHandler
	DeleteTopic DeleteTopicCommandHandler
}

func NewTopicCommands(
	createTopic CreateTopicCommandHandler,
	updateTopic UpdateTopicCommandHandler,
	deleteTopic DeleteTopicCommandHandler,
) *Commands {
	return &Commands{
		CreateTopic: createTopic,
		UpdateTopic: updateTopic,
		DeleteTopic: deleteTopic,
	}
}
