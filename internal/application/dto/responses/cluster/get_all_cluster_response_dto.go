package cluster

import "gallery-service/internal/application/dto/responses"

type GetAllClusterResponseDto struct {
	Pagination responses.Pagination    `json:"pagination"`
	Clusters   []GetClusterResponseDto `json:"clusters"`
}

type GetClusterResponseDto struct {
	ID          string `json:"id"`
	ClusterName string `json:"cluster_name"`
	ImageKey    string `json:"image_key"`
	ImageURL    string `json:"image_url"`
}
