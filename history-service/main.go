package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"history-service/controllers"
	"history-service/models"
	"history-service/utils"
	"log"
	"net/http"
	"os"
)

var db *gorm.DB

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}
	dsn := "host=" + utils.GetEnv("DB_HOST", "localhost") +
		" user=" + utils.GetEnv("DB_USER", "postgres") +
		" password=" + utils.GetEnv("DB_PASSWORD", "") +
		" dbname=" + utils.GetEnv("DB_NAME", "mvpdiploma") +
		" port=" + utils.GetEnv("DB_PORT", "5432") +
		" sslmode=disable"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("AuthService: Failed to connect Database %v", err)
	}

	_ = db.AutoMigrate(&models.QueryHistory{})
}

func main() {
	r := gin.Default()

	historyController := controllers.New(db)

	history := r.Group("/api/history")
	{
		history.GET("", historyController.GetHistory)
		history.POST("", historyController.AddHistory)
		history.DELETE("", historyController.ClearHistory)
	}
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5007"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
