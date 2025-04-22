package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type QueryHistory struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	DatabaseUUID uuid.UUID `json:"database_uuid" gorm:"type:uuid;not null"`
	QueryType    string    `json:"query_type" gorm:"type:varchar(50);not null"`
	Query        string    `json:"query" gorm:"type:text;not null"`
	Result       string    `json:"result" gorm:"type:jsonb"`
	Timestamp    time.Time `json:"timestamp" gorm:"autoCreateTime"`
	Success      bool      `json:"success" gorm:"not null;default:false"`
}

func (q *QueryHistory) BeforeCreate(tx *gorm.DB) (err error) {
	q.ID = uuid.New()
	return
}
