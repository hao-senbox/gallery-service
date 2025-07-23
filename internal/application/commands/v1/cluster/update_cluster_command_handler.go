package cluster

import (
	"context"
	"gallery-service/internal/domain/models"
	"gallery-service/internal/domain/repository"
	"gallery-service/pkg/zap"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UpdateClusterCommandHandler interface {
	Handle(ctx context.Context, command *UpdateClusterCommand) error
}

type updateClusterHandler struct {
	log         zap.Logger
	clusterRepo repository.ClusterRepository
	folderRepo  repository.FolderRepository
}

func NewUpdateClusterHandler(
	log zap.Logger,
	clusterRepo repository.ClusterRepository,
	folderRepo repository.FolderRepository,
) *updateClusterHandler {
	return &updateClusterHandler{
		log:         log,
		clusterRepo: clusterRepo,
		folderRepo:  folderRepo,
	}
}

func (u *updateClusterHandler) Handle(ctx context.Context, command *UpdateClusterCommand) error {
	cluster, err := u.clusterRepo.GetByID(ctx, command.ID)
	if err != nil {
		return err
	}

	folderID, err := primitive.ObjectIDFromHex(command.FolderID)
	if err != nil {
		return errors.New("invalid folder id")
	}

	// check if cluster exist
	exist, err := u.folderRepo.Exists(ctx, bson.M{"_id": folderID})
	if err != nil {
		return err
	}

	if !exist {
		return errors.New("folder not found")
	}

	t := models.Cluster{
		ID:             cluster.ID,
		ClusterName:    command.ClusterName,
		Title:          command.Title,
		Note:           command.Note,
		Image:          command.Image,
		LanguageConfig: command.LanguageConfig,
		FolderID:       folderID,
		CreatedAt:      cluster.CreatedAt,
		UpdatedAt:      time.Now(),
	}

	// Save to database
	return u.clusterRepo.Update(ctx, &t)
}
