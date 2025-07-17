package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Automatically override with environment variables if available
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error loading config file: %w", err)
	}

	return nil
}
