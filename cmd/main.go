package main

import (
	"log"

	"github.com/resumelens/authservice/internal/config"
	"github.com/resumelens/authservice/internal/db"
	"github.com/resumelens/authservice/internal/gcs"
	"github.com/resumelens/authservice/internal/handler"
	"github.com/resumelens/authservice/internal/routes"
	uploader "github.com/resumelens/authservice/resume-uploader"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// 2. Initialize Dependencies
	db.ConnectDatabase()
	gcs.InitClient()

	// 3. Initialize Services
	resumeService := uploader.NewService(gcs.GCSClient, cfg.GCSBucketName)

	// 4. Initialize Handlers
	resumeHandler := handler.NewResumeHandler(resumeService)

	// 5. Setup Router
	r := routes.SetupRouter(resumeHandler)

	// 6. Start Server
	port := cfg.Port
	if port == "" {
		port = "8080" // Just the number
	}

	log.Printf("Server running on port %s", port)

	// --- THIS IS THE FIX ---
	// We add the colon here to create the correct address format ":8080"
	r.Run(":" + port)
}
