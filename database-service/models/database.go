package models

type Row struct {
	Name     string
	DataType string
	Null     bool
}

type Table struct {
	Name string
	Rows []Row
}

type Database struct {
	ID     string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID string `gorm:"not null;type:uuid" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID" json:"user"`
	Path   string `gorm:"type:varchar(255);unique;not null" json:"file_path"`
}
