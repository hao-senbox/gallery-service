package folder

type Commands struct {
	CreateFolder CreateFolderCommandHandler
	UpdateFolder UpdateFolderCommandHandler
	DeleteFolder DeleteFolderCommandHandler
}

func NewFolderCommands(
	createFolder CreateFolderCommandHandler,
	updateFolder UpdateFolderCommandHandler,
	deleteFolder DeleteFolderCommandHandler,
) *Commands {
	return &Commands{
		CreateFolder: createFolder,
		UpdateFolder: updateFolder,
		DeleteFolder: deleteFolder,
	}
}
