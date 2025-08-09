package mongodb

import (
	"context"
	"fmt"
	"gallery-service/pkg/zap"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	URI               string           `mapstructure:"uri"`
	Username          string           `mapstructure:"username"`
	Password          string           `mapstructure:"password"`
	Db                string           `mapstructure:"db"`
	ConnectionTimeout time.Duration    `mapstructure:"connection_timeout"`
	MaxConnIdleTime   time.Duration    `mapstructure:"max_conn_idle_time"`
	MinPoolSize       int              `mapstructure:"min_pool_size"`
	MaxPoolSize       int              `mapstructure:"max_pool_size"`
	Collections       MongoCollections `mapstructure:"collections"`
}

type MongoCollections struct {
	Cluster string `mapstructure:"cluster" validate:"required"`
	Folder  string `mapstructure:"folder" validate:"required"`
	Topic   string `mapstructure:"topic" validate:"required"`
}

// Client represents a service that interacts with MongoDB.
type Client interface {
	Health() map[string]string
	Close()
	GetClient() *mongo.Client
}

type service struct {
	client *mongo.Client
	log    zap.Logger
	cfg    *Config
}

// DBConfig returns the MongoDB client configuration.
func DBConfig(cfg *Config) *options.ClientOptions {
	// Fetching config values
	username := cfg.Username
	password := cfg.Password

	// connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s&authMechanism=%s", username, password, "127.0.0.1", 17017, cfg.Db, "admin", "SCRAM-SHA-1")

	// Create client options
	clientOptions := options.Client().ApplyURI(cfg.URI).
		SetAuth(options.Credential{
			Username:      username,
			Password:      password,
			AuthSource:    "admin",
			AuthMechanism: "SCRAM-SHA-1",
		}).
		SetConnectTimeout(cfg.ConnectionTimeout * time.Second).
		SetMaxConnIdleTime(cfg.MaxConnIdleTime * time.Minute).
		SetMinPoolSize(uint64(cfg.MinPoolSize)).
		SetMaxPoolSize(uint64(cfg.MaxPoolSize))

	// Additional configuration options can go here

	return clientOptions
}

// NewMongoDBConn initializes a new MongoDB client or returns the existing one.
func NewMongoDBConn(ctx context.Context, log zap.Logger, cfg *Config) *service {
	// Create MongoDB client and connect
	client, err := mongo.Connect(ctx, DBConfig(cfg))
	if err != nil {
		log.Fatalf("(mongodb.New) Failed to create MongoDB and connect client, error: {%v}", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("(mongodb.New) Failed to ping MongoDB, error: {%v}", err)
	}

	log.Info("Connected to MongoDB!")

	return &service{
		client: client, // Store the client instance
		log:    log,
		cfg:    cfg,
	}
}

// GetClient returns the underlying *mongo.Client instance.
func (s *service) GetClient() *mongo.Client {
	return s.client
}

// Health checks the health of the MongoDB connection.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the MongoDB server
	err := s.client.Ping(ctx, nil)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		s.log.Errorf("(mongodb.Health) DATABASE DOWN: {%v}", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Gather MongoDB server stats (optional)
	var serverStatus bson.M
	err = s.client.Database(s.cfg.Db).RunCommand(ctx, bson.D{{"serverStatus", 1}}).Decode(&serverStatus)
	if err != nil {
		stats["message"] = "Failed to get server stats"
	} else {
		stats["message"] = "MongoDB server stats fetched"
	}

	// MongoDB health metrics can include:
	// - Connections
	// - Operations count
	// - Memory usage, etc.
	stats["server_status"] = fmt.Sprintf("%v", serverStatus)

	return stats
}

// Close closes the MongoDB client connection.
func (s *service) Close() {
	if s.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := s.client.Disconnect(ctx)
		if err != nil {
			s.log.Errorf("(mongodb.Close) Failed to disconnect MongoDB: {%v}", err)
		} else {
			s.log.Info("Disconnected from MongoDB.")
		}
	}
}
