package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/resumelens/authservice/internal/handler"
	"github.com/resumelens/authservice/internal/middleware"
)

func SetupRouter(jobApplicationHandler *handler.JobApplicationHandler, authHandler *handler.AuthHandler) *gin.Engine {
	router := gin.Default()

	router.MaxMultipartMemory = 30 << 20

	api := router.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Server is running",
				"service": "Job Application Backend",
			})
		})

		api.POST("/signup", authHandler.Signup)
		api.POST("/login", authHandler.Login)
		api.GET("/validate-invite", authHandler.ValidateInvite)
		api.POST("/accept-invite", authHandler.AcceptInvite)
		api.POST("/refresh-token", authHandler.RefreshToken)

		secured := api.Group("/")
		secured.Use(middleware.JWTAuthMiddleware())
		{
			secured.POST("/invite", authHandler.Invite)
			secured.POST("/upload-resume", jobApplicationHandler.UploadResume)
			secured.POST("/upload-cover-letter", jobApplicationHandler.UploadCoverLetter)
		}
	}

	return router
}
