package main

import (
	"log"

	"github.com/resumelens/authservice/internal/config"
	"github.com/resumelens/authservice/internal/db"
	"github.com/resumelens/authservice/internal/gcs"
	"github.com/resumelens/authservice/internal/handler"
	"github.com/resumelens/authservice/internal/routes"
	"github.com/resumelens/authservice/internal/services"
	"github.com/resumelens/authservice/internal/utils"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	db.ConnectDatabase(cfg)
	utils.InitJWT(cfg)
	gcs.InitClient(cfg.GoogleProjectID, cfg.GoogleCredentialsFile)

	// Services
	jobApplicationService := services.NewJobApplicationService(gcs.GCSClient, cfg.GCSBucketName)
	authService := services.NewAuthService(cfg)
	jobHostingService := services.NewJobHostingService(cfg)

	// Handlers
	jobApplicationHandler := handler.NewJobApplicationHandler(jobApplicationService)
	authHandler := handler.NewAuthHandler(authService)
	jobHostingHandler := handler.NewJobHostingHandler(jobHostingService)

	// Routes
	r := routes.SetupRouter(jobApplicationHandler, authHandler, jobHostingHandler)

	port := cfg.Port
	if port == "" {
		port = "8000"
	}

	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
