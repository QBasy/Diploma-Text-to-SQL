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
	UserID uint   // Внешний ключ для связи с пользователем
	UUID   string `gorm:"unique"` // UUID базы данных
	Name   string // Название базы данных
	Path   string // Путь к SQLite файлу
}

type DatabaseRequest struct {
	SQLQuery string `json:"sql_query"`
}
