package topic

import "gallery-service/internal/domain/models"

type CreateTopicReqDto struct {
	TopicName      string                       `json:"topic_name" validate:"required"`
	IsPublished    bool                         `json:"is_published"`
	LanguageConfig []models.TopicLanguageConfig `json:"language_config" validate:"required"`
}
