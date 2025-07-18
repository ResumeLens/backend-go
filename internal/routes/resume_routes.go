package routes

import (
	"github.com/resumelens/authservice/internal/handler" // Use your module name from go.mod

	"github.com/gin-gonic/gin"
)

// RegisterResumeRoutes sets up the routes for the resume upload feature.
func RegisterResumeRoutes(router *gin.RouterGroup, h *handler.ResumeHandler) {
	router.POST("/upload-resume", h.Upload)
}
