package topic

import "gallery-service/internal/application/dto/responses"

type GetAllTopicResponseDto struct {
	Pagination responses.Pagination  `json:"pagination"`
	Topics     []GetTopicResponseDto `json:"topics"`
}

type GetTopicResponseDto struct {
	ID          string `json:"id"`
	ClusterName string `json:"cluster_name"`
	ImageKey    string `json:"image_key"`
	ImageURL    string `json:"image_url"`
}
