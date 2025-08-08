package topic

import "gallery-service/internal/domain/models"

type UpdateTopicCommand struct {
	ID             string
	FileName       string                       `json:"file_name"`
	IsPublished    bool                         `json:"is_published"`
	LanguageConfig []models.TopicLanguageConfig `json:"language_config"`
}

func NewUpdateTopicCommand(
	id string,
	fileName string,
	isPublished bool,
	languageConfig []models.TopicLanguageConfig,
) *UpdateTopicCommand {
	return &UpdateTopicCommand{
		ID:             id,
		FileName:       fileName,
		IsPublished:    isPublished,
		LanguageConfig: languageConfig,
	}
}
