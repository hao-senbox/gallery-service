package topic

import (
	"gallery-service/internal/domain/models"
	"time"
)

type GetAllTopicResponseDto struct {
	//Pagination responses.Pagination  `json:"pagination"`
	Topics []GetTopicResponseDto `json:"topics"`
}

type GetTopicResponseDto struct {
	ID             string                       `json:"id"`
	FileName       string                       `json:"file_name"`
	IsPublished    bool                         `json:"is_published"`
	LanguageConfig []models.TopicLanguageConfig `json:"language_config"`
	CreatedAt      time.Time                    `json:"created_at"`
	UpdatedAt      time.Time                    `json:"updated_at"`
}
