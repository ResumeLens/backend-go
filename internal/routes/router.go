package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/resumelens/authservice/internal/middleware"

	"time"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()


	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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

