package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ppondeu/go-post-api/internal/handler"
	"github.com/ppondeu/go-post-api/internal/repository"
	"github.com/ppondeu/go-post-api/internal/routes"
	"github.com/ppondeu/go-post-api/internal/usecase"

	"github.com/ppondeu/go-post-api/config"
	database "github.com/ppondeu/go-post-api/internal/db"
)

func main() {
	cfg := config.LoadConfig()
	db := database.ConnectDatabase(cfg)
	validate := validator.New()
	userRepo := repository.NewUserRepositoryDB(db)
	userService := usecase.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, validate)

	authService := usecase.NewAuthService(userService, []byte(cfg.ACCESS_SECRET), []byte(cfg.REFRESH_SECRET))
	authHandler := handler.NewAuthHandler(authService, validate)

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	routes.SetupUserRouter(router, userHandler)
	routes.SetupAuthRouter(router, authHandler)

	fmt.Printf("Server running on port %v", cfg.SERVER_PORT)
	router.Run(":" + cfg.SERVER_PORT)
}
