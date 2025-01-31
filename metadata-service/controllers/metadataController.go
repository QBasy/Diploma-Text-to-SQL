package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"metadata-service/models"
	"net/http"
)

type MetadataController struct {
	db *gorm.DB
}

func NewMetadataController(db *gorm.DB) *MetadataController {
	return &MetadataController{db: db}
}

func (md *MetadataController) isUserAuthorized(c *gin.Context, userIDParam string, databaseUUID string) bool {
	currentUserID := c.GetString("user_id")

	if currentUserID != userIDParam {
		return false
	}

	var user models.User
	if err := md.db.Preload("Databases").Where("id = ?", currentUserID).First(&user).Error; err != nil {
		return false
	}

	// Ensure the user owns the requested database
	for _, userDatabase := range user.Databases {
		if userDatabase.UUID == databaseUUID {
			return true
		}
	}

	return false
}

func (md *MetadataController) GetMetadata(c *gin.Context) {
	userID := c.Param("user_id")
	databaseUUID := c.Param("database_uuid")

	if !md.isUserAuthorized(c, userID, databaseUUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this database"})
		return
	}

	var metadata models.DatabaseMetadata
	if err := md.db.Where("database_uuid = ?", databaseUUID).Preload("Tables.Columns").First(&metadata).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Metadata not found"})
		return
	}

	c.JSON(http.StatusOK, metadata)
}

func (md *MetadataController) AddMetadata(c *gin.Context) {
	userID := c.Param("user_id")
	databaseUUID := c.Param("database_uuid")

	if !md.isUserAuthorized(c, userID, databaseUUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to modify this database"})
		return
	}

	var metadata models.DatabaseMetadata
	if err := c.ShouldBindJSON(&metadata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	metadata.DatabaseUUID = databaseUUID
	if err := md.db.Create(&metadata).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store metadata"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Metadata stored successfully"})
}

func (md *MetadataController) UpdateMetadata(c *gin.Context) {
	userID := c.Param("user_id")
	databaseUUID := c.Param("database_uuid")

	if !md.isUserAuthorized(c, userID, databaseUUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this database"})
		return
	}

	var updatedMetadata models.DatabaseMetadata
	if err := c.ShouldBindJSON(&updatedMetadata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var existingMetadata models.DatabaseMetadata
	if err := md.db.Where("database_uuid = ?", databaseUUID).First(&existingMetadata).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Metadata not found"})
		return
	}

	existingMetadata.Tables = updatedMetadata.Tables
	if err := md.db.Save(&existingMetadata).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Metadata updated successfully"})
}

func (md *MetadataController) DeleteMetadata(c *gin.Context) {
	userID := c.Param("user_id")
	databaseUUID := c.Param("database_uuid")

	if !md.isUserAuthorized(c, userID, databaseUUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this database"})
		return
	}

	if err := md.db.Where("database_uuid = ?", databaseUUID).Delete(&models.DatabaseMetadata{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Metadata deleted successfully"})
}
