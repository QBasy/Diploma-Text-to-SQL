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

type CustomDatabase struct {
	gorm.Model
	UserUUID    string `gorm:"type:uuid;not null", json:"userUUID" json:"userUUID,omitempty"`
	User        User   `gorm:"foreignKey:UserUUID;references:UUID" json:"user"`
	UUID        string `gorm:"unique" json:"UUID,omitempty"`
	Name        string `json:"name,omitempty"`
	DBType      string `json:"DBType,omitempty"` // "postgres", "sqlite", "mysql"
	Host        string `json:"host,omitempty"`
	Port        int    `json:"port,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Database    string `json:"database,omitempty"`
	SSLMode     string `json:"SSLMode,omitempty"`
	Description string `json:"description,omitempty"`
}
