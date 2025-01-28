package controllers

import "gorm.io/gorm"

type MetadataController struct {
	db gorm.DB
}

func NewMetadataController(db gorm.DB) *MetadataController {
	return &MetadataController{db: db}
}
