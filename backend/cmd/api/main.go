package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"ticket-system/backend/internal/config"
	"ticket-system/backend/internal/delivery"
	"ticket-system/backend/internal/repository"
	"ticket-system/backend/internal/usecase"
)

func main() {
	cfg := config.LoadConfig()

	db, err := repository.NewDB(cfg.DBUrl)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	authService := usecase.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := delivery.NewAuthHandler(authService)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/register", authHandler.Register)

	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
