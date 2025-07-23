package folder

type UpdateFolderCommand struct {
	ID                 string
	FolderName         string
	FolderThumbnailKey string
	FolderThumbnailURL string
	ParentID           *string
}

func NewUpdateFolderCommand(
	id string,
	folderName string,
	folderThumbnailKey string,
	folderThumbnailURL string,
	parentID *string,
) *UpdateFolderCommand {
	return &UpdateFolderCommand{
		ID:                 id,
		FolderName:         folderName,
		FolderThumbnailKey: folderThumbnailKey,
		FolderThumbnailURL: folderThumbnailURL,
		ParentID:           parentID,
	}
}
