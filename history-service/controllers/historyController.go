package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"history-service/models"
	"history-service/utils"
	"log"
	"math"
	"net/http"
	"strconv"
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

	// Parsing pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	// Validate pagination params
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 50 {
		perPage = 10
	}

	databaseUUID := c.Query("database_uuid")
	queryType := c.Query("query_type")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	searchQuery := c.Query("search")

	sortBy := c.DefaultQuery("sort_by", "timestamp")
	sortDir := c.DefaultQuery("sort_dir", "DESC")

	validSortFields := map[string]bool{
		"timestamp":  true,
		"query_type": true,
		"success":    true,
	}

	if !validSortFields[sortBy] {
		sortBy = "timestamp"
	}

	if sortDir != "ASC" && sortDir != "DESC" {
		sortDir = "DESC"
	}

	query := ctrl.DB.Where("user_id = ?", userUUID)

	if databaseUUID != "" {
		query = query.Where("database_uuid = ?", databaseUUID)
	}

	if queryType != "" {
		query = query.Where("query_type = ?", queryType)
	}

	if startDate != "" {
		query = query.Where("timestamp >= ?", startDate)
	}

	if endDate != "" {
		query = query.Where("timestamp <= ?", endDate)
	}

	if searchQuery != "" {
		query = query.Where("query ILIKE ?", "%"+searchQuery+"%")
	}

	var total int64
	if err := query.Model(&models.QueryHistory{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count history items"})
		return
	}

	var histories []models.QueryHistory
	if err := query.Order(fmt.Sprintf("%s %s", sortBy, sortDir)).
		Limit(perPage).
		Offset((page - 1) * perPage).
		Find(&histories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      histories,
		"page":      page,
		"per_page":  perPage,
		"total":     total,
		"last_page": int(math.Ceil(float64(total) / float64(perPage))),
	})
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
