package main

import (
	"fmt"

	"github.com/ppondeu/go-post-api/api"
	"github.com/ppondeu/go-post-api/internal/adapter/http"
	"github.com/ppondeu/go-post-api/internal/repository"
	"github.com/ppondeu/go-post-api/internal/usecase"

	"github.com/ppondeu/go-post-api/config"
	database "github.com/ppondeu/go-post-api/db"
)

func main() {
	cfg := config.LoadConfig()
	db := database.ConnectDatabase(cfg)

	userRepo := repository.NewUserRepositoryDB(db)
	userService := usecase.NewUserService(userRepo)
	userHandler := http.NewUserHandler(userService)

	handlers := http.Handlers{
		UserHandler: userHandler,
	}
	router := api.SetupRouter(handlers)

	fmt.Printf("Server running on port %v", cfg.SERVER_PORT)
	router.Run(":" + cfg.SERVER_PORT)
}
