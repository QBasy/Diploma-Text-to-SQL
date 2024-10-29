package models

import "gorm.io/gorm"

type User struct {
	gorm.DB
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
