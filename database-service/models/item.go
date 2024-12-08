package models

type Item struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	OwnerID uint   `json:"owner_id"`
	Owner   User   `gorm:"foreignKey:OwnerID" json:"owner"`
}
