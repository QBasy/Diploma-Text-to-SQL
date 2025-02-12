package controllers

import (
	"database-service/models"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"path/filepath"
)

type DatabaseController struct {
	db *gorm.DB
}

func NewDatabaseController(db *gorm.DB) *DatabaseController {
	return &DatabaseController{db: db}
}

type CreateDatabaseRequest struct {
	UserUUID string `json:"user_uuid"`
	Name     string `json:"name"`
}

type ExecuteSQLRequest struct {
	UserUUID     string `json:"user_uuid"`
	DatabaseUUID string `json:"database_uuid"`
	Query        string `json:"query"`
}

func (dc *DatabaseController) CreateDatabase(c *gin.Context) {
	var request CreateDatabaseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Using GORM to find the user in PostgreSQL
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

	if err := dc.db.Exec("INSERT INTO user_databases (user_id, uuid, name, path) VALUES (?, ?, ?, ?)", user.ID, databaseUUID, request.Name, dbPath).Error; err != nil {
		log.Printf("Failed to create database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save database metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database created successfully", "uuid": databaseUUID, "path": dbPath})
}

func (dc *DatabaseController) ExecuteSQL(c *gin.Context) {
	var request ExecuteSQLRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userDatabase models.UserDatabase
	if err := dc.db.Where("uuid = ?", request.DatabaseUUID).First(&userDatabase).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database not found"})
		return
	}

	sqliteDB, err := sql.Open("sqlite", userDatabase.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to SQLite database"})
		return
	}
	defer sqliteDB.Close()

	rows, err := sqliteDB.Query(request.Query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	cols, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnsPtrs := make([]interface{}, len(cols))
		for i := range columns {
			columnsPtrs[i] = &columns[i]
		}

		if err := rows.Scan(columnsPtrs...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowMap := make(map[string]interface{})
		for i, col := range cols {
			rowMap[col] = columns[i]
		}
		result = append(result, rowMap)
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (dc *DatabaseController) GetDatabaseSchema(c *gin.Context) {
	databaseUUID := c.Param("database_uuid")

	var userDatabase models.UserDatabase
	if err := dc.db.Where("uuid = ?", databaseUUID).First(&userDatabase).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database not found"})
		return
	}

	sqliteDB, err := sql.Open("sqlite", userDatabase.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to SQLite database"})
		return
	}
	defer sqliteDB.Close()

	rows, err := sqliteDB.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get schema"})
		return
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tables = append(tables, tableName)
	}

	schema := make(map[string][]string)
	for _, table := range tables {
		columnsRows, err := sqliteDB.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer columnsRows.Close()

		var columns []string
		for columnsRows.Next() {
			var cid int
			var name string
			var typeOfCol string
			var notNull bool
			var defaultVal interface{}
			var primaryKey bool

			if err := columnsRows.Scan(&cid, &name, &typeOfCol, &notNull, &defaultVal, &primaryKey); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			columns = append(columns, name)
		}
		schema[table] = columns
	}

	c.JSON(http.StatusOK, gin.H{"schema": schema})
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
