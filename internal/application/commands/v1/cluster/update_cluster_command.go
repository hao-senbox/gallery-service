package cluster

import "gallery-service/internal/domain/models"

type UpdateClusterCommand struct {
	ID             string
	ClusterName    string
	Title          string
	Note           string
	Image          models.ImageConfig
	LanguageConfig []models.LanguageConfig
	FolderID       string
}

func NewUpdateClusterCommand(
	id string,
	clusterName string,
	title string,
	note string,
	image models.ImageConfig,
	languageConfig []models.LanguageConfig,
	folderID string,
) *UpdateClusterCommand {
	return &UpdateClusterCommand{
		ID:             id,
		ClusterName:    clusterName,
		Title:          title,
		Note:           note,
		Image:          image,
		LanguageConfig: languageConfig,
		FolderID:       folderID,
	}
}
