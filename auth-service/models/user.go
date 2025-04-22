package models

import (
	"gorm.io/gorm"
	"time"
)

type UserDatabase struct {
	gorm.Model
	UserUUID string `gorm:"type:uuid;not null"`                  // UUID внешнего ключа
	User     User   `gorm:"foreignKey:UserUUID;references:UUID"` // Явная связь через UUID
	UUID     string `gorm:"unique"`
	Name     string
	Path     string
}

type User struct {
	gorm.Model
	UUID         string `gorm:"type:uuid;unique;default:uuid_generate_v4()"`
	Username     string `gorm:"unique"`
	Email        string `gorm:"unique"`
	PasswordHash string
	Databases    []UserDatabase `gorm:"foreignKey:UserUUID;references:UUID"`
}

type EmailVerificationToken struct {
	gorm.Model
	UserUUID string `gorm:"type:uuid;uniqueIndex"`
	Token    string `gorm:"uniqueIndex"`
	Expiry   time.Time
}

type PasswordResetToken struct {
	gorm.Model
	UserUUID string `gorm:"type:uuid;unique"`
	Token    string `gorm:"unique"`
	Expiry   time.Time
}
