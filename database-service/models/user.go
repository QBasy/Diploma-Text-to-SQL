package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID         string `gorm:"type:uuid;unique;default:uuid_generate_v4()"` // UUID с default значением
	Username     string `gorm:"unique"`
	Email        string `gorm:"unique"`
	PasswordHash string
	Databases    []UserDatabase `gorm:"foreignKey:UserUUID;references:UUID"` // Связываем через UUID
}
