package services

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
)

type JobApplicationService struct {
	gcsClient  *storage.Client
	bucketName string
}

type Metadata struct {
	ResumeFilename      string    `json:"resume_filename,omitempty"`
	CoverLetterFilename string    `json:"cover_letter_filename,omitempty"`
	LastUpdated         time.Time `json:"last_updated"`
}

func NewJobApplicationService(gcsClient *storage.Client, bucketName string) *JobApplicationService {
	return &JobApplicationService{
		gcsClient:  gcsClient,
		bucketName: bucketName,
	}
}

func (s *JobApplicationService) UploadResume(ctx context.Context, file multipart.File, handler *multipart.FileHeader, orgID, jobID, candidateID string) error {
	if filepath.Ext(handler.Filename) == ".zip" {
		return s.processZip(ctx, file, handler, orgID, jobID, candidateID)
	}

	ext := filepath.Ext(handler.Filename)
	objectName := s.buildObjectPath(orgID, jobID, candidateID, "resume", ext)

	if err := s.uploadToGCS(ctx, objectName, file); err != nil {
		return fmt.Errorf("failed to upload resume: %w", err)
	}

	metadataUpdate := map[string]interface{}{"resume_filename": handler.Filename}
	return s.updateMetadata(ctx, orgID, jobID, candidateID, metadataUpdate)
}

func (s *JobApplicationService) UploadCoverLetter(ctx context.Context, file multipart.File, handler *multipart.FileHeader, orgID, jobID, candidateID string) error {
	if filepath.Ext(handler.Filename) == ".zip" {
		return fmt.Errorf("zip files are not supported for cover letters")
	}

	ext := filepath.Ext(handler.Filename)
	objectName := s.buildObjectPath(orgID, jobID, candidateID, "cover_letter", ext)

	if err := s.uploadToGCS(ctx, objectName, file); err != nil {
		return fmt.Errorf("failed to upload cover letter: %w", err)
	}

	metadataUpdate := map[string]interface{}{"cover_letter_filename": handler.Filename}
	return s.updateMetadata(ctx, orgID, jobID, candidateID, metadataUpdate)
}

func (s *JobApplicationService) buildObjectPath(orgID, jobID, candidateID, fileType, ext string) string {
	return fmt.Sprintf("org-%s/job-%s/candidate-%s/%s%s", orgID, jobID, candidateID, fileType, ext)
}

func (s *JobApplicationService) updateMetadata(ctx context.Context, orgID, jobID, candidateID string, updates map[string]interface{}) error {
	metadataPath := fmt.Sprintf("org-%s/job-%s/candidate-%s/metadata.json", orgID, jobID, candidateID)
	obj := s.gcsClient.Bucket(s.bucketName).Object(metadataPath)

	reader, err := obj.NewReader(ctx)
	currentMetadata := make(map[string]interface{})

	if err == nil {
		defer reader.Close()
		data, readErr := ioutil.ReadAll(reader)
		if readErr != nil {
			return fmt.Errorf("failed to read existing metadata: %w", readErr)
		}
		if unmarshalErr := json.Unmarshal(data, &currentMetadata); unmarshalErr != nil {
			return fmt.Errorf("failed to parse existing metadata: %w", unmarshalErr)
		}
	} else if err != storage.ErrObjectNotExist {
		return fmt.Errorf("failed to download metadata: %w", err)
	}

	for key, value := range updates {
		currentMetadata[key] = value
	}
	currentMetadata["last_updated"] = time.Now().UTC().Format(time.RFC3339)

	updatedData, marshalErr := json.MarshalIndent(currentMetadata, "", "  ")
	if marshalErr != nil {
		return fmt.Errorf("failed to marshal updated metadata: %w", marshalErr)
	}

	return s.uploadToGCS(ctx, metadataPath, bytes.NewReader(updatedData))
}

func (s *JobApplicationService) uploadToGCS(ctx context.Context, objectName string, reader io.Reader) error {
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

func (s *JobApplicationService) processZip(ctx context.Context, zipFile multipart.File, zipHandler *multipart.FileHeader, orgID, jobID, candidateID string) error {
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

		objectName := s.buildObjectPath(orgID, jobID, candidateID, "resume", filepath.Ext(file.Name))
		if err := s.uploadToGCS(ctx, objectName, zippedFile); err != nil {
			log.Printf("Failed to upload %s from zip: %v", file.Name, err)
		}
	}

	return nil
}
