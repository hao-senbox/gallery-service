package config

import (
	"fmt"
	"gallery-service/pkg/constants"
	"gallery-service/pkg/kafka"
	"gallery-service/pkg/mongodb"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
)

// AppConfiguration holds the application-specific configuration
type AppConfiguration struct {
	Name        string    `mapstructure:"name"`
	Version     string    `mapstructure:"version"`
	Environment string    `mapstructure:"environment"`
	API         APIConfig `mapstructure:"api"`
}

// APIConfig holds API-related configurations
type APIConfig struct {
	Rest RestConfig `mapstructure:"rest"`
}

// RestConfig holds REST API server configurations
type RestConfig struct {
	Host    string        `mapstructure:"host"`
	Port    string        `mapstructure:"port"`
	Setting SettingConfig `mapstructure:"setting"`
}

// SettingConfig holds settings for the REST API
type SettingConfig struct {
	Debug               bool     `mapstructure:"debug"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
}

type Consul struct {
	Host string `mapstructure:"host" validate:"required"`
	Port string `mapstructure:"port" validate:"required"`
}

type Registry struct {
	Host string `mapstructure:"host" validate:"required"`
}

// Config is the overall configuration structure
type Config struct {
	App      AppConfiguration `mapstructure:"app"`
	Mongo    *mongodb.Config  `mapstructure:"mongo"`
	Consul   Consul           `mapstructure:"consul" validate:"required"`
	Registry Registry         `mapstructure:"registry" validate:"required"`
	Kafka    kafka.Config     `mapstructure:"kafka" validate:"required"`
}

// LoadConfig reads the configuration from a file
func LoadConfig(cfg *viper.Viper, configPath string) (*Config, error) {
	if configPath != "" {
		// Use config file from the flag.
		cfg.SetConfigFile(configPath)
	} else {
		// Search config in home directory with name ".config.dev" (without extension).
		cfg.AddConfigPath(".")
		cfg.SetConfigType(constants.Yaml)
		cfg.SetConfigName(".config.dev")
	}

	cfg.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := cfg.ReadInConfig(); err == nil {

		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", cfg.ConfigFileUsed())

		fmt.Println("--- Configuration read from file ---")
		for s, i := range cfg.AllSettings() {
			fmt.Printf("\t%s = %s\n", s, i)
		}
		fmt.Println("---")
	}

	config := &Config{}
	if err := cfg.Unmarshal(&config); err != nil {
		return nil, errors.Errorf("\n(LoadConfig) Unable to decode into struct, err: {%v}\n", err)
	}

	return config, nil
}
