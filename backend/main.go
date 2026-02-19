package main

import (
	"log"
	"obucon/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Printf("Server starting on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
