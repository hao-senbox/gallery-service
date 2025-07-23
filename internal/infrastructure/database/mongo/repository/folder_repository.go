package repository

import (
	"context"
	"fmt"
	"gallery-service/config"
	"gallery-service/internal/application/dto/responses"
	"gallery-service/internal/application/dto/responses/folder"
	"gallery-service/internal/application/mappers"
	"gallery-service/internal/domain/models"
	"gallery-service/pkg/utils"
	"gallery-service/pkg/zap"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type folderRepository struct {
	log zap.Logger
	cfg *config.Config
	db  *mongo.Client
}

var (
	folderRepo *folderRepository
)

func NewFolderRepository(log zap.Logger, cfg *config.Config, db *mongo.Client) *folderRepository {
	if folderRepo == nil {
		folderRepo = &folderRepository{log: log, cfg: cfg, db: db}
	}

	return folderRepo
}

func (c *folderRepository) Insert(ctx context.Context, folder *models.Folder) (string, error) {
	insertResult, err := c.getFoldersCollection().InsertOne(ctx, folder, &options.InsertOneOptions{})
	if err != nil {
		c.log.Errorf("(FolderRepository.Insert) Error inserting cluster: %v", err)
		return "", err
	}

	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (c *folderRepository) Update(ctx context.Context, folder *models.Folder) error {
	req := make(bson.M)
	req["folder_name"] = folder.FolderName
	req["folder_thumbnail_key"] = folder.FolderThumbnailKey
	req["folder_thumbnail_url"] = folder.FolderThumbnailURL
	req["parent_id"] = folder.ParentID

	result, err := c.getFoldersCollection().UpdateOne(
		ctx,
		bson.M{"_id": folder.ID},
		bson.M{"$set": req})

	if err != nil {
		return fmt.Errorf("(FolderRepository.Update) failed to update: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("(FolderRepository.Update) no folder found with ID: %s", folder.ID.Hex())
	}

	return nil
}

func (c *folderRepository) GetAll(ctx context.Context, pq *utils.Pagination) (*folder.GetAllFolderResponseDto, error) {
	if pq.Page <= 0 {
		pq.Page = 1
	}

	if pq.Size <= 0 {
		pq.Size = 10000
	}

	// Prepare pagination options
	skip := int64((pq.Page - 1) * pq.Size)
	limit := int64(pq.Size)

	aggPipeline := []bson.M{
		{
			"$skip": skip,
		},
		{
			"$limit": limit,
		},
	}

	// Perform the search query on the folders collection
	cursor, err := c.getFoldersCollection().Aggregate(ctx, aggPipeline)
	if err != nil {
		c.log.Errorf("(FolderRepository.Search) Error fetching folders: %v", err)
		return nil, errors.Wrap(err, "mongoRepository.Find")
	}
	defer cursor.Close(ctx)

	// Create a slice to hold the search results
	var folders []*models.Folder
	if err := cursor.All(ctx, &folders); err != nil {
		c.log.Errorf("(FolderRepository.Search) Error fetching folders: %v", err)
		return nil, errors.Wrap(err, "cursor.All")
	}

	// Prepare pagination response
	totalCount := len(folders)

	return &folder.GetAllFolderResponseDto{
		Pagination: responses.Pagination{
			TotalCount: int64(totalCount),
			TotalPages: int64(pq.GetTotalPages(totalCount)),
			Page:       int64(pq.GetPage()),
			Size:       int64(pq.GetSize()),
			HasMore:    pq.GetHasMore(totalCount),
		},
		Folders: mappers.GetAllFoldersFromModels(folders),
	}, nil
}

func (c *folderRepository) GetByID(ctx context.Context, folderID string) (*models.Folder, error) {
	c.log.Infof("(FolderRepository.GetByID) FolderID: %s", folderID)
	objectId, _ := primitive.ObjectIDFromHex(folderID)
	var folder models.Folder
	if err := c.getFoldersCollection().FindOne(ctx, bson.M{"_id": objectId}).Decode(&folder); err != nil {
		c.log.Errorf("(FolderRepository.GetByID) Error fetching cluster: %v", err)
		return nil, err
	}

	return &folder, nil
}

func (c *folderRepository) Search(ctx context.Context, query map[string]interface{}, pq *utils.Pagination) (*folder.GetAllFolderResponseDto, error) {
	// Prepare pagination options
	skip := int64((pq.Page - 1) * pq.Size)
	limit := int64(pq.Size)

	queryFilter := bson.M{}

	if query["keyword"] != nil {
		queryFilter = bson.M{
			"$or": []bson.M{
				{"folder_name": bson.M{"$regex": primitive.Regex{Pattern: query["keyword"].(string), Options: "i"}}},
			},
		}
	}

	aggPipeline := []bson.M{
		{
			"$skip": skip,
		},
		{
			"$limit": limit,
		},
	}

	if len(queryFilter) > 0 {
		aggPipeline = append(aggPipeline, bson.M{
			"$match": queryFilter,
		})
	} else {
		c.log.Warnf("(TaskRepository.Search) Query filter is empty, skipping $match stage.")
	}

	c.log.Infof("(TaskRepository.Search) Searching for query: %v", aggPipeline)

	// Perform the search query on the folders collection
	cursor, err := c.getFoldersCollection().Aggregate(ctx, aggPipeline)
	if err != nil {
		c.log.Errorf("(FolderRepository.Search) Error fetching folders: %v", err)
		return nil, errors.Wrap(err, "mongoRepository.Find")
	}
	defer cursor.Close(ctx)

	// Create a slice to hold the search results
	var folders []*models.Folder
	if err := cursor.All(ctx, &folders); err != nil {
		c.log.Errorf("(FolderRepository.Search) Error fetching folders: %v", err)
		return nil, errors.Wrap(err, "cursor.All")
	}

	// Prepare pagination response
	totalCount := len(folders)

	return &folder.GetAllFolderResponseDto{
		Pagination: responses.Pagination{
			TotalCount: int64(totalCount),
			TotalPages: int64(pq.GetTotalPages(totalCount)),
			Page:       int64(pq.GetPage()),
			Size:       int64(pq.GetSize()),
			HasMore:    pq.GetHasMore(totalCount),
		},
		Folders: mappers.GetAllFoldersFromModels(folders),
	}, nil
}

func (c *folderRepository) Delete(ctx context.Context, folderID string) (bool, error) {
	objectId, _ := primitive.ObjectIDFromHex(folderID)
	res, err := c.getFoldersCollection().DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		c.log.Errorf("(FolderRepository.Delete) Error deleting cluster: %v", err)
		return false, err
	}

	return res.DeletedCount > 0, nil
}

func (c *folderRepository) Exists(ctx context.Context, query map[string]interface{}) (bool, error) {
	count, err := c.getFoldersCollection().CountDocuments(ctx, query)
	if err != nil {
		c.log.Errorf("(FolderRepository.Exists) Error counting folders: %v", err)
		return false, err
	}

	return count > 0, nil
}

func (c *folderRepository) getFoldersCollection() *mongo.Collection {
	return c.db.Database(c.cfg.Mongo.Db).Collection(c.cfg.Mongo.Collections.Folder)
}
