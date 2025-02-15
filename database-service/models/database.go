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
	UserUUID string `gorm:"type:uuid;not null"`                  // UUID внешнего ключа
	User     User   `gorm:"foreignKey:UserUUID;references:UUID"` // Явная связь через UUID
	UUID     string `gorm:"unique"`
	Name     string
	Path     string
}

type DatabaseRequest struct {
	SQLQuery string `json:"sql_query"`
}
