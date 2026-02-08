package main

import (
	"backend-ping-pong-app/internal/config"
	"backend-ping-pong-app/internal/database"
	//"backend-ping-pong-app/internal/handlers"
	//"backend-ping-pong-app/internal/middleware"
	//"backend-ping-pong-app/internal/repository"
	//"backend-ping-pong-app/internal/service"

	//"github.com/gin-gonic/gin"
	log "github.com/jeanphorn/log4go"
)

func initLogger() {
	log.AddFilter("stdout", log.DEBUG, log.NewConsoleLogWriter())
}

func main() {
	initLogger()

	cfg := config.Load()
	log.Debug(
		"DB env: user=%s password=%t",
		cfg.Database.User,
		cfg.Database.Password != "",
	)

	db, err := database.OpenPostgresDB(cfg.Database)
	if err != nil {
		log.Error("Database connection failed: %v", err)
	}
	defer db.Close()

	//repo := repository.NewRepository(db)
	//svc := service.NewService(repo)
	//
	//router := gin.New()
	//router.Use(
	//	gin.Logger(),
	//	gin.Recovery(),
	//	middleware.CORS(),
	//)

	//handlers.RegisterRoutes(router, svc)

	log.Info("🚀 Server running on :%s", cfg.App.Port)
	//if err := router.Run(":" + cfg.App.Port); err != nil {
	//	log.Fatal("Server failed: %v", err)
	//}
}
