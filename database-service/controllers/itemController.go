package controllers

import (
	"database-service/models"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
)

type ItemController struct {
	DB *gorm.DB
}

func (ic *ItemController) CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ic.DB.Create(&item).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (ic *ItemController) GetItems(w http.ResponseWriter, r *http.Request) {
	var items []models.Item
	if err := ic.DB.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(items)
}
