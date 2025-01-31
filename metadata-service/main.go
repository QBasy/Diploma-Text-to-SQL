package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"metadata-service/controllers"
	"metadata-service/middleware"
	"metadata-service/models"
	"metadata-service/utils"
)

var db *gorm.DB

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	dsn := "host=" + utils.GetEnv("POSTGRES_HOST", "") +
		" user=" + utils.GetEnv("POSTGRES_USER", "") +
		" password=" + utils.GetEnv("POSTGRES_PASSWORD", "") +
		" dbname=" + utils.GetEnv("POSTGRES_DB", "") +
		" port=" + utils.GetEnv("POSTGRES_PORT", "") +
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

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware)
	{
		protected.GET("/metadata/:user_id/:database_uuid", metadataController.GetMetadata)
		protected.POST("/metadata/:user_id/:database_uuid", metadataController.AddMetadata)
		protected.PUT("/metadata/:user_id/:database_uuid", metadataController.UpdateMetadata)
		protected.DELETE("/metadata/:user_id/:database_uuid", metadataController.DeleteMetadata)
	}

	r.Run(utils.GetEnv("PORT", "5005"))
}
