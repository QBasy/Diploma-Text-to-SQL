package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"metadata-service/models"
	"os"
)

var db *gorm.DB

func init() {
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
	db.AutoMigrate(&models.UserMetadata{}, &models.DatabaseMetadata{})
}

func main() {
	r := gin.Default()

	metadataController := controllers.NewMetadataController(db)

	log.Fatalf(r.Run(os.Getenv("PORT")))
}
