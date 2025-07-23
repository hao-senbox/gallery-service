package folder

type DeleteFolderCommand struct {
	ID string `json:"id" validate:"required"`
}

func NewDeleteFolderCommand(id string) *DeleteFolderCommand {
	return &DeleteFolderCommand{
		ID: id,
	}
}
