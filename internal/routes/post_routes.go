package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ppondeu/go-post-api/internal/handler"
)

func SetupPostRouter(router *gin.Engine, postHander *handler.PostHandler) {
	post := router.Group("api/posts")
	{
		post.GET("/", postHander.GetAllPosts)
		post.GET("/:id", postHander.GetPostByID)
		post.GET("/user/:id", postHander.GetPostsByUserID)
		post.POST("/", postHander.CreatePost)
		post.PATCH("/:id", postHander.UpdatePost)
		post.DELETE("/:id", postHander.DeletePost)

		post.GET("/tags", postHander.GetTags)
	}
}
