package main

import (
	"archive/zip" //importing the requirements to handle zip files
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func processZip(zipFile multipart.File, zipHandler *multipart.FileHeader) error {
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

		log.Printf("Extracting file from zip: %s", file.Name)

		zippedFile, err := file.Open()
		if err != nil {
			log.Printf("Error opening file %s in zip: %v", file.Name, err)
			continue
		}
		defer zippedFile.Close()

		targetPath := filepath.Join("uploads", file.Name)
		dst, err := os.Create(targetPath)
		if err != nil {
			log.Printf("Error creating destination file for %s: %v", file.Name, err)
			continue
		}
		defer dst.Close()

		if _, err := io.Copy(dst, zippedFile); err != nil {
			log.Printf("Error copying content for %s: %v", file.Name, err)
		}
	}
	return nil
}

func saveFile(file multipart.File, handle *multipart.FileHeader) error {

	dst, err := os.Create(filepath.Join("uploads", handle.Filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return err
	}

	return nil
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("ZIP file detected: %s. Processing...", handler.Filename)
		if err := processZip(file, handler); err != nil {
			http.Error(w, "Failed to process ZIP file", http.StatusInternalServerError)
			log.Printf("Error processing zip: %v", err)
			return
		}
		fmt.Fprintf(w, "Successfully processed ZIP file: %s\n", handler.Filename)
	} else {

		log.Printf("Single file detected: %s. Saving...", handler.Filename)
		if err := saveFile(file, handler); err != nil {
			http.Error(w, "Error saving the file", http.StatusInternalServerError)
			log.Println("Error saving file:", err)
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
