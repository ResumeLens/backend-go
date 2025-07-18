package handler

import (
	// Use your module name from go.mod
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

// Upload is the specific method for handling file uploads in Gin.
func (h *ResumeHandler) Upload(c *gin.Context) {
	file, handler, err := c.Request.FormFile("myFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not retrieve file from request"})
		return
	}
	defer file.Close()

	// The core logic is delegated to the service
	if err := h.service.UploadFile(c.Request.Context(), file, handler); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File processed successfully."})
}
