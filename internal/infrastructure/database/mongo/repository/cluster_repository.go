package repository

import (
	"context"
	"fmt"
	"gallery-service/config"
	"gallery-service/internal/application/dto/responses"
	"gallery-service/internal/application/dto/responses/cluster"
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

type clusterRepository struct {
	log zap.Logger
	cfg *config.Config
	db  *mongo.Client
}

var (
	clusterRepo *clusterRepository
)

func NewClusterRepository(log zap.Logger, cfg *config.Config, db *mongo.Client) *clusterRepository {
	if clusterRepo == nil {
		clusterRepo = &clusterRepository{log: log, cfg: cfg, db: db}
	}

	return clusterRepo
}

func (p *clusterRepository) Insert(ctx context.Context, cluster *models.Cluster) (string, error) {
	insertResult, err := p.getClustersCollection().InsertOne(ctx, cluster, &options.InsertOneOptions{})
	if err != nil {
		p.log.Errorf("(ClusterRepository.Insert) Error inserting cluster: %v", err)
		return "", err
	}

	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (p *clusterRepository) Update(ctx context.Context, cluster *models.Cluster) error {
	req := make(bson.M)
	req["cluster_name"] = cluster.ClusterName
	req["title"] = cluster.Title
	req["note"] = cluster.Note
	req["image"] = cluster.Image
	req["language_config"] = cluster.LanguageConfig
	req["folder_id	"] = cluster.FolderID
	req["created_at"] = cluster.CreatedAt
	req["updated_at"] = cluster.UpdatedAt

	result, err := p.getClustersCollection().UpdateOne(
		ctx,
		bson.M{"_id": cluster.ID},
		bson.M{"$set": req})

	if err != nil {
		return fmt.Errorf("(ClusterRepository.Update) failed to update: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("(ClusterRepository.Update) no cluster found with ID: %s", cluster.ID.Hex())
	}

	return nil
}

func (p *clusterRepository) GetAll(ctx context.Context, pq *utils.Pagination) (*cluster.GetAllClusterResponseDto, error) {
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

	// Perform the search query on the clusters collection
	cursor, err := p.getClustersCollection().Aggregate(ctx, aggPipeline)
	if err != nil {
		p.log.Errorf("(ClusterRepository.GetAll) Error fetching clusters: %v", err)
		return nil, errors.Wrap(err, "mongoRepository.Find")
	}
	defer cursor.Close(ctx)

	// Create a slice to hold the search results
	var clusters []*models.Cluster
	if err := cursor.All(ctx, &clusters); err != nil {
		p.log.Errorf("(ClusterRepository.GetAll) Error fetching clusters: %v", err)
		return nil, errors.Wrap(err, "cursor.All")
	}

	//Prepare pagination response
	totalCount := len(clusters)

	return &cluster.GetAllClusterResponseDto{
		Pagination: responses.Pagination{
			TotalCount: int64(totalCount),
			TotalPages: int64(pq.GetTotalPages(totalCount)),
			Page:       int64(pq.GetPage()),
			Size:       int64(pq.GetSize()),
			HasMore:    pq.GetHasMore(totalCount),
		},
		Clusters: mappers.GetAllClustersFromModels(clusters),
	}, nil
}

func (p *clusterRepository) GetAllByFolderID(ctx context.Context, folderID string, pq *utils.Pagination) (*cluster.GetAllClusterResponseDto, error) {
	// Debug: Log input parameters
	p.log.Infof("(ClusterRepository.GetAllByFolderID) Input - folderID: '%s', page: %d, size: %d", folderID, pq.Page, pq.Size)

	if pq.Page <= 0 {
		pq.Page = 1
	}

	if pq.Size <= 0 {
		pq.Size = 10000
	}
	folderObjectID, err := primitive.ObjectIDFromHex(folderID)
	if err != nil {
		p.log.Errorf("(ClusterRepository.GetAllByFolderID) Error parsing folderID: %v", err)
		return nil, errors.Wrap(err, "primitive.ObjectIDFromHex")
	}
	// Debug: Log processed pagination
	p.log.Infof("(ClusterRepository.GetAllByFolderID) Processed pagination - page: %d, size: %d", pq.Page, pq.Size)

	// Prepare pagination options
	skip := int64((pq.Page - 1) * pq.Size)
	limit := int64(pq.Size)

	// Debug: Log skip and limit values
	p.log.Infof("(ClusterRepository.GetAllByFolderID) Skip: %d, Limit: %d", skip, limit)

	aggPipeline := []bson.M{
		{
			"$match": bson.M{
				"folder_id": folderObjectID,
			},
		},
		{
			"$skip": skip,
		},
		{
			"$limit": limit,
		},
	}

	// Debug: Log aggregation pipeline
	p.log.Infof("(ClusterRepository.GetAllByFolderID) Aggregation pipeline: %+v", aggPipeline)

	// Debug: First, let's check if there are ANY documents in the collection
	totalDocsInCollection, err := p.getClustersCollection().CountDocuments(ctx, bson.M{})
	if err != nil {
		p.log.Errorf("(ClusterRepository.GetAllByFolderID) Error counting total documents: %v", err)
	} else {
		p.log.Infof("(ClusterRepository.GetAllByFolderID) Total documents in clusters collection: %d", totalDocsInCollection)
	}

	// Debug: Check if there are documents with this folderID BEFORE pagination
	countBeforePagination, err := p.getClustersCollection().CountDocuments(ctx, bson.M{
		"folder_id": folderID,
	})
	if err != nil {
		p.log.Errorf("(ClusterRepository.GetAllByFolderID) Error counting documents with folderID: %v", err)
	} else {
		p.log.Infof("(ClusterRepository.GetAllByFolderID) Documents matching folderID '%s': %d", folderID, countBeforePagination)
	}

	// Debug: Let's also check for variations of the folderID (in case of data inconsistency)
	// Check for documents with similar folder_id patterns
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id": "$folder_id",
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"count": -1},
		},
		{
			"$limit": 10,
		},
	}

	cursor, err := p.getClustersCollection().Aggregate(ctx, pipeline)
	if err == nil {
		var folderStats []bson.M
		if err := cursor.All(ctx, &folderStats); err == nil {
			p.log.Infof("(ClusterRepository.GetAllByFolderID) Top folder_id values in collection: %+v", folderStats)
		}
		cursor.Close(ctx)
	}

	// Perform the search query on the clusters collection
	cursor, err = p.getClustersCollection().Aggregate(ctx, aggPipeline)
	if err != nil {
		p.log.Errorf("(ClusterRepository.GetAllByFolderID) Error fetching clusters: %v", err)
		return nil, errors.Wrap(err, "mongoRepository.Aggregate")
	}
	defer cursor.Close(ctx)

	var clusters []*models.Cluster
	if err := cursor.All(ctx, &clusters); err != nil {
		p.log.Errorf("(ClusterRepository.GetAllByFolderID) Error reading cursor: %v", err)
		return nil, errors.Wrap(err, "cursor.All")
	}

	// Debug: Log the number of clusters found
	p.log.Infof("(ClusterRepository.GetAllByFolderID) Found %d clusters after aggregation", len(clusters))

	// Debug: If we found clusters, log the first one (without sensitive data)
	if len(clusters) > 0 {
		p.log.Infof("(ClusterRepository.GetAllByFolderID) First cluster - ID: %s, FolderID: %s", 
			clusters[0].ID, clusters[0].FolderID)
	}

	// Count query for total (this is the same as countBeforePagination, but keeping for consistency)
	count, err := p.getClustersCollection().CountDocuments(ctx, bson.M{
		"folder_id": folderID,
	})
	if err != nil {
		p.log.Errorf("(ClusterRepository.GetAllByFolderID) Error counting clusters: %v", err)
		return nil, errors.Wrap(err, "CountDocuments")
	}

	// Debug: Log final results
	p.log.Infof("(ClusterRepository.GetAllByFolderID) Final results - Total count: %d, Clusters returned: %d", count, len(clusters))

	result := &cluster.GetAllClusterResponseDto{
		Pagination: responses.Pagination{
			TotalCount: count,
			TotalPages: int64(pq.GetTotalPages(int(count))),
			Page:       int64(pq.GetPage()),
			Size:       int64(pq.GetSize()),
			HasMore:    pq.GetHasMore(int(count)),
		},
		Clusters: mappers.GetAllClustersFromModels(clusters),
	}

	// Debug: Log the final response structure (without actual data)
	p.log.Infof("(ClusterRepository.GetAllByFolderID) Response - TotalCount: %d, TotalPages: %d, Page: %d, Size: %d, HasMore: %v, ClustersCount: %d", 
		result.Pagination.TotalCount, 
		result.Pagination.TotalPages, 
		result.Pagination.Page, 
		result.Pagination.Size, 
		result.Pagination.HasMore,
		len(result.Clusters))

	return result, nil
}

