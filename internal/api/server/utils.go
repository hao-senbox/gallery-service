package server

import (
	"context"
	"fmt"
	"gallery-service/config"
	serviceErrors "gallery-service/pkg/service_errors"
	"gallery-service/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

const (
	waitShotDownDuration = 3 * time.Second
)

func (s *server) mongoMigrationUp(ctx context.Context) {
	// Create the "cluster" collection
	err := s.mongoClient.Database(s.cfg.Mongo.Db).CreateCollection(ctx, s.cfg.Mongo.Collections.Cluster)
	if err != nil {
		if !utils.CheckErrMessages(err, serviceErrors.ErrMsgMongoCollectionAlreadyExists) {
			s.log.Warnf("(CreateCollection) err: {%v}", err)
		}
	}

	// Create the "folders" collection
	err = s.mongoClient.Database(s.cfg.Mongo.Db).CreateCollection(ctx, s.cfg.Mongo.Collections.Folder)
	if err != nil {
		if !utils.CheckErrMessages(err, serviceErrors.ErrMsgMongoCollectionAlreadyExists) {
			s.log.Warnf("(CreateCollection) err: {%v}", err)
		}
	}

	// Create indexes on the "cluster" collection
	{
		indexes, err := s.mongoClient.Database(s.cfg.Mongo.Db).Collection(s.cfg.Mongo.Collections.Cluster).Indexes().CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.D{
					{"cluster_name", "text"},
					{"title", "text"},
					{"note", "text"},
				},
				Options: options.Index().SetSparse(true).SetName(fmt.Sprintf("%s.text_index", s.cfg.Mongo.Collections.Cluster)),
			},
		})
		if err != nil && !utils.CheckErrMessages(err, serviceErrors.ErrMsgAlreadyExists) {
			s.log.Warnf("(CreateMany) err: {%v}", err)
		}
		s.log.Infof("(CreatedIndexes) indexes: {%v}", indexes)
	}

	// Create indexes on the "folders" collection
	{
		indexes, err := s.mongoClient.Database(s.cfg.Mongo.Db).Collection(s.cfg.Mongo.Collections.Folder).Indexes().CreateMany(ctx, []mongo.IndexModel{
			{
				Keys:    bson.D{{"folder_name", 1}},
				Options: options.Index().SetSparse(true).SetName(fmt.Sprintf("%s.%s_index", s.cfg.Mongo.Collections.Folder, "folder_name")),
			},
			{
				Keys:    bson.D{{"folder_thumbnail_key", 1}},
				Options: options.Index().SetUnique(true).SetName(fmt.Sprintf("%s.%s_unique_index", s.cfg.Mongo.Collections.Folder, "folder_thumbnail_key")),
			},
		})
		if err != nil && !utils.CheckErrMessages(err, serviceErrors.ErrMsgAlreadyExists) {
			s.log.Warnf("(CreateMany) err: {%v}", err)
		}
		s.log.Infof("(CreatedIndexes) indexes: {%v}", indexes)
	}

	list, err := s.mongoClient.Database(s.cfg.Mongo.Db).Collection(s.cfg.Mongo.Collections.Cluster).Indexes().List(ctx)
	if err != nil {
		s.log.Warnf("(initMongoDBCollections) [List] err: {%v}", err)
	}

	if list != nil {
		var results []bson.M
		if err := list.All(ctx, &results); err != nil {
			s.log.Warnf("(All) err: {%v}", err)
		}
		s.log.Infof("(indexes) results: {%#v}", results)
	}

	collections, err := s.mongoClient.Database(s.cfg.Mongo.Db).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		s.log.Warnf("(ListCollections) err: {%v}", err)
	}
	s.log.Infof("(Collections) created collections: {%v}", collections)
}

func (s *server) waitShootDown(duration time.Duration) {
	go func() {
		time.Sleep(duration)
		s.doneCh <- struct{}{}
	}()
}

func GetMicroserviceName(cfg *config.Config) string {
	return fmt.Sprintf("(%s)", strings.ToUpper(cfg.App.Name))
}

func GetMicroserviceVersion(cfg *config.Config) string {
	return fmt.Sprintf("(%s)", strings.ToUpper(cfg.App.Version))
}
