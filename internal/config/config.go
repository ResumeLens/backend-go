package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application.
// The mapstructure tags tell Viper how to map settings from your .env file.
type Config struct {
	Port          string `mapstructure:"PORT"`
	GCSBucketName string `mapstructure:"GCS_BUCKET_NAME"`
	// Your teammates' other config fields would also go here.
}

// LoadConfig now returns a pointer to the Config struct and an error.
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error loading config file: %w", err)
	}

	// This new part decodes the loaded settings into our struct.
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &config, nil
}
