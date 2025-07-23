package cluster

import (
	"gallery-service/internal/domain/models"
)

type CreateClusterCommand struct {
	ClusterName    string
	Title          string
	Note           string
	Image          models.ImageConfig
	LanguageConfig []models.LanguageConfig
	FolderID       string
}

func NewCreateClusterCommand(
	clusterName string,
	title string,
	note string,
	image models.ImageConfig,
	languageConfig []models.LanguageConfig,
	folderID string,
) *CreateClusterCommand {
	return &CreateClusterCommand{
		ClusterName:    clusterName,
		Title:          title,
		Note:           note,
		Image:          image,
		LanguageConfig: languageConfig,
		FolderID:       folderID,
	}
}
