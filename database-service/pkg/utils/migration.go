package utils

import (
	models2 "database-service/internal/models"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models2.User{}, &models2.UserDatabase{}, &models2.CustomDatabase{}); err != nil {
		log.Println("Error during AutoMigrate:", err)
		return err
	}

	if err := db.Migrator().CreateConstraint(&models2.UserDatabase{}, "UserUUID"); err != nil {
		log.Println("Error adding foreign key:", err)
		return err
	}

	return nil
}
