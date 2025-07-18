package gcs

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
)

// GCSClient holds the global client instance for Google Cloud Storage.
var GCSClient *storage.Client

// InitClient creates the connection to Google Cloud Storage.
func InitClient() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}
	GCSClient = client
	log.Println("Google Cloud Storage client initialized.")
}
