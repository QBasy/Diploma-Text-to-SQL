package utils

import (
	"database-service/models"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.User{}, &models.UserDatabase{}); err != nil {
		log.Println("Error during AutoMigrate:", err)
		return err
	}

	if err := db.Migrator().CreateConstraint(&models.UserDatabase{}, "UserUUID"); err != nil {
		log.Println("Error adding foreign key:", err)
		return err
	}

	return nil
}
