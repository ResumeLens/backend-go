package main

import (
	"log"

	"github.com/resumelens/authservice/internal/config"
	"github.com/resumelens/authservice/internal/db"
	"github.com/resumelens/authservice/internal/routes"
	"github.com/spf13/viper"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Config error: %s", err)
	}

	db.ConnectDatabase()

	r := routes.SetupRouter()

	port := viper.GetString("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
