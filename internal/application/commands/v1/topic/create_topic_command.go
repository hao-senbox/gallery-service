package topic

import (
	"gallery-service/internal/domain/models"
)

type CreateTopicCommand struct {
	FileName       string                       `json:"file_name"`
	IsPublished    bool                         `json:"is_published"`
	LanguageConfig []models.TopicLanguageConfig `json:"language_config"`
}

func NewCreateTopicCommand(
	fileName string,
	isPublished bool,
	languageConfig []models.TopicLanguageConfig,
) *CreateTopicCommand {
	return &CreateTopicCommand{
		FileName:       fileName,
		IsPublished:    isPublished,
		LanguageConfig: languageConfig,
	}
}
