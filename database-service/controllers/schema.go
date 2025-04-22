package controllers

import (
	"database-service/models"
	"database-service/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (dc *DatabaseController) GetDatabaseSchema(c *gin.Context) {
	userUUID, err := utils.GetUserUUIDFromRequest(c)
	log.Println("GetUserUUIDFromRequest:", userUUID)
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
		var foreignKeys []ForeignKeyInfo

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
			var refTable string
			var from string
			var to string
			var onUpdate string
			var onDelete string
			var match string
			if err := fkRows.Scan(&id, &seq, &refTable, &from, &to, &onUpdate, &onDelete, &match); err != nil {
				log.Printf("Error scanning foreign key info: %v", err)
				continue
			}
			foreignKeys = append(foreignKeys, ForeignKeyInfo{
				Table: refTable,
				From:  from,
				To:    to,
			})
		}

		schema[table] = map[string]interface{}{
			"columns":     columns,
			"primaryKey":  primaryKey,
			"foreignKeys": foreignKeys,
		}
		log.Printf("Foreign keys for table %s: %+v", table, foreignKeys)
	}

	log.Printf("Final Schema: %+v", schema)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"schema": schema,
	})
}

func (dc *DatabaseController) GetFullDatabaseSchema(c *gin.Context) {
	userUUID, err := utils.GetUserUUIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	rows.Close()

	var fullSchema FullSchemaResponse

	for _, table := range tables {
		columnsQuery := fmt.Sprintf("PRAGMA table_info(%s)", table)
		columnsRows, err := sqliteDB.Query(columnsQuery)
		if err != nil {
			log.Printf("Error fetching columns for table %s: %v", table, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var columns []ColumnInfo
		var primaryKey string

		for columnsRows.Next() {
			var cid int
			var name, typeOfCol string
			var notNull, primaryKeyFlag bool

			if err := columnsRows.Scan(&cid, &name, &typeOfCol, &notNull, new(interface{}), &primaryKeyFlag); err != nil {
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
		columnsRows.Close()

		fkQuery := fmt.Sprintf("PRAGMA foreign_key_list(%s)", table)
		fkRows, err := sqliteDB.Query(fkQuery)
		if err != nil {
			log.Printf("Error fetching foreign keys for table %s: %v", table, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

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
		fkRows.Close()

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
