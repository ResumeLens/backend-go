package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/resumelens/authservice/internal/handler"
	"github.com/resumelens/authservice/internal/middleware"
)

// The function now accepts 'resumeHandler' as a parameter.
func SetupRouter(resumeHandler *handler.ResumeHandler) *gin.Engine {
	router := gin.Default()

	// This sets a memory limit for file uploads.
	router.MaxMultipartMemory = 30 << 20 // 30 MiB

	api := router.Group("/api/v1")
	{
		// Your existing public routes
		api.POST("/signup", SignupHandler)
		api.POST("/login", LoginHandler)
		api.GET("/validate-invite", ValidateInviteHandler)
		api.POST("/accept-invite", AcceptInviteHandler)
		api.POST("/refresh-token", RefreshTokenHandler)

		// This group is protected by your JWT middleware
		secured := api.Group("/")
		secured.Use(middleware.JWTAuthMiddleware())
		{
			// Your existing secured routes
			secured.POST("/invite", InviteHandler)
			secured.POST("/upload-resume", resumeHandler.Upload)
			secured.POST("/upload-cover-letter", resumeHandler.UploadCoverLetter)

			// --- ADD THIS NEW LINE ---
			// This creates the new route for the metadata.
			secured.POST("/upload-metadata", resumeHandler.UploadMetadata)
		}
	}

	return router
}
