package uploader

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"

	"cloud.google.com/go/storage"
)

// Service defines the interface for the resume upload feature.
type Service struct {
	gcsClient  *storage.Client
	bucketName string
}

// NewService creates a new instance of the upload service.
func NewService(gcsClient *storage.Client, bucketName string) *Service {
	return &Service{
		gcsClient:  gcsClient,
		bucketName: bucketName,
	}
}

// --- MODIFIED FUNCTION SIGNATURE ---
// UploadFile now accepts orgID, jobID, and candidateID.
func (s *Service) UploadFile(ctx context.Context, file multipart.File, handler *multipart.FileHeader, orgID, jobID, candidateID string) error {
	// For now, we assume single file uploads for path creation.
	// ZIP file processing logic will need to be adapted if candidates in a zip belong to different jobs.
	if filepath.Ext(handler.Filename) == ".zip" {
		// We pass the IDs to the zip processor.
		return s.processZip(ctx, file, handler, orgID, jobID, candidateID)
	}

	// --- NEW GCS PATH LOGIC ---
	// Get the file extension (e.g., ".pdf").
	ext := filepath.Ext(handler.Filename)
	// Construct the object path according to the plan.
	objectName := fmt.Sprintf("organization-%s/job-%s/candidate-%s/resume%s", orgID, jobID, candidateID, ext)

	return s.uploadToGCS(ctx, objectName, file)
}

func (s *Service) uploadToGCS(ctx context.Context, objectName string, reader io.Reader) error {
	obj := s.gcsClient.Bucket(s.bucketName).Object(objectName)
	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, reader); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}
	log.Printf("Successfully uploaded %s", objectName)
	return nil
}

// processZip is also updated to accept the IDs.
func (s *Service) processZip(ctx context.Context, zipFile multipart.File, zipHandler *multipart.FileHeader, orgID, jobID, candidateID string) error {
	// This is a simplified implementation. A real-world scenario might need to handle
	// multiple candidates within a single zip, each with their own ID.
	reader, err := zip.NewReader(zipFile, zipHandler.Size)
	if err != nil {
		return fmt.Errorf("could not read zip file: %w", err)
	}
	for _, file := range reader.File {
		if ext := filepath.Ext(file.Name); ext != ".pdf" && ext != ".docx" {
			continue
		}
		zippedFile, err := file.Open()
		if err != nil {
			log.Printf("Error opening %s in zip: %v", file.Name, err)
			continue
		}
		defer zippedFile.Close()

		// --- NEW GCS PATH LOGIC FOR ZIPPED FILES ---
		// We use a placeholder for candidateID from the filename, but use the provided org/job IDs.
		// A more robust solution would parse candidate info from the filename or an associated manifest.
		objectName := fmt.Sprintf("organization-%s/job-%s/candidate-%s/resume%s", orgID, jobID, file.Name, filepath.Ext(file.Name))
		if err := s.uploadToGCS(ctx, objectName, zippedFile); err != nil {
			log.Printf("Failed to upload %s from zip: %v", file.Name, err)
		}
	}
	return nil
}
