package models

import "gorm.io/gorm"

type Row struct {
	Name     string
	DataType string
	Null     bool
}

type Table struct {
	Name string
	Rows []Row
}

type UserDatabase struct {
	gorm.Model
	UserUUID string `gorm:"type:uuid;not null"`
	User     User   `gorm:"foreignKey:UserUUID;references:UUID"`
	UUID     string `gorm:"unique"`
	Name     string
	Path     string
}
