package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ppondeu/go-post-api/internal/handler"
)

func SetupUserRouter(router *gin.Engine, userHandler *handler.UserHandler) {
	user := router.Group("api/users")
	{
		user.GET("/", userHandler.GetAllUsers)

		user.GET("/email/:email", userHandler.GetUserByEmail)
		user.GET("/username/:username", userHandler.GetUserByUsername)
		user.GET("/:id", userHandler.GetUserByID)

		user.POST("/", userHandler.CreateUser)
		user.PATCH("/:id", userHandler.UpdateUser)
		user.DELETE("/:id", userHandler.DeleteUser)
		user.PATCH("/session/:id", userHandler.UpdateUserSession)

		user.GET("/test", userHandler.GetUsersWithRelation)
		user.GET("/test/:id", userHandler.GetUserWithRelation)

		user.GET("/bookmarks/:id", userHandler.GetUserBookmarks)
	}
}
