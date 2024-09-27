package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppondeu/go-post-api/internal/adapter/http"
)

func SetupRouter(handers http.Handlers) *gin.Engine {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	SetupUserRouter(router, handers.UserHandler)

	return router
}
