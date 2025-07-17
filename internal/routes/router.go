package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/resumelens/authservice/internal/middleware"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.POST("/signup", SignupHandler)
		api.POST("/login", LoginHandler)
		api.GET("/validate-invite", ValidateInviteHandler)
		api.POST("/accept-invite", AcceptInviteHandler)
		api.POST("/refresh-token", RefreshTokenHandler)

		secured := api.Group("/")
		secured.Use(middleware.JWTAuthMiddleware())
		{
			secured.POST("/invite", InviteHandler)
		}
	}

	return router
}
