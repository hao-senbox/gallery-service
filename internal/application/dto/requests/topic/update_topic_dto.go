package topic

import "gallery-service/internal/domain/models"

type UpdateTopicReqDto struct {
	ID             string                       `json:"id" validate:"required"`
	TopicName      string                       `json:"topic_name" validate:"required"`
	Title          string                       `json:"title" validate:"required"`
	Note           string                       `json:"note" validate:"required"`
	Images         []models.TopicImageConfig    `json:"images" validate:"required"`
	LanguageConfig []models.TopicLanguageConfig `json:"language_config" validate:"required"`
	FolderID       string                       `json:"folder_id" validate:"required"`
}
