package models

type Item struct {
	ID      string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	OwnerID uint   `json:"owner_id"`
	Owner   User   `gorm:"foreignKey:OwnerID" json:"owner"`
}
