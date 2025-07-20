package db

import (
	"fmt"
	"log"

	"github.com/resumelens/authservice/internal/config"
	"github.com/resumelens/authservice/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(cfg *config.Config) {
	dsn := cfg.DatabaseURL
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = database

	fmt.Println("Database connected.")

	migrateDatabase()
}

func migrateDatabase() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Organization{},
		&models.Invite{},
		&models.Candidate{},
		&models.JobApplication{},
		&models.Job{},
		&models.JobAnalytics{},
	)
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	fmt.Println("Database migrated successfully.")
}
