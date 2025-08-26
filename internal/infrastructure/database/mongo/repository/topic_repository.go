package repository

import (
	"context"
	"fmt"
	"gallery-service/config"
	"gallery-service/internal/application/dto/responses/topic"
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

type topicRepository struct {
	log zap.Logger
	cfg *config.Config
	db  *mongo.Client
}

var (
	topicRepo *topicRepository
)

func NewTopicRepository(log zap.Logger, cfg *config.Config, db *mongo.Client) *topicRepository {
	if topicRepo == nil {
		topicRepo = &topicRepository{log: log, cfg: cfg, db: db}
	}

	return topicRepo
}

func (p *topicRepository) Insert(ctx context.Context, topic *models.Topic) (string, error) {
	insertResult, err := p.getTopicsCollection().InsertOne(ctx, topic, &options.InsertOneOptions{})
	if err != nil {
		p.log.Errorf("(topicRepository.Insert) Error inserting topic: %v", err)
		return "", err
	}

	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (p *topicRepository) Update(ctx context.Context, topic *models.Topic) error {
	req := bson.M{
		"file_name":       topic.TopicName,
		"is_published":    topic.IsPublished,
		"language_config": topic.LanguageConfig,
		"created_at":      topic.CreatedAt,
		"updated_at":      topic.UpdatedAt,
	}

	result, err := p.getTopicsCollection().UpdateOne(
		ctx,
		bson.M{"_id": topic.ID},
		bson.M{"$set": req},
	)
	if err != nil {
		return fmt.Errorf("(topicRepository.Update) failed to update: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("(topicRepository.Update) no topic found with ID: %s", topic.ID.Hex())
	}

	return nil
}

func (p *topicRepository) GetAll(ctx context.Context, pq *utils.Pagination) (*topic.GetAllTopicResponseDto, error) {
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

	// Perform the search query on the topics collection
	cursor, err := p.getTopicsCollection().Aggregate(ctx, aggPipeline)
	if err != nil {
		p.log.Errorf("(topicRepository.GetAll) Error fetching topics: %v", err)
		return nil, errors.Wrap(err, "mongoRepository.Find")
	}
	defer cursor.Close(ctx)

	// Create a slice to hold the search results
	var topics []*models.Topic
	if err := cursor.All(ctx, &topics); err != nil {
		p.log.Errorf("(topicRepository.GetAll) Error fetching topics: %v", err)
		return nil, errors.Wrap(err, "cursor.All")
	}

	//Prepare pagination response
	//totalCount := len(topics)

	return &topic.GetAllTopicResponseDto{
		// Pagination: responses.Pagination{
		// 	TotalCount: int64(totalCount),
		// 	TotalPages: int64(pq.GetTotalPages(totalCount)),
		// 	Page:       int64(pq.GetPage()),
		// 	Size:       int64(pq.GetSize()),
		// 	HasMore:    pq.GetHasMore(totalCount),
		// },
		Topics: mappers.GetTopicsFromModels(topics),
	}, nil
}

func (p *topicRepository) GetByID(ctx context.Context, topicID string) (*models.Topic, error) {
	p.log.Infof("(topicRepository.GetByID) topicID: %s", topicID)
	objectId, _ := primitive.ObjectIDFromHex(topicID)
	var topic models.Topic
	if err := p.getTopicsCollection().FindOne(ctx, bson.M{"_id": objectId}).Decode(&topic); err != nil {
		p.log.Errorf("(topicRepository.GetByID) Error fetching topic: %v", err)
		return nil, err
	}

	return &topic, nil
}

func (p *topicRepository) Search(ctx context.Context, query map[string]interface{}, pq *utils.Pagination) (*topic.GetAllTopicResponseDto, error) {
	// Prepare pagination options
	skip := int64((pq.Page - 1) * pq.Size)
	limit := int64(pq.Size)

	queryFilter := bson.M{}

	if query["keyword"] != nil {
		queryFilter = bson.M{
			"$or": []bson.M{
				{"topic_name": bson.M{"$regex": primitive.Regex{Pattern: query["keyword"].(string), Options: "i"}}},
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
		p.log.Warnf("(topicRepository.Search) Query filter is empty, skipping $match stage.")
	}

	p.log.Infof("(topicRepository.Search) Searching for query: %v", aggPipeline)

	// Perform the search query on the topics collection
	cursor, err := p.getTopicsCollection().Aggregate(ctx, aggPipeline)
	if err != nil {
		p.log.Errorf("(topicRepository.Search) Error fetching topcs: %v", err)
		return nil, errors.Wrap(err, "mongoRepository.Find")
	}
	defer cursor.Close(ctx)

	// Create a slice to hold the search results
	var topics []*models.Topic
	if err := cursor.All(ctx, &topics); err != nil {
		p.log.Errorf("(topicRepository.Search) Error fetching topics: %v", err)
		return nil, errors.Wrap(err, "cursor.All")
	}

	// Prepare pagination response
	//totalCount := len(topics)

	return &topic.GetAllTopicResponseDto{
		// Pagination: responses.Pagination{
		// 	TotalCount: int64(totalCount),
		// 	TotalPages: int64(pq.GetTotalPages(totalCount)),
		// 	Page:       int64(pq.GetPage()),
		// 	Size:       int64(pq.GetSize()),
		// 	HasMore:    pq.GetHasMore(totalCount),
		// },
		Topics: mappers.GetTopicsFromModels(topics),
	}, nil
}

func (p *topicRepository) Delete(ctx context.Context, topicID string) (bool, error) {
	objectId, _ := primitive.ObjectIDFromHex(topicID)
	res, err := p.getTopicsCollection().DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		p.log.Errorf("(topicRepository.Delete) Error deleting topic: %v", err)
		return false, err
	}

	return res.DeletedCount > 0, nil
}

func (p *topicRepository) Exists(ctx context.Context, query map[string]interface{}) (bool, error) {
	count, err := p.getTopicsCollection().CountDocuments(ctx, query)
	if err != nil {
		p.log.Errorf("(topicRepository.Exists) Error counting topics: %v", err)
		return false, err
	}

	return count > 0, nil
}

func (p *topicRepository) GetAll4App(ctx context.Context) (*topic.GetAllTopicForAppResponseDto, error) {
	cur, err := p.getTopicsCollection().Find(ctx, bson.M{})
	if err != nil {
		p.log.Errorf("(topicRepository.GetAll4App) Error fetching topics: %v", err)
		return nil, errors.Wrap(err, "find topics failed")
	}
	defer cur.Close(ctx)

	var topics []*models.Topic
	if err := cur.All(ctx, &topics); err != nil {
		p.log.Errorf("(topicRepository.GetAll4App) Error decoding topics: %v", err)
		return nil, errors.Wrap(err, "cursor.All")
	}

	var topicForAppList []topic.TopicForAppResponseDto
	for _, t := range topics {
		topicForAppList = append(topicForAppList, topic.TopicForAppResponseDto{
			ID:        t.ID.Hex(),
			TopicName: t.TopicName,
		})
	}

	return &topic.GetAllTopicForAppResponseDto{
		Topics: topicForAppList,
	}, nil
}

func (p *topicRepository) getTopicsCollection() *mongo.Collection {
	return p.db.Database(p.cfg.Mongo.Db).Collection(p.cfg.Mongo.Collections.Topic)
}
