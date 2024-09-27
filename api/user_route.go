package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppondeu/go-post-api/internal/adapter/http"
)

func SetupUserRouter(router *gin.Engine, userHandler *http.UserHandler) {
	user := router.Group("api/users")
	{
		user.GET("/", userHandler.GetAllUsers)

		user.GET("/email/:email", userHandler.GetUserByEmail)
		user.GET("/username/:username", userHandler.GetUserByUsername)
		user.GET("/:id", userHandler.GetUserByID)

		user.POST("/", userHandler.CreateUser)
		user.PATCH("/:id", userHandler.UpdateUser)
		user.DELETE("/:id", userHandler.DeleteUser)
		// user.PATCH("/:id/session", userHandler.UpdateUserSession)
	}
}
