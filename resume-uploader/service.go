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

// UploadFile handles uploading the main resume file.
func (s *Service) UploadFile(ctx context.Context, file multipart.File, handler *multipart.FileHeader, orgID, jobID, candidateID string) error {
	if filepath.Ext(handler.Filename) == ".zip" {
		return s.processZip(ctx, file, handler, orgID, jobID, candidateID)
	}

	ext := filepath.Ext(handler.Filename)
	objectName := fmt.Sprintf("org-%s/job-%s/candidate-%s/resume%s", orgID, jobID, candidateID, ext)

	return s.uploadToGCS(ctx, objectName, file)
}

// --- NEW FUNCTION for Cover Letter ---
// UploadCoverLetter handles uploading a cover letter.
func (s *Service) UploadCoverLetter(ctx context.Context, file multipart.File, handler *multipart.FileHeader, orgID, jobID, candidateID string) error {
	// We don't support zip for cover letters to keep it simple.
	if filepath.Ext(handler.Filename) == ".zip" {
		return fmt.Errorf("zip files are not supported for cover letters")
	}

	ext := filepath.Ext(handler.Filename)
	// The object name is now hardcoded to 'cover_letter'.
	objectName := fmt.Sprintf("org-%s/job-%s/candidate-%s/cover_letter%s", orgID, jobID, candidateID, ext)

	return s.uploadToGCS(ctx, objectName, file)
}

// --- NEW FUNCTION for Metadata ---
// UploadMetadata handles uploading a metadata JSON file.
func (s *Service) UploadMetadata(ctx context.Context, jsonData io.Reader, orgID, jobID, candidateID string) error {
	// The object name is hardcoded to 'metadata.json'.
	objectName := fmt.Sprintf("org-%s/job-%s/candidate-%s/metadata.json", orgID, jobID, candidateID)
	return s.uploadToGCS(ctx, objectName, jsonData)
}

// uploadToGCS is a private helper function for GCS uploads.
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

// processZip handles zip files for resumes.
func (s *Service) processZip(ctx context.Context, zipFile multipart.File, zipHandler *multipart.FileHeader, orgID, jobID, candidateID string) error {
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

		objectName := fmt.Sprintf("org-%s/job-%s/candidate-%s/resume%s", orgID, jobID, file.Name, filepath.Ext(file.Name))
		if err := s.uploadToGCS(ctx, objectName, zippedFile); err != nil {
			log.Printf("Failed to upload %s from zip: %v", file.Name, err)
		}
	}
	return nil
}
