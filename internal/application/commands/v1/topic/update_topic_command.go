package topic

import "gallery-service/internal/domain/models"

type UpdateTopicCommand struct {
	ID             string
	TopicName      string                       `json:"topic_name"`
	IsPublished    bool                         `json:"is_published"`
	LanguageConfig []models.TopicLanguageConfig `json:"language_config"`
}

func NewUpdateTopicCommand(
	id string,
	topicName string,
	isPublished bool,
	languageConfig []models.TopicLanguageConfig,
) *UpdateTopicCommand {
	return &UpdateTopicCommand{
		ID:             id,
		TopicName:      topicName,
		IsPublished:    isPublished,
		LanguageConfig: languageConfig,
	}
}
