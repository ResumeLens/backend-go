package handler

import (
	// Use your module name
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
	// --- NEW: READ IDs FROM FORM ---
	// We get the IDs from the form-data body of the request.
	orgID := c.PostForm("organization_id")
	jobID := c.PostForm("job_id")
	candidateID := c.PostForm("candidate_id")

	// Basic validation to ensure IDs are provided.
	if orgID == "" || jobID == "" || candidateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization_id, job_id, and candidate_id are required"})
		return
	}
	// --- END NEW ---

	file, handler, err := c.Request.FormFile("myFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not retrieve file from request"})
		return
	}
	defer file.Close()

	// --- MODIFIED: PASS IDs TO SERVICE ---
	// The core logic is delegated to the service, now with the new IDs.
	if err := h.service.UploadFile(c.Request.Context(), file, handler, orgID, jobID, candidateID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File processed and stored successfully."})
}
