package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"history-service/models"
	"log"
	"net/http"
)

type HistoryController struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *HistoryController {
	return &HistoryController{DB: db}
}

func (ctrl *HistoryController) AddHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	var history models.QueryHistory

	if err := c.ShouldBindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	history.UserID = parseUUID(userID)
	if err := ctrl.DB.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "History saved successfully"})
}

func (ctrl *HistoryController) GetHistory(c *gin.Context) {
	userID := c.GetString("user_id")

	var histories []models.QueryHistory
	if err := ctrl.DB.Where("user_id = ?", userID).Order("timestamp DESC").Find(&histories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}

	c.JSON(http.StatusOK, histories)
}

func (ctrl *HistoryController) ClearHistory(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := ctrl.DB.Where("user_id = ?", userID).Delete(&models.QueryHistory{}).Error; err != nil {
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
