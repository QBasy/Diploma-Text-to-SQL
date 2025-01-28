package controllers

import (
	"database-service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

// Запрос для создания базы данных
type CreateDatabaseRequest struct {
	UserUUID string `json:"user_uuid"`
	Name     string `json:"name"`
}

// Запрос для выполнения SQL
type ExecuteSQLRequest struct {
	UserUUID     string `json:"user_uuid"`
	DatabaseUUID string `json:"database_uuid"`
	Query        string `json:"query"`
}

// Создание новой базы данных для пользователя
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

	dbPath := filepath.Join("static", "user_databases", user.UUID, databaseUUID+".sqlite")

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	_, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SQLite database"})
		return
	}

	userDatabase := models.UserDatabase{
		UserID: user.ID,
		UUID:   databaseUUID,
		Name:   request.Name,
		Path:   dbPath,
	}
	if err := dc.db.Create(&userDatabase).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save database metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database created successfully", "uuid": databaseUUID, "path": dbPath})
}

// Выполнение SQL-запроса на базе данных пользователя
func (dc *DatabaseController) ExecuteSQL(c *gin.Context) {
	var request ExecuteSQLRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Поиск информации о базе данных в PostgreSQL
	var userDatabase models.UserDatabase
	if err := dc.db.Where("uuid = ?", request.DatabaseUUID).First(&userDatabase).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database not found"})
		return
	}

	// Подключение к SQLite базе данных
	sqliteDB, err := gorm.Open(sqlite.Open(userDatabase.Path), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to SQLite database"})
		return
	}

	// Выполнение SQL-запроса
	var result []map[string]interface{}
	if err := sqliteDB.Raw(request.Query).Scan(&result).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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

	sqliteDB, err := gorm.Open(sqlite.Open(userDatabase.Path), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to SQLite database"})
		return
	}

	var tables []string
	sqliteDB.Raw("SELECT name FROM sqlite_master WHERE type='table'").Scan(&tables)

	schema := make(map[string][]string)
	for _, table := range tables {
		var columns []string
		sqliteDB.Raw("PRAGMA table_info(" + table + ")").Scan(&columns)
		schema[table] = columns
	}

	c.JSON(http.StatusOK, gin.H{"schema": schema})
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
