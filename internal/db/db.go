package db

import (
	"fmt"
	"log"

	"github.com/resumelens/authservice/internal/models"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := viper.GetString("DB_URL")
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = database

	fmt.Println("Database connected.")

	migrateDatabase()
}

func migrateDatabase() {
	err := DB.AutoMigrate(&models.User{}, &models.Organization{}, &models.Role{}, &models.Invite{})
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	fmt.Println("Database migrated successfully.")
}
