package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	uploader "github.com/resumelens/authservice/resume-uploader"
)

// ResumeHandler handles HTTP requests for the resume uploader feature.
type ResumeHandler struct {
	service *uploader.Service
}

// NewResumeHandler creates a new handler for the resume feature.
func NewResumeHandler(s *uploader.Service) *ResumeHandler {
	return &ResumeHandler{service: s}
}

// Upload handles the resume file upload.
func (h *ResumeHandler) Upload(c *gin.Context) {
	orgID := c.PostForm("organization_id")
	jobID := c.PostForm("job_id")
	candidateID := c.PostForm("candidate_id")

	if orgID == "" || jobID == "" || candidateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization_id, job_id, and candidate_id are required"})
		return
	}

	file, handler, err := c.Request.FormFile("myFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not retrieve file from request"})
		return
	}
	defer file.Close()

	if err := h.service.UploadFile(c.Request.Context(), file, handler, orgID, jobID, candidateID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process resume file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resume processed and stored successfully."})
}

// --- NEW HANDLER for Cover Letter ---
// UploadCoverLetter handles the cover letter file upload.
func (h *ResumeHandler) UploadCoverLetter(c *gin.Context) {
	orgID := c.PostForm("organization_id")
	jobID := c.PostForm("job_id")
	candidateID := c.PostForm("candidate_id")

	if orgID == "" || jobID == "" || candidateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization_id, job_id, and candidate_id are required"})
		return
	}

	// We look for a different form field name: "coverLetterFile"
	file, handler, err := c.Request.FormFile("coverLetterFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not retrieve cover letter file from request"})
		return
	}
	defer file.Close()

	// We call the new service method
	if err := h.service.UploadCoverLetter(c.Request.Context(), file, handler, orgID, jobID, candidateID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process cover letter file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cover letter stored successfully."})
}

// --- NEW HANDLER for Metadata ---
// UploadMetadata handles the metadata JSON upload.
