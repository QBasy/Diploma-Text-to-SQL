package models

import (
	"gorm.io/gorm"
	"time"
)

type UserDatabase struct {
	gorm.Model
	UserID uint
	UUID   string `gorm:"unique"`
	Name   string
	Path   string
}

type User struct {
	gorm.Model
	UUID         string `gorm:"unique"`
	Username     string `gorm:"unique"`
	Email        string `gorm:"unique"`
	PasswordHash string
	Databases    []UserDatabase
}

type PasswordResetToken struct {
	gorm.Model
	UserID uint
	Token  string `gorm:"unique"`
	Expiry time.Time
}
