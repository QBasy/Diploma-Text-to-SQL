package utils

import (
	"gorm.io/gorm"
	"history-service/models"
	"log"
)

func MigrateHistory(db *gorm.DB) {
	err := db.AutoMigrate(&models.QueryHistory{})
	if err != nil {
		log.Fatalf("Failed to migrate History table: %v", err)
	}
}
