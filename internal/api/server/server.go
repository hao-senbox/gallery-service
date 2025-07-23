package server

import (
	"context"
	"gallery-service/config"
	"gallery-service/internal/api/rest/middlewares"
	"gallery-service/pkg/consul"
	httpPkg "gallery-service/pkg/http"
	"gallery-service/pkg/mongodb"
	"gallery-service/pkg/zap"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server interface {
	StartServer() error
}

type server struct {
	fiber        *fiber.App
	log          zap.Logger
	cfg          *config.Config
	mw           middlewares.MiddlewareManager
	mongoClient  *mongo.Client
	consulClient *api.Client
	doneCh       chan struct{}
}

const (
	maxHeaderBytes = 1 << 20         // 1 MB
	stackSize      = 1 << 10         // 1 KB
	bodyLimit      = 2 * 1024 * 1024 // 2 MB
	readTimeout    = 15 * time.Second
	writeTimeout   = 15 * time.Second
	gzipLevel      = 5
)

func New(log zap.Logger, cfg *config.Config) (*server, error) {
	server := &server{
		fiber: fiber.New(fiber.Config{
			ServerHeader: cfg.App.Name,
			AppName:      cfg.App.Name,
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return httpPkg.NewBadRequestError(c, err.Error(), cfg.App.API.Rest.Setting.DebugErrorsResponse)
			},
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			BodyLimit:    bodyLimit,
		}),
		log:    log,
		cfg:    cfg,
		doneCh: make(chan struct{}),
	}

	return server, nil
}

func (s *server) StartServer() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	s.mw = middlewares.NewMiddlewareManager(s.log, s.cfg)

	mongoDBClient := mongodb.NewMongoDBConn(ctx, s.log, s.cfg.Mongo)
	s.mongoClient = mongoDBClient.GetClient()
	defer mongoDBClient.Close()

	s.mongoMigrationUp(ctx)

	consulConn := consul.NewConsulConn(s.log, s.cfg)
	s.consulClient = consulConn.Connect()
	defer consulConn.Deregister()

	//if err := consumers.NewEngine(s.cfg, s.log, s.mongoClient).Start(ctx); err != nil {
	//	s.log.DPanicf("Failed to start consumers: {%v}", err)
	//	return err
	//}

	errorChan := make(chan error, 1)
	go s.start(errorChan)

	select {
	case <-ctx.Done():
		s.log.Infof("%s shutting down the server", GetMicroserviceName(s.cfg))
		s.waitShootDown(waitShotDownDuration)
	case err := <-errorChan:
		s.log.Fatalf("Failed to start HTTP server: {%v}", err)
		cancel()
		return err
	}

	<-s.doneCh
	s.log.Infof("%s server exited properly", GetMicroserviceName(s.cfg))

	return nil
}
