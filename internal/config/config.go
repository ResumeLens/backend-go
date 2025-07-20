package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port                  string `mapstructure:"PORT"`
	GCSBucketName         string `mapstructure:"GCS_BUCKET_NAME"`
	GoogleProjectID       string `mapstructure:"GOOGLE_PROJECT_ID"`
	GoogleCredentialsFile string `mapstructure:"GOOGLE_APPLICATION_CREDENTIALS"`

	DatabaseURL string `mapstructure:"DB_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
	JWTExpiry   int    `mapstructure:"JWT_EXPIRY"`

	SMTPHost       string `mapstructure:"SMTP_HOST"`
	SMTPPort       string `mapstructure:"SMTP_PORT"`
	SMTPUser       string `mapstructure:"SMTP_USER"`
	SMTPPass       string `mapstructure:"SMTP_PASS"`
	SMTPSenderName string `mapstructure:"SMTP_SENDER_NAME"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error loading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Validate required fields
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Set defaults
	if config.Port == "" {
		config.Port = "8080"
	}
	if config.JWTExpiry == 0 {
		config.JWTExpiry = 60
	}

	return &config, nil
}

func validateConfig(cfg *Config) error {
	required := map[string]string{
		"GCS_BUCKET_NAME":                cfg.GCSBucketName,
		"GOOGLE_PROJECT_ID":              cfg.GoogleProjectID,
		"GOOGLE_APPLICATION_CREDENTIALS": cfg.GoogleCredentialsFile,
		"DB_URL":                         cfg.DatabaseURL,
		"JWT_SECRET":                     cfg.JWTSecret,
	}

	for name, value := range required {
		if value == "" {
			return fmt.Errorf("%s is required", name)
		}
	}

	return nil
}
