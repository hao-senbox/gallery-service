package cluster

import (
	"gallery-service/internal/domain/models"
)

type CreateClusterReqDto struct {
	ClusterName    string                  `json:"cluster_name" validate:"required"`
	Title          string                  `json:"title" validate:"required"`
	Note           string                  `json:"note" validate:"required"`
	Image          models.ImageConfig      `json:"image" validate:"required"`
	LanguageConfig []models.LanguageConfig `json:"language_config" validate:"required"`
	FolderID       string                  `json:"folder_id" validate:"required"`
}
