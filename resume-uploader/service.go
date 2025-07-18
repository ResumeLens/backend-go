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

// UploadFile manages the core logic for uploading a file or processing a ZIP.
func (s *Service) UploadFile(ctx context.Context, file multipart.File, handler *multipart.FileHeader) error {
	if filepath.Ext(handler.Filename) == ".zip" {
		return s.processZip(ctx, file, handler)
	}
	return s.uploadToGCS(ctx, handler.Filename, file)
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

func (s *Service) processZip(ctx context.Context, zipFile multipart.File, zipHandler *multipart.FileHeader) error {
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
		if err := s.uploadToGCS(ctx, file.Name, zippedFile); err != nil {
			log.Printf("Failed to upload %s from zip: %v", file.Name, err)
		}
	}
	return nil
}
