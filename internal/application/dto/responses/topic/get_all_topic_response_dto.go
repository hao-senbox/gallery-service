package topic

import (
	"gallery-service/internal/application/dto/responses"
	"gallery-service/internal/domain/models"
)

type GetAllTopicResponseDto struct {
	Pagination responses.Pagination  `json:"pagination"`
	Topics     []GetTopicResponseDto `json:"topics"`
}

type GetTopicResponseDto struct {
	ID        string `json:"id"`
	TopicName string `json:"topic_name"`
	ImageKey  string `json:"image_key"`
	ImageURL  string `json:"image_url"`
	Images    []models.TopicImageConfig
}
