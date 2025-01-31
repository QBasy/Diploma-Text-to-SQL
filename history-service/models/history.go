package models

import (
	"time"
)

type QueryHistory struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID       string    `json:"user_id" gorm:"type:uuid"`
	DatabaseUUID string    `json:"database_uuid" gorm:"type:uuid"`
	QueryType    string    `json:"query_type"`
	Query        string    `json:"query"`
	Result       string    `json:"result" gorm:"type:jsonb"`
	Timestamp    time.Time `json:"timestamp" gorm:"autoCreateTime"`
}
