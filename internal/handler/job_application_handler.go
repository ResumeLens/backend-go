package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/resumelens/authservice/internal/services"
)

type JobApplicationHandler struct {
	service *services.JobApplicationService
}

func NewJobApplicationHandler(s *services.JobApplicationService) *JobApplicationHandler {
	return &JobApplicationHandler{service: s}
}

func (h *JobApplicationHandler) UploadResume(c *gin.Context) {
	orgID := c.PostForm("organization_id")
	jobID := c.PostForm("job_id")
	candidateID := c.PostForm("candidate_id")

	if orgID == "" || jobID == "" || candidateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization_id, job_id, and candidate_id are required"})
		return
	}

	file, handler, err := c.Request.FormFile("resumeFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not retrieve file from request"})
		return
	}
	defer file.Close()

	if err := h.service.UploadResume(c.Request.Context(), file, handler, orgID, jobID, candidateID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process resume file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resume processed and stored successfully."})
}

func (h *JobApplicationHandler) UploadCoverLetter(c *gin.Context) {
	orgID := c.PostForm("organization_id")
	jobID := c.PostForm("job_id")
	candidateID := c.PostForm("candidate_id")

	if orgID == "" || jobID == "" || candidateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization_id, job_id, and candidate_id are required"})
		return
	}

	file, handler, err := c.Request.FormFile("coverLetterFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not retrieve cover letter file from request"})
		return
	}
	defer file.Close()

	if err := h.service.UploadCoverLetter(c.Request.Context(), file, handler, orgID, jobID, candidateID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process cover letter file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cover letter stored successfully."})
}
