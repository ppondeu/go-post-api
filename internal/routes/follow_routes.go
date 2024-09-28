package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ppondeu/go-post-api/internal/handler"
)

func SetupFollowRouter(router *gin.Engine, userHandler *handler.FollowHandler) {
	follow := router.Group("api/follow")
	{
		follow.POST("/", userHandler.Follow)
		follow.DELETE("/", userHandler.Unfollow)
		follow.GET("/followers/:id", userHandler.GetFollowers)
		follow.GET("/followed/:id", userHandler.GetFollowedUsers)
	}
}
