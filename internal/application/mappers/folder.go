package mappers

import (
	"gallery-service/internal/application/dto/responses/folder"
	"gallery-service/internal/domain/models"
)

func GetAllFoldersFromModel(f *models.Folder) folder.GetFolderResponseDto {
	var parentID string
	if f.ParentID != nil {
		parentID = f.ParentID.Hex()
	}

	return folder.GetFolderResponseDto{
		ID:                 f.ID.Hex(),
		FolderName:         f.FolderName,
		FolderThumbnailKey: f.FolderThumbnailKey,
		FolderThumbnailURL: f.FolderThumbnailURL,
		ParentID:           parentID,
	}
}

func GetAllFoldersFromModels(folders []*models.Folder) []folder.GetFolderResponseDto {
	res := make([]folder.GetFolderResponseDto, 0, len(folders))
	for _, f := range folders {
		res = append(res, GetAllFoldersFromModel(f))
	}
	return res
}
