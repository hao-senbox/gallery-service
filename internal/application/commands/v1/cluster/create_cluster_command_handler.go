package cluster

import (
	"context"
	"gallery-service/internal/domain/models"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/kafka"
	"gallery-service/pkg/zap"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CreateClusterCommandHandler interface {
	Handle(ctx context.Context, command *CreateClusterCommand) (*string, error)
}

type createClusterHandler struct {
	cfg         kafka.Config
	log         zap.Logger
	clusterRepo repository.ClusterRepository
	folderRepo  repository.FolderRepository
}

func NewCreateClusterHandler(
	cfg kafka.Config,
	log zap.Logger,
	clusterRepo repository.ClusterRepository,
	folderRepo repository.FolderRepository,
) *createClusterHandler {
	return &createClusterHandler{
		cfg:         cfg,
		log:         log,
		clusterRepo: clusterRepo,
		folderRepo:  folderRepo,
	}
}

func (c *createClusterHandler) Handle(ctx context.Context, command *CreateClusterCommand) (*string, error) {
	id := primitive.NewObjectID()
	folderID, err := primitive.ObjectIDFromHex(command.FolderID)
	if err != nil {
		return nil, errors.New("invalid folder id")
	}

	// check if folder exist
	exist, err := c.folderRepo.Exists(ctx, bson.M{"_id": folderID})
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("folder not found")
	}

	cluster := models.Cluster{
		ID:             id,
		ClusterName:    command.ClusterName,
		Title:          command.Title,
		Note:           command.Note,
		Image:          command.Image,
		LanguageConfig: command.LanguageConfig,
		FolderID:       folderID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Save to database
	clusterID, err := c.clusterRepo.Insert(ctx, &cluster)
	if err != nil {
		return nil, err
	}

	return &clusterID, nil
}
