package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ppondeu/go-post-api/internal/handler"
	"github.com/ppondeu/go-post-api/internal/middleware"
	"github.com/ppondeu/go-post-api/internal/usecase"
)

func SetupAuthRouter(router *gin.Engine, authHandler *handler.AuthHandler, jwtService *usecase.JwtService) {
	auth := router.Group("api/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", middleware.ValidateRefreshToken(*jwtService), authHandler.Logout)
		auth.POST("/refresh_token", middleware.ValidateRefreshToken(*jwtService), authHandler.RefreshToken)
	}
}
