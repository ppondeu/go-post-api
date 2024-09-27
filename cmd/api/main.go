package main

import (
	"fmt"

	"github.com/ppondeu/go-post-api/api"
	"github.com/ppondeu/go-post-api/config"
	// database "github.com/ppondeu/go-post-api/db"
	"github.com/ppondeu/go-post-api/logs"
)

func main() {
	cfg := config.LoadConfig()
	// db := database.ConnectDatabase(cfg)
	router := api.SetupRouter()

	logs.Info(fmt.Sprintf("Server running on port %v", cfg.SERVER_PORT))
	router.Run(":" + cfg.SERVER_PORT)
}
