package mappers

import (
	"gallery-service/internal/application/dto/responses/cluster"
	"gallery-service/internal/domain/models"
)

func GetAllClustersFromModel(c *models.Cluster) cluster.GetClusterResponseDto {
	return cluster.GetClusterResponseDto{
		ID:          c.ID.Hex(),
		ClusterName: c.ClusterName,
		ImageKey:    c.Image.ImageKey,
		ImageURL:    c.Image.ImageURL,
	}
}

func GetAllClustersFromModels(clusters []*models.Cluster) []cluster.GetClusterResponseDto {
	res := make([]cluster.GetClusterResponseDto, 0, len(clusters))
	for _, c := range clusters {
		res = append(res, GetAllClustersFromModel(c))
	}
	return res
}
