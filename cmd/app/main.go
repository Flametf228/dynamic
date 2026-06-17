package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"zzz/database"
	"zzz/internal/config"
	"zzz/internal/handler"
	"zzz/internal/repository"
	"zzz/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	taskRepo := repository.NewTaskRepo(db)
	dataRepo := repository.NewDataRepo(db)
	statsRepo := repository.NewStatsRepo(db)

	taskService := service.NewTaskService(taskRepo, dataRepo, &cfg)
	statsService := service.NewStatsService(statsRepo)

	taskHandler := handler.NewTaskHandler(taskService)
	statsHandler := handler.NewStatsHandler(statsService)

	router := gin.Default()
	router.GET("/stats", statsHandler.GetStats)
	router.POST("/download", taskHandler.Download)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
