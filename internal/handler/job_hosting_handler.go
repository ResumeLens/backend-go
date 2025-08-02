package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/resumelens/authservice/internal/services"
)

type JobHostingHandler struct {
	jobHostingService *services.JobHostingService
}

func NewJobHostingHandler(jobHostingService *services.JobHostingService) *JobHostingHandler {
	return &JobHostingHandler{jobHostingService: jobHostingService}
}

func (h *JobHostingHandler) CreateJob(c *gin.Context) {
	var req services.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	orgID, exists := c.Get("organizationID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in token"})
		return
	}

	response, statusCode := h.jobHostingService.CreateJob(req, orgID.(string), userID.(string))
	c.JSON(statusCode, response)
}

func (h *JobHostingHandler) GetJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID is required"})
		return
	}

	response, statusCode := h.jobHostingService.GetJob(id)
	c.JSON(statusCode, response)
}

func (h *JobHostingHandler) GetJobsForOrganization(c *gin.Context) {
	orgID, _ := c.Get("organizationID")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	response, status := h.jobHostingService.GetJobsForOrganization(orgID.(string), page, pageSize)
	c.JSON(status, response)
}

func (h *JobHostingHandler) UpdateJob(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID is required"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields provided for update"})
		return
	}

	response, status := h.jobHostingService.UpdateJob(jobID, updates)
	c.JSON(status, response)
}

func (h *JobHostingHandler) DeleteJob(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID is required"})
		return
	}

	response, status := h.jobHostingService.DeleteJob(jobID)
	c.JSON(status, response)
}

func (h *JobHostingHandler) GetJobApplications(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	response, status := h.jobHostingService.GetJobApplications(jobID, page, pageSize)
	c.JSON(status, response)
}
