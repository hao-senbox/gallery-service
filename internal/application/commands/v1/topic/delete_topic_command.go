package topic

type DeleteTopicCommand struct {
	ID string `json:"id" validate:"required"`
}

func NewDeleteTopicCommand(id string) *DeleteTopicCommand {
	return &DeleteTopicCommand{
		ID: id,
	}
}
