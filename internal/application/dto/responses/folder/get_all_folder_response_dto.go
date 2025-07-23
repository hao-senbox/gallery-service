package folder

import "gallery-service/internal/application/dto/responses"

type GetAllFolderResponseDto struct {
	Pagination responses.Pagination   `json:"pagination"`
	Folders    []GetFolderResponseDto `json:"folders"`
}

type GetFolderResponseDto struct {
	ID                 string `json:"id"`
	FolderName         string `json:"folder_name"`
	FolderThumbnailKey string `json:"folder_thumbnail_key"`
	FolderThumbnailURL string `json:"folder_thumbnail_url"`
	ParentID           string `json:"parent_id"`
}
