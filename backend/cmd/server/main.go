package main

import (
	"log"
	"obucon/internal/config"
	"obucon/internal/database"
	"obucon/internal/handlers"
	"obucon/internal/nlp/japanese"
	"obucon/internal/repository"
	"obucon/internal/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if err := database.RunMigrations(cfg); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	tokenizer, err := japanese.NewTokenizer()
	if err != nil {
		log.Fatalf("Failed to initialize tokenizer: %v", err)
	}

	analysisService := services.NewAnalysisService(
		tokenizer,
		nil,
		nil,
	)
	analysisHandler := handlers.NewAnalysisHandler(analysisService)

	// router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"}, // will change upon dpeloyment
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth endpoints
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)
	r.POST("/auth/logout", authHandler.Logout)

	// Protected endpoints (require authentication)
	protected := r.Group("/")
	protected.Use(handlers.AuthMiddleware(authService))
	{
		protected.GET("/auth/me", authHandler.GetMe)
		protected.POST("/analyze", analysisHandler.AnalyzeText)
	}

	log.Printf("Server starting on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
