package repository

import (
	"context"
	"gallery-service/internal/application/dto/responses/cluster"
	"gallery-service/internal/application/dto/responses/folder"
	"gallery-service/internal/domain/models"
	"gallery-service/pkg/utils"
)

type ClusterRepository interface {
	Insert(ctx context.Context, cluster *models.Cluster) (string, error)
	Update(ctx context.Context, cluster *models.Cluster) error
	GetAll(ctx context.Context, pq *utils.Pagination) (*cluster.GetAllClusterResponseDto, error)
	GetAllByFolderID(ctx context.Context, folderID string, pq *utils.Pagination) (*cluster.GetAllClusterResponseDto, error)
	GetByID(ctx context.Context, clusterID string) (*models.Cluster, error)
	Search(ctx context.Context, query map[string]interface{}, pq *utils.Pagination) (*cluster.GetAllClusterResponseDto, error)
	Delete(ctx context.Context, clusterID string) (bool, error)
	Exists(ctx context.Context, query map[string]interface{}) (bool, error)
}

type FolderRepository interface {
	Insert(ctx context.Context, folder *models.Folder) (string, error)
	Update(ctx context.Context, folder *models.Folder) error
	GetAll(ctx context.Context, pq *utils.Pagination) (*folder.GetAllFolderResponseDto, error)
	GetByID(ctx context.Context, folderID string) (*models.Folder, error)
	Search(ctx context.Context, query map[string]interface{}, pq *utils.Pagination) (*folder.GetAllFolderResponseDto, error)
	Delete(ctx context.Context, folderID string) (bool, error)
	Exists(ctx context.Context, query map[string]interface{}) (bool, error)
}
