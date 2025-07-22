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
	err := DB.AutoMigrate(&models.User{}, &models.Organization{}, &models.Invite{}, &models.Candidate{}, &models.JobApplication{}, &models.Role{})
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	fmt.Println("Database migrated successfully.")
}

// CheckUserPermission checks if a user has a specific permission based on their role.
func CheckUserPermission(userID string, permission string) (bool, error) {
	// Fetch the user
	var user models.User
	if err := DB.First(&user, "id = ?", userID).Error; err != nil {
		return false, err
	}

	// Fetch the role
	var role models.Role
	if err := DB.First(&role, "id = ?", user.RoleID).Error; err != nil {
		return false, err
	}

	switch permission {
	case "home":
		return role.HomePermission, nil
	case "create_job":
		return role.CreateJobPermission, nil
	case "view_job":
		return role.ViewJobPermission, nil
	case "iam":
		return role.IamPermission, nil
	default:
		return false, nil
	}
}
