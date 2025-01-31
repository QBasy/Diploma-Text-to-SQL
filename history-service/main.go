package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"history-service/controllers"
	"history-service/models"
	"log"
	"os"
)

var db *gorm.DB

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}
	dsn := "host=" + os.Getenv("POSTGRES_HOST") +
		" user=" + os.Getenv("POSTGRES_USER") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_DB") +
		" port=" + os.Getenv("POSTGRES_PORT") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("AuthService: Failed to connect Database %v", err)
	}

	_ = db.AutoMigrate(&models.QueryHistory{})
}

func main() {
	r := gin.Default()

	historyController := controllers.NewHistoryController(db)

	history := r.Group("/api/history")
	{
		history.GET("", historyController.GetHistory)
		history.POST("", historyController.AddHistory)
		history.DELETE("", historyController.ClearHistory)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5007"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
