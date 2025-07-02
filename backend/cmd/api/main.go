package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.Run(":" + port)
} 