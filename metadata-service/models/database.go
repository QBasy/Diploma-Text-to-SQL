package models

import "gorm.io/gorm"

type DatabaseMetadata struct {
	gorm.Model
	DatabaseUUID string `gorm:"unique"`
	Tables       []TableMetadata
}

type TableMetadata struct {
	gorm.Model
	DatabaseMetadataID uint
	Name               string
	Columns            []ColumnMetadata
}

type ColumnMetadata struct {
	gorm.Model
	TableMetadataID uint
	Name            string
	Type            string
	IsNullable      bool
	IsPrimaryKey    bool
}
