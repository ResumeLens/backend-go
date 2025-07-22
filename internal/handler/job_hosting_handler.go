package handler

import (
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, statusCode := h.jobHostingService.CreateJob(req)
	c.JSON(statusCode, response)
}

func (h *JobHostingHandler) GetJob(c *gin.Context) {
	id := c.Param("id")
	response, statusCode := h.jobHostingService.GetJob(id)
	c.JSON(statusCode, response)
}
