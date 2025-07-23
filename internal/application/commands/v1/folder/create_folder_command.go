package folder

type CreateFolderCommand struct {
	FolderName         string
	FolderThumbnailKey string
	FolderThumbnailURL string
	ParentID           *string
}

func NewCreateFolderCommand(
	folderName string,
	folderThumbnailKey string,
	folderThumbnailURL string,
	parentID *string,
) *CreateFolderCommand {
	return &CreateFolderCommand{
		FolderName:         folderName,
		FolderThumbnailKey: folderThumbnailKey,
		FolderThumbnailURL: folderThumbnailURL,
		ParentID:           parentID,
	}
}
