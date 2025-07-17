package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.POST("/signup", SignupHandler)
		api.POST("/login", LoginHandler)
	}

	return router
}
