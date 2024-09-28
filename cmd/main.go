package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ppondeu/go-post-api/internal/handler"
	"github.com/ppondeu/go-post-api/internal/repository"
	"github.com/ppondeu/go-post-api/internal/routes"
	"github.com/ppondeu/go-post-api/internal/usecase"
	"github.com/ppondeu/go-post-api/internal/validate"

	"github.com/ppondeu/go-post-api/config"
	database "github.com/ppondeu/go-post-api/internal/db"
)

func main() {
	cfg := config.LoadConfig()
	db := database.ConnectDatabase(cfg)
	validate := validate.NewValidator()

	userRepo := repository.NewUserRepositoryDB(db)
	userService := usecase.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, validate)

	followRepo := repository.NewFollowRepositoryDB(db)
	followService := usecase.NewFollowService(followRepo, userService)
	followHandler := handler.NewFollowHandler(followService, validate)

	jwtService := usecase.NewJwtService([]byte(cfg.ACCESS_SECRET), []byte(cfg.REFRESH_SECRET))
	authService := usecase.NewAuthService(userService, jwtService)
	authHandler := handler.NewAuthHandler(authService, validate)

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	routes.SetupUserRouter(router, userHandler)
	routes.SetupAuthRouter(router, authHandler, &jwtService)
	routes.SetupFollowRouter(router, followHandler)
	fmt.Printf("Server running on port %v", cfg.SERVER_PORT)
	router.Run(":" + cfg.SERVER_PORT)
}
