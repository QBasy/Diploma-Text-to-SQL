package controllers

import (
	"database-service/models"
	"database-service/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
)

type AddCustomDatabaseRequest struct {
	Name        string `json:"name"`
	DBType      string `json:"db_type"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Database    string `json:"database"`
	SSLMode     string `json:"ssl_mode"`
	Description string `json:"description"`
}

func (dc *DatabaseController) AddCustomDatabase(c *gin.Context) {
	var req AddCustomDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	userUUID, err := utils.GetUserUUIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	customDB := models.CustomDatabase{
		UserUUID:    userUUID,
		UUID:        uuid.New().String(),
		Name:        req.Name,
		DBType:      req.DBType,
		Host:        req.Host,
		Port:        req.Port,
		Username:    req.Username,
		Password:    req.Password,
		Database:    req.Database,
		SSLMode:     req.SSLMode,
		Description: req.Description,
	}

	if err := dc.db.Create(&customDB).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save custom database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Custom database added successfully", "uuid": customDB.UUID})
}

func (dc *DatabaseController) ListCustomDatabases(c *gin.Context) {
	userUUID, err := utils.GetUserUUIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var databases []models.CustomDatabase
	if err := dc.db.Where("user_uuid = ?", userUUID).Find(&databases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch custom databases"})
		return
	}

	c.JSON(http.StatusOK, databases)
}

func (dc *DatabaseController) DeleteCustomDatabase(c *gin.Context) {
	dbUUID := c.Param("uuid")
	userUUID, err := utils.GetUserUUIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := dc.db.Where("uuid = ? AND user_uuid = ?", dbUUID, userUUID).Delete(&models.CustomDatabase{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Custom database deleted"})
}

func (dc *DatabaseController) GetCustomDatabaseSchema(c *gin.Context) {
	dbUUID := c.Query("database_uuid")
	userUUID, err := utils.GetUserUUIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	dbConn, dbTypes, err := GetDatabaseConnection(dc.db, userUUID, dbUUID)
	if err != nil {
		log.Printf("DB connection error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer dbConn.Close()

	var tableQuery string
	switch dbTypes[0] {
	case "sqlite":
		tableQuery = "SELECT name FROM sqlite_master WHERE type='table'"
	case "postgres":
		tableQuery = "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
	case "mysql":
		tableQuery = "SHOW TABLES"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported DB type"})
		return
	}

	rows, err := dbConn.Query(tableQuery)
	if err != nil {
		log.Printf("Table query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tables"})
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

	schema := make(map[string]map[string]interface{})

	for _, table := range tables {
		var columns []string
		var primaryKey string
		var foreignKeys []ForeignKeyInfo

		switch dbTypes[0] {
		case "sqlite":
			colRows, err := dbConn.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
			if err != nil {
				log.Printf("PRAGMA table_info error: %v", err)
				continue
			}
			for colRows.Next() {
				var cid int
				var name, colType string
				var notNull, pk int
				var defaultVal sql.NullString
				if err := colRows.Scan(&cid, &name, &colType, &notNull, &defaultVal, &pk); err != nil {
					log.Printf("Error scanning column info: %v", err)
					continue
				}
				columns = append(columns, name)
				if pk > 0 {
					primaryKey = name
				}
			}
			colRows.Close()

			fkRows, err := dbConn.Query(fmt.Sprintf("PRAGMA foreign_key_list(%s)", table))
			if err == nil {
				for fkRows.Next() {
					var _, _, refTable, from, to string
					dummy := new(string)
					if err := fkRows.Scan(new(int), new(int), &refTable, &from, &to, dummy, dummy, dummy); err != nil {
						log.Printf("Error scanning FK: %v", err)
						continue
					}
					foreignKeys = append(foreignKeys, ForeignKeyInfo{Table: refTable, From: from, To: to})
				}
				fkRows.Close()
			}

		case "postgres":
			colRows, err := dbConn.Query(`
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_schema = 'public' AND table_name = $1`, table)
			if err != nil {
				log.Printf("Postgres column query error: %v", err)
				continue
			}
			for colRows.Next() {
				var name, dataType string
				if err := colRows.Scan(&name, &dataType); err != nil {
					log.Printf("Error scanning column: %v", err)
					continue
				}
				columns = append(columns, name)
			}
			colRows.Close()

			pkRows, err := dbConn.Query(`
		SELECT a.attname
		FROM pg_index i
		JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
		WHERE i.indrelid = $1::regclass AND i.indisprimary`, table)
			if err != nil {
				log.Printf("Error getting primary key: %v", err)
			} else {
				defer pkRows.Close()
				if pkRows.Next() {
					if err := pkRows.Scan(&primaryKey); err != nil {
						log.Printf("Error scanning primary key: %v", err)
					}
				}
			}

			fkQuery := `
		SELECT kcu.column_name, ccu.table_name AS foreign_table, ccu.column_name AS foreign_column
		FROM information_schema.table_constraints AS tc
		JOIN information_schema.key_column_usage AS kcu ON tc.constraint_name = kcu.constraint_name
		JOIN information_schema.constraint_column_usage AS ccu ON ccu.constraint_name = tc.constraint_name
		WHERE constraint_type = 'FOREIGN KEY' AND tc.table_name = $1`
			fkRows, err := dbConn.Query(fkQuery, table)
			if err == nil {
				for fkRows.Next() {
					var columnName, refTable, refCol string
					if err := fkRows.Scan(&columnName, &refTable, &refCol); err != nil {
						log.Printf("Error scanning FK: %v", err)
						continue
					}
					foreignKeys = append(foreignKeys, ForeignKeyInfo{Table: refTable, From: columnName, To: refCol})
				}
				fkRows.Close()
			}

		case "mysql":
			colRows, err := dbConn.Query(fmt.Sprintf("SHOW COLUMNS FROM %s", table))
			if err != nil {
				log.Printf("MySQL column query error: %v", err)
				continue
			}
			for colRows.Next() {
				var field, colType, null, key, defVal, extra string
				if err := colRows.Scan(&field, &colType, &null, &key, &defVal, &extra); err != nil {
					log.Printf("Error scanning column: %v", err)
					continue
				}
				columns = append(columns, field)
				if strings.ToUpper(key) == "PRI" {
					primaryKey = field
				}
			}
			colRows.Close()

			fkRows, err := dbConn.Query(fmt.Sprintf(`
				SELECT kcu.column_name, kcu.referenced_table_name, kcu.referenced_column_name
				FROM information_schema.key_column_usage AS kcu
				JOIN information_schema.table_constraints AS tc ON kcu.constraint_name = tc.constraint_name
				WHERE tc.constraint_type = 'FOREIGN KEY' AND kcu.table_name = '%s'`, table))
			if err == nil {
				for fkRows.Next() {
					var columnName, refTable, refCol string
					if err := fkRows.Scan(&columnName, &refTable, &refCol); err != nil {
						log.Printf("Error scanning FK: %v", err)
						continue
					}
					foreignKeys = append(foreignKeys, ForeignKeyInfo{Table: refTable, From: columnName, To: refCol})
				}
				fkRows.Close()
			}
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

func (dc *DatabaseController) GetFullCustomDatabaseSchema(c *gin.Context) {
	dbUUID := c.Query("database_uuid")
	userUUID, err := utils.GetUserUUIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	dbConn, dbTypes, err := GetDatabaseConnection(dc.db, userUUID, dbUUID)
	if err != nil {
		log.Printf("Error getting database connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return
	}
	defer dbConn.Close()

	var tables []string
	var errQueryTables error

	switch dbTypes[0] {
	case "sqlite":
		rows, err := dbConn.Query("SELECT name FROM sqlite_master WHERE type='table'")
		errQueryTables = err
		if rows != nil {
			defer rows.Close()
			for rows.Next() {
				var tableName string
				if err := rows.Scan(&tableName); err != nil {
					log.Printf("Error scanning table name: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				tables = append(tables, tableName)
			}
		}
	case "postgres":
		rows, err := dbConn.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
		errQueryTables = err
		if rows != nil {
			defer rows.Close()
			for rows.Next() {
				var tableName string
				if err := rows.Scan(&tableName); err != nil {
					log.Printf("Error scanning table name: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				tables = append(tables, tableName)
			}
		}
	case "mysql":
		rows, err := dbConn.Query("SHOW TABLES")
		errQueryTables = err
		if rows != nil {
			defer rows.Close()
			for rows.Next() {
				var tableName string
				if err := rows.Scan(&tableName); err != nil {
					log.Printf("Error scanning table name: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				tables = append(tables, tableName)
			}
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unsupported database type: %s", dbTypes[0])})
		return
	}

	if errQueryTables != nil {
		log.Printf("Error fetching tables: %v", errQueryTables)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get schema"})
		return
	}

	var fullSchema FullSchemaResponse
	fullSchema.Tables = make([]TableInfo, 0, len(tables))

	for _, table := range tables {
		var columns []ColumnInfo
		var primaryKey string

		switch dbTypes[0] {
		case "sqlite":
			columnsRows, err := dbConn.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
			if err != nil {
				log.Printf("Error fetching columns for table %s: %v", table, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer columnsRows.Close()
			for columnsRows.Next() {
				var cid int
				var name, typeOfCol string
				var notNull, primaryKeyFlag bool
				var defaultVal interface{}

				if err := columnsRows.Scan(&cid, &name, &typeOfCol, &notNull, &defaultVal, &primaryKeyFlag); err != nil {
					log.Printf("Error scanning column info: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				columns = append(columns, ColumnInfo{
					Name: name,
					Type: typeOfCol,
				})
				if primaryKeyFlag {
					primaryKey = name
				}
			}

			fkRows, err := dbConn.Query(fmt.Sprintf("PRAGMA foreign_key_list(%s)", table))
			if err == nil && fkRows != nil {
				defer fkRows.Close()
				for fkRows.Next() {
					var refTable, fromCol, toCol string
					if err := fkRows.Scan(new(int), new(int), &refTable, &fromCol, &toCol, new(string), new(string), new(string)); err != nil {
						log.Printf("Error scanning foreign key info: %v", err)
						continue
					}
					for i := range columns {
						if columns[i].Name == fromCol {
							columns[i].IsForeignKey = true
							columns[i].ReferencedTable = refTable
							columns[i].ReferencedColumn = toCol
							break
						}
					}
				}
			}

		case "postgres":
			columnsRows, err := dbConn.Query(`SELECT column_name, data_type
		FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = $1`, table)
			if err != nil {
				log.Printf("Error fetching columns for table %s: %v", table, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer columnsRows.Close()

			for columnsRows.Next() {
				var columnName, dataType string
				if err := columnsRows.Scan(&columnName, &dataType); err != nil {
					log.Printf("Error scanning column info: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				columns = append(columns, ColumnInfo{
					Name: columnName,
					Type: dataType,
				})
			}

			var primaryKey string
			pkRows, err := dbConn.Query(`
		SELECT a.attname
		FROM pg_index i
		JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
		WHERE i.indrelid = $1::regclass AND i.indisprimary`, table)
			if err != nil {
				log.Printf("Error fetching primary key for table %s: %v", table, err)
			} else {
				defer pkRows.Close()
				if pkRows.Next() {
					if err := pkRows.Scan(&primaryKey); err != nil {
						log.Printf("Error scanning primary key: %v", err)
					}
				}
			}

			fkRows, err := dbConn.Query(`
		SELECT kcu.column_name, ccu.table_name AS referenced_table, ccu.column_name AS referenced_column
		FROM information_schema.table_constraints AS tc
		JOIN information_schema.key_column_usage AS kcu
			ON tc.constraint_name = kcu.constraint_name AND tc.table_schema = kcu.table_schema AND tc.table_name = kcu.table_name
		JOIN information_schema.referential_constraints AS rc
			ON tc.constraint_name = rc.constraint_name AND tc.table_schema = rc.constraint_schema AND tc.table_name = rc.table_name
		JOIN information_schema.key_column_usage AS ccu
			ON rc.unique_constraint_name = ccu.constraint_name AND rc.unique_constraint_schema = ccu.constraint_schema AND rc.unique_constraint_table_name = ccu.table_name AND kcu.ordinal_position = ccu.ordinal_position
		WHERE tc.constraint_type = 'FOREIGN KEY' AND kcu.table_name = $1`, table)
			if err == nil && fkRows != nil {
				defer fkRows.Close()
				for fkRows.Next() {
					var fkColumn, referencedTable, referencedColumn string
					if err := fkRows.Scan(&fkColumn, &referencedTable, &referencedColumn); err != nil {
						log.Printf("Error scanning foreign key info: %v", err)
						continue
					}
					for i := range columns {
						if columns[i].Name == fkColumn {
							columns[i].IsForeignKey = true
							columns[i].ReferencedTable = referencedTable
							columns[i].ReferencedColumn = referencedColumn
							break
						}
					}
				}
			}

			if primaryKey != "" {
				for i := range columns {
					if columns[i].Name == primaryKey {
						columns[i].IsForeignKey = true
						columns[i].ReferencedTable = "PRIMARY KEY"
						columns[i].ReferencedColumn = primaryKey
						break
					}
				}
			}

		case "mysql":
			columnsRows, err := dbConn.Query(fmt.Sprintf("SHOW COLUMNS FROM %s", table))
			if err != nil {
				log.Printf("Error fetching columns for table %s: %v", table, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer columnsRows.Close()
			for columnsRows.Next() {
				var field, colType, key string
				if err := columnsRows.Scan(&field, &colType, new(string), &key, new(interface{}), new(string)); err != nil {
					log.Printf("Error scanning column info: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				columns = append(columns, ColumnInfo{
					Name: field,
					Type: colType,
				})
				if strings.ToUpper(key) == "PRI" {
					primaryKey = field
				}
			}

			fkRows, err := dbConn.Query(fmt.Sprintf(`SELECT kcu.column_name, kcu.referenced_table_name, kcu.referenced_column_name
				FROM information_schema.key_column_usage AS kcu
				JOIN information_schema.table_constraints AS tc ON kcu.constraint_name = tc.constraint_name
				WHERE tc.constraint_type = 'FOREIGN KEY' AND kcu.table_name = '%s'`, table))
			if err == nil && fkRows != nil {
				defer fkRows.Close()
				for fkRows.Next() {
					var fkColumn, referencedTable, referencedColumn string
					if err := fkRows.Scan(&fkColumn, &referencedTable, &referencedColumn); err != nil {
						log.Printf("Error scanning foreign key info: %v", err)
						continue
					}
					for i := range columns {
						if columns[i].Name == fkColumn {
							columns[i].IsForeignKey = true
							columns[i].ReferencedTable = referencedTable
							columns[i].ReferencedColumn = referencedColumn
							break
						}
					}
				}
			}

		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unsupported database type: %s", dbTypes[0])})
			return
		}

		fullSchema.Tables = append(fullSchema.Tables, TableInfo{
			Name:       table,
			Columns:    columns,
			PrimaryKey: primaryKey,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"schema": fullSchema.Tables,
	})
}
