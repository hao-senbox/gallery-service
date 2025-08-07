package topic

import (
	"gallery-service/internal/domain/models"
)

type CreateTopicReqDto struct {
	TopicName      string                     `json:"topic_name" validate:"required"`
	Title          string                     `json:"title" validate:"required"`
	Note           string                     `json:"note" validate:"required"`
	Images         []models.TopicImageConfig  `json:"images" validate:"required"`
	LanguageConfig models.TopicLanguageConfig `json:"language_config" validate:"required"`
}
