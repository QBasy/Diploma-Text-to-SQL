package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UUID         string `gorm:"unique"`
	Username     string `gorm:"unique"`
	Email        string `gorm:"unique"`
	PasswordHash string
	Databases    []UserDatabase
}

type UserDatabase struct {
	gorm.Model
	UserID uint
	UUID   string `gorm:"unique"`
	Name   string
	Path   string
}
