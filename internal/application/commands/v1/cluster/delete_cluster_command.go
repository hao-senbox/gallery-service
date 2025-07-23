package cluster

type DeleteClusterCommand struct {
	ID string `json:"id" validate:"required"`
}

func NewDeleteClusterCommand(id string) *DeleteClusterCommand {
	return &DeleteClusterCommand{
		ID: id,
	}
}
