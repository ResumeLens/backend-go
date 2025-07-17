//create a service account for gcs neural hive 
// get a service key 
// download as json and save it locally 
//before running the main.go in the terminal type $env:GOOGLE_APPLICATION_CREDENTIALS="/location path/credentials.json" (for microsoft for apple go fine it ) considering you named the json file as credentials.json
//also go to the bucket then permissions find ur service account and give permission -> 1.Storage Object Admin 2.Storage Object Creator
// now run it 

package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"cloud.google.com/go/storage"
)

var bucketName = "resume-files-bucket"

func uploadToGCS(ctx context.Context, client *storage.Client, reader io.Reader, objectName string) error {
	obj := client.Bucket(bucketName).Object(objectName)
	wc := obj.NewWriter(ctx)

	// We copy from the reader to the GCS writer.
	if _, err := io.Copy(wc, reader); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}

	log.Printf("Successfully uploaded %s to GCS bucket %s", objectName, bucketName)
	return nil
}

func processZip(ctx context.Context, client *storage.Client, zipFile multipart.File, zipHandler *multipart.FileHeader) error {
	reader, err := zip.NewReader(zipFile, zipHandler.Size)
	if err != nil {
		return fmt.Errorf("could not read zip file: %w", err)
	}

	for _, file := range reader.File {
		ext := filepath.Ext(file.Name)
		if ext != ".pdf" && ext != ".docx" && ext != ".txt" {
			log.Printf("Skipping file in zip: %s (unsupported type)", file.Name)
			continue
		}

		zippedFile, err := file.Open()
		if err != nil {
			log.Printf("Error opening file %s in zip: %v", file.Name, err)
			continue
		}
		defer zippedFile.Close()

		if err := uploadToGCS(ctx, client, zippedFile, file.Name); err != nil {
			log.Printf("Failed to upload %s to GCS: %v", file.Name, err)
		}
	}
	return nil
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		http.Error(w, "Failed to create GCS client", http.StatusInternalServerError)
		log.Printf("storage.NewClient: %v", err)
		return
	}
	defer client.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 30<<20)

	if err := r.ParseMultipartForm(30 << 20); err != nil {
		http.Error(w, "The uploaded file is too big.", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if filepath.Ext(handler.Filename) == ".zip" {
		if err := processZip(ctx, client, file, handler); err != nil {
			http.Error(w, "Failed to process ZIP file", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Successfully processed ZIP file: %s\n", handler.Filename)
	} else {

		if err := uploadToGCS(ctx, client, file, handler.Filename); err != nil {
			http.Error(w, "Failed to upload to GCS", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Successfully Uploaded and Saved File: %s\n", handler.Filename)
	}
}

func main() {
	http.HandleFunc("/upload", uploadHandler)

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
