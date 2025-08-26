package topic

import (
	"gallery-service/internal/domain/models"
)

type CreateTopicCommand struct {
	TopicName      string                       `json:"topic_name"`
	IsPublished    bool                         `json:"is_published"`
	LanguageConfig []models.TopicLanguageConfig `json:"language_config"`
}

func NewCreateTopicCommand(
	topicName string,
	isPublished bool,
	languageConfig []models.TopicLanguageConfig,
) *CreateTopicCommand {
	return &CreateTopicCommand{
		TopicName:      topicName,
		IsPublished:    isPublished,
		LanguageConfig: languageConfig,
	}
}
