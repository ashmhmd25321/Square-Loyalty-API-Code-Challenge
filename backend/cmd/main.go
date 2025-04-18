package main

import (
	"log"
	"os"

	"loyalty-app/internal/api/handlers"
	"loyalty-app/internal/middleware"
	"loyalty-app/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	squareService := services.NewSquareService()

	dbService := services.NewMemoryService()
	log.Println("Running with in-memory storage")

	loyaltyHandler := handlers.NewLoyaltyHandler(squareService, dbService)
	authHandler := handlers.NewAuthHandler(dbService)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.POST("/register", authHandler.Register)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/earn", loyaltyHandler.EarnPoints)
			protected.POST("/redeem", loyaltyHandler.RedeemPoints)
			protected.GET("/balance", loyaltyHandler.GetBalance)
			protected.GET("/history", loyaltyHandler.GetHistory)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
