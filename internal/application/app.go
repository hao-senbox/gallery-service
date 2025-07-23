package application

import (
	"fmt"
	"gallery-service/config"
	"gallery-service/internal/api/server"
	"gallery-service/pkg/zap"
	"github.com/spf13/viper"
)

// App holds the application configuration and services
type App struct {
	logger zap.Logger
	server server.Server
}

// New initializes a new application instance
func New(configPath string) (*App, error) {
	// Load configuration
	cfg := viper.New()
	c, err := config.LoadConfig(cfg, configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize logger
	logger, err := zap.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}
	logger.WithName(server.GetMicroserviceName(c))

	// Initialize the server
	s, err := server.New(logger, c)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize server: %w", err)
	}

	return &App{
		logger: logger,
		server: s,
	}, nil
}

// Run starts the application
func (a *App) Run() error {
	return a.server.StartServer()
}
