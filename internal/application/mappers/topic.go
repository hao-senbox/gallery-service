package mappers

import (
	"gallery-service/internal/application/dto/responses/topic"
	"gallery-service/internal/domain/models"
)

func GetAllTopicsFromModel(c *models.Topic) topic.GetTopicResponseDto {
	return topic.GetTopicResponseDto{
		ID:        c.ID.Hex(),
		TopicName: c.TopicName,
		ImageKey:  "",
		ImageURL:  "",
		Images:    c.Images,
	}
}

func GetAllTopicsFromModels(topics []*models.Topic) []topic.GetTopicResponseDto {
	res := make([]topic.GetTopicResponseDto, 0, len(topics))
	for _, c := range topics {
		res = append(res, GetAllTopicsFromModel(c))
	}
	return res
}
