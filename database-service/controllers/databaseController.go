package controllers

import (
	"database-service/models"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"path/filepath"
)

func (dc *DatabaseController) CreateDatabase(c *gin.Context) {
	var request CreateDatabaseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := dc.db.Where("uuid = ?", request.UserUUID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	databaseUUID := uuid.New().String()

	dbFile := databaseUUID + ".sqlite"
	dbPath := filepath.Join("databases", "user_databases", user.UUID, dbFile)

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Printf("Failed to create database directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	sqliteDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Printf("Failed to open database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SQLite database"})
		return
	}
	defer sqliteDB.Close()

	_, err = sqliteDB.Exec("CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY)")
	if err != nil {
		log.Printf("Failed to create table: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize database"})
		return
	}

	if err := dc.db.Exec("INSERT INTO user_databases (user_uuid, uuid, name, path) VALUES (?, ?, ?, ?)", user.UUID, databaseUUID, request.Name, dbPath).Error; err != nil {
		log.Printf("Failed to create database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save database metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database created successfully", "uuid": databaseUUID, "path": dbPath})
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
