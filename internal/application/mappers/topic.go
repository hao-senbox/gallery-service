package mappers

import (
	"gallery-service/internal/application/dto/responses/topic"
	"gallery-service/internal/domain/models"
)

func GetTopicFromModel(c *models.Topic) topic.GetTopicResponseDto {
	return topic.GetTopicResponseDto{
		ID:             c.ID.Hex(),
		TopicName:      c.TopicName,
		IsPublished:    c.IsPublished,
		LanguageConfig: c.LanguageConfig,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}

func GetTopicsFromModels(topics []*models.Topic) []topic.GetTopicResponseDto {
	res := make([]topic.GetTopicResponseDto, 0, len(topics))
	for _, c := range topics {
		res = append(res, GetTopicFromModel(c))
	}
	return res
}
