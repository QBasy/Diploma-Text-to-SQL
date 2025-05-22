package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID         string `gorm:"type:uuid;unique;default:uuid_generate_v4()"`
	Username     string `gorm:"unique"`
	Email        string `gorm:"unique"`
	PasswordHash string
	Role         string
	Databases    []UserDatabase `gorm:"foreignKey:UserUUID;references:UUID"`
}
