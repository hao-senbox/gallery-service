package folder

type UpdateFolderReqDto struct {
	ID                 string  `json:"id" validate:"required"`
	FolderName         string  `json:"folder_name" validate:"required"`
	FolderThumbnailKey string  `json:"folder_thumbnail_key" validate:"required"`
	FolderThumbnailURL string  `json:"folder_thumbnail_url" validate:"required"`
	ParentID           *string `json:"parent_id"`
}
