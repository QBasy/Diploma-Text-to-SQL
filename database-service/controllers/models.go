package controllers

import "gorm.io/gorm"

type ForeignKeyInfo struct {
	Table string `json:"table"`
	From  string `json:"from"`
	To    string `json:"to"`
}

type ColumnInfo struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	IsForeignKey     bool   `json:"isForeignKey"`
	ReferencedTable  string `json:"referencedTable,omitempty"`
	ReferencedColumn string `json:"referencedColumn,omitempty"`
}

type TableInfo struct {
	Name       string       `json:"name"`
	Columns    []ColumnInfo `json:"columns"`
	PrimaryKey string       `json:"primaryKey"`
}

type FullSchemaResponse struct {
	Tables []TableInfo `json:"tables"`
}

type DatabaseRequest struct {
	SQLQuery string `json:"sql_query"`
}

type CreateDatabaseRequest struct {
	UserUUID string `json:"user_uuid"`
	Name     string `json:"name"`
}

type VisualisationRequest struct {
	Query string `json:"query"`
}

type ExecuteSQLRequest struct {
	Query string `json:"query"`
}

type DatabaseController struct {
	db *gorm.DB
}

func NewDatabaseController(db *gorm.DB) *DatabaseController {
	return &DatabaseController{db: db}
}
