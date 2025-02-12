package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"metadata-service/controllers"
	"metadata-service/middleware"
	"metadata-service/models"
	"metadata-service/utils"
	"net/http"
)

var db *gorm.DB

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	dsn := "host=" + utils.GetEnv("DB_HOST", "") +
		" user=" + utils.GetEnv("DB_USER", "") +
		" password=" + utils.GetEnv("DB_PASSWORD", "") +
		" dbname=" + utils.GetEnv("DB_NAME", "") +
		" port=" + utils.GetEnv("DB_PORT", "") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("AuthService: Failed to connect Database %v", err)
	}
	_ = db.AutoMigrate(&models.DatabaseMetadata{}, &models.TableMetadata{}, &models.ColumnMetadata{})
}

func main() {
	r := gin.Default()

	metadataController := controllers.NewMetadataController(db)

	protected := r.Group("/api/metadata")
	protected.Use(middleware.AuthMiddleware)
	{
		protected.GET("/:user_id/:database_uuid", metadataController.GetMetadata)
		protected.POST("/:user_id/:database_uuid", metadataController.AddMetadata)
		protected.PUT("/:user_id/:database_uuid", metadataController.UpdateMetadata)
		protected.DELETE("/:user_id/:database_uuid", metadataController.DeleteMetadata)
	}
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Run(fmt.Sprintf(":%v", utils.GetEnv("SERVER_PORT", "5005")))
}
