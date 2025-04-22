package main

import (
	"API/middleware"
	"API/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	r := gin.Default()

	r.Use(middleware.RateLimiter())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	authRoutes(r)

	databaseRoutes(r)

	textToSQLRoutes(r)

	historyRoutes(r)

	r.GET("/api/health", handleHealthCheck)

	if err := r.Run(utils.GetEnv("PORT", ":5001")); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