func (p *clusterRepository) GetByID(ctx context.Context, clusterID string) (*models.Cluster, error) {
	p.log.Infof("(clusterRepository.GetByID) ClusterID: %s", clusterID)
	objectId, _ := primitive.ObjectIDFromHex(clusterID)
	var cluster models.Cluster
	if err := p.getClustersCollection().FindOne(ctx, bson.M{"_id": objectId}).Decode(&cluster); err != nil {
		p.log.Errorf("(clusterRepository.GetByID) Error fetching cluster: %v", err)
		return nil, err
	}

	return &cluster, nil
}

func (p *clusterRepository) Search(ctx context.Context, query map[string]interface{}, pq *utils.Pagination) (*cluster.GetAllClusterResponseDto, error) {
	// Prepare pagination options
	skip := int64((pq.Page - 1) * pq.Size)
	limit := int64(pq.Size)

	queryFilter := bson.M{}

	if query["keyword"] != nil {
		queryFilter = bson.M{
			"$or": []bson.M{
				{"cluster_name": bson.M{"$regex": primitive.Regex{Pattern: query["keyword"].(string), Options: "i"}}},
				{"title": bson.M{"$regex": primitive.Regex{Pattern: query["keyword"].(string), Options: "i"}}},
				{"note": bson.M{"$regex": primitive.Regex{Pattern: query["keyword"].(string), Options: "i"}}},
			},
		}
	}

	p.log.Debugf("Starting %v", queryFilter)

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
		p.log.Warnf("(ClusterRepository.Search) Query filter is empty, skipping $match stage.")
	}

	p.log.Infof("(ClusterRepository.Search) Searching for query: %v", aggPipeline)

	// Perform the search query on the clusters collection
	cursor, err := p.getClustersCollection().Aggregate(ctx, aggPipeline)
	if err != nil {
		p.log.Errorf("(ClusterRepository.Search) Error fetching clusters: %v", err)
		return nil, errors.Wrap(err, "mongoRepository.Find")
	}
	defer cursor.Close(ctx)

	// Create a slice to hold the search results
	var clusters []*models.Cluster
	if err := cursor.All(ctx, &clusters); err != nil {
		p.log.Errorf("(ClusterRepository.Search) Error fetching clusters: %v", err)
		return nil, errors.Wrap(err, "cursor.All")
	}

	// Prepare pagination response
	totalCount := len(clusters)

	return &cluster.GetAllClusterResponseDto{
		Pagination: responses.Pagination{
			TotalCount: int64(totalCount),
			TotalPages: int64(pq.GetTotalPages(totalCount)),
			Page:       int64(pq.GetPage()),
			Size:       int64(pq.GetSize()),
			HasMore:    pq.GetHasMore(totalCount),
		},
		Clusters: mappers.GetAllClustersFromModels(clusters),
	}, nil
}

func (p *clusterRepository) Delete(ctx context.Context, clusterID string) (bool, error) {
	objectId, _ := primitive.ObjectIDFromHex(clusterID)
	res, err := p.getClustersCollection().DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		p.log.Errorf("(ClusterRepository.Delete) Error deleting cluster: %v", err)
		return false, err
	}

	return res.DeletedCount > 0, nil
}

func (p *clusterRepository) Exists(ctx context.Context, query map[string]interface{}) (bool, error) {
	count, err := p.getClustersCollection().CountDocuments(ctx, query)
	if err != nil {
		p.log.Errorf("(ClusterRepository.Exists) Error counting clusters: %v", err)
		return false, err
	}

	return count > 0, nil
}

func (p *clusterRepository) getClustersCollection() *mongo.Collection {
	return p.db.Database(p.cfg.Mongo.Db).Collection(p.cfg.Mongo.Collections.Cluster)
}
