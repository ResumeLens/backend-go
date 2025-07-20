package gcs

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var GCSClient *storage.Client

func InitClient(projectID, credentialsFile string) {
	ctx := context.Background()

	var opts []option.ClientOption

	if credentialsFile != "" {
		if _, err := os.Stat(credentialsFile); err == nil {
			opts = append(opts, option.WithCredentialsFile(credentialsFile))
			log.Printf("Using GCS credentials from file: %s", credentialsFile)
		} else {
			log.Printf("Warning: Credentials file not found: %s", credentialsFile)
		}
	}

	if projectID != "" {
		log.Printf("Using GCS project ID: %s", projectID)
	}

	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}
	GCSClient = client
	log.Println("Google Cloud Storage client initialized.")
}
