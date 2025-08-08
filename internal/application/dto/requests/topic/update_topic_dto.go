package topic

import "gallery-service/internal/domain/models"

type UpdateTopicReqDto struct {
	ID             string                       `json:"id" validate:"required"`
	FileName       string                       `json:"file_name" validate:"required"`
	IsPublished    bool                         `json:"is_published"`
	LanguageConfig []models.TopicLanguageConfig `json:"language_config" validate:"required"`
}
