package controllers

import (
	"database-service/models"
	"database-service/utils"
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
	"strings"
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

func (dc *DatabaseController) ExecuteSQL(c *gin.Context) {
	userUUID, err := utils.GetUserUUIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request ExecuteSQLRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	var userDatabase models.UserDatabase
	if err := dc.db.Where("user_uuid = ?", userUUID).First(&userDatabase).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database not found"})
		return
	}

	sqliteDB, err := sql.Open("sqlite", userDatabase.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer sqliteDB.Close()

	if !(isSelectQuery(request.Query) || isCreateQuery(request.Query) || isDropQuery(request.Query) || isInsertQuery(request.Query) || isUpdateQuery(request.Query) || isDeleteQuery(request.Query)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query type not allowed"})
		return
	}

	if isSelectQuery(request.Query) {
		rows, err := sqliteDB.Query(request.Query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SQL query"})
			return
		}
		defer rows.Close()

		cols, err := rows.Columns()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch column information"})
			return
		}

		var result []map[string]interface{}
		for rows.Next() {
			columns := make([]interface{}, len(cols))
			columnPtrs := make([]interface{}, len(cols))
			for i := range columns {
				columnPtrs[i] = &columns[i]
			}

			if err := rows.Scan(columnPtrs...); err != nil {
				log.Printf("Error scanning row: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process row data"})
				return
			}

			rowMap := make(map[string]interface{})
			for i, col := range cols {
				rowMap[col] = columns[i]
			}
			result = append(result, rowMap)
		}

		c.JSON(http.StatusOK, gin.H{
			"row_count": len(result),
			"columns":   cols,
			"result":    result,
		})
		return
	}

	_, err = sqliteDB.Exec(request.Query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to execute SQL query"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Query executed successfully"})
}
func (dc *DatabaseController) GetDatabaseSchema(c *gin.Context) {
	userUUID, err := utils.GetUserUUIDFromRequest(c)
	if err != nil {
		return
	}

	var userDatabase models.UserDatabase
	if err := dc.db.Where("user_uuid = ?", userUUID).First(&userDatabase).Error; err != nil {
		log.Printf("Error fetching database for UUID %s: %v", userUUID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Database not found"})
		return
	}

	sqliteDB, err := sql.Open("sqlite", userDatabase.Path)
	if err != nil {
		log.Printf("Error connecting to SQLite database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to SQLite database"})
		return
	}
	defer sqliteDB.Close()

	rows, err := sqliteDB.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		log.Printf("Error fetching tables: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get schema"})
		return
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Printf("Error scanning table name: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tables = append(tables, tableName)
	}

	if len(tables) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"schema":  map[string][]string{},
			"message": "The database exists but does not contain any tables or connections.",
		})
		return
	}

	schema := make(map[string]map[string]interface{})
	for _, table := range tables {
		columnsRows, err := sqliteDB.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
		if err != nil {
			log.Printf("Error fetching columns for table %s: %v", table, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var columns []string
		var primaryKey string
		var foreignKeys []string

		for columnsRows.Next() {
			var cid int
			var name string
			var typeOfCol string
			var notNull bool
			var defaultVal interface{}
			var primaryKeyFlag bool

			if err := columnsRows.Scan(&cid, &name, &typeOfCol, &notNull, &defaultVal, &primaryKeyFlag); err != nil {
				log.Printf("Error scanning column info: %v", err)
				columnsRows.Close()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			columns = append(columns, name)
			if primaryKeyFlag {
				primaryKey = name
			}
		}
		columnsRows.Close()

		// Fetch foreign keys
		fkRows, err := sqliteDB.Query(fmt.Sprintf("PRAGMA foreign_key_list(%s)", table))
		if err != nil {
			log.Printf("Error fetching foreign keys for table %s: %v", table, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer fkRows.Close()

		for fkRows.Next() {
			var id int
			var seq int
			var table string
			var from string
			var to string
			var onUpdate string
			var onDelete string
			var match string
			if err := fkRows.Scan(&id, &seq, &table, &from, &to, &onUpdate, &onDelete, &match); err != nil {
				log.Printf("Error scanning foreign key info: %v", err)
				continue
			}
			foreignKeys = append(foreignKeys, fmt.Sprintf("%s -> %s(%s)", table, to, from))
		}

		schema[table] = map[string]interface{}{
			"columns":     columns,
			"primaryKey":  primaryKey,
			"foreignKeys": foreignKeys,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"schema": schema,
	})
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func isSelectQuery(query string) bool {
	return strings.HasPrefix(strings.ToUpper(query), "SELECT")
}

func isCreateQuery(query string) bool {
	return strings.HasPrefix(strings.ToUpper(query), "CREATE")
}

func isDropQuery(query string) bool {
	return strings.HasPrefix(strings.ToUpper(query), "DROP")
}

func isInsertQuery(query string) bool {
	return strings.HasPrefix(strings.ToUpper(query), "INSERT")
}

func isUpdateQuery(query string) bool {
	return strings.HasPrefix(strings.ToUpper(query), "UPDATE")
}

func isDeleteQuery(query string) bool {
	return strings.HasPrefix(strings.ToUpper(query), "DELETE")
}
