package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"history-service/models"
	"history-service/utils"
	"log"
	"net/http"
	"time"
)

type HistoryController struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *HistoryController {
	return &HistoryController{DB: db}
}

func (ctrl *HistoryController) AddHistory(c *gin.Context) {
	userUUID, err := utils.GetUserUUIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var databaseUUID uuid.UUID
	row := ctrl.DB.Raw("SELECT uuid FROM user_databases WHERE user_uuid = ? LIMIT 1", userUUID).Row()
	if err := row.Scan(&databaseUUID); err != nil {
		log.Printf("Error fetching database UUID for user %s: %v", userUUID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Database not found"})
		return
	}

	var history models.QueryHistory
	if err := c.ShouldBindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	history.UserID = parseUUID(userUUID)
	history.DatabaseUUID = databaseUUID
	if history.Timestamp.IsZero() {
		history.Timestamp = time.Now()
	}

	if err := ctrl.DB.Create(&history).Error; err != nil {
		fmt.Println("DB ERROR:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "History saved successfully"})
}

func (ctrl *HistoryController) GetHistory(c *gin.Context) {
	userUUID, err := utils.GetUserUUIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var histories []models.QueryHistory
	if err := ctrl.DB.Where("user_id = ?", userUUID).Order("timestamp DESC").Find(&histories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}

	c.JSON(http.StatusOK, histories)
}

func (ctrl *HistoryController) ClearHistory(c *gin.Context) {
	userUUID, err := utils.GetUserUUIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Where("user_id = ?", userUUID).Delete(&models.QueryHistory{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "History cleared successfully"})
}

func parseUUID(uuidStr string) uuid.UUID {
	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		log.Printf("Invalid UUID: %s, error: %v", uuidStr, err)
		return uuid.Nil
	}
	return parsedUUID
}
