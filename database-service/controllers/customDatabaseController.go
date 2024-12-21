package controllers

import (
	"database-service/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type CustomDatabaseController struct {
	DB          *gorm.DB
	RawDatabase *sql.DB
}

func (ctrl CustomDatabaseController) CreateTable(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)

	if !ok {
		http.Error(w, "User is not Authorized", http.StatusUnauthorized)
	}
	var request struct {
		Table models.Table `json:"table"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userDatabase models.Database
	if err := ctrl.DB.Where(&models.Database{UserID: userID}).First(&userDatabase).Error; err != nil {
		http.Error(w, "Failed to find User Database", http.StatusInternalServerError)
		return
	}

	sqlDB, err := sql.Open("sqlite3", userDatabase.Path)
	if err != nil {
		http.Error(w, "Failed to open User Database", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	createTableSQL := generateTable(request.Table)
	if createTableSQL == "" {
		http.Error(w, "Failed to generate SQL", http.StatusInternalServerError)
		return
	}

	if _, err := sqlDB.Exec(createTableSQL); err != nil {
		http.Error(w, "Failed to create Table", http.StatusInternalServerError)
		return
	}
}

func (ctrl CustomDatabaseController) GetDatabase(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)

	if !ok {
		http.Error(w, "User is not Authorized", http.StatusUnauthorized)
	}

	var userDatabase models.Database
	if err := ctrl.DB.Where(&models.Database{UserID: userID}).First(&userDatabase).Error; err != nil {
		http.Error(w, "Failed to find User Database", http.StatusInternalServerError)
		return
	}

	sqlDB, err := sql.Open("sqlite3", userDatabase.Path)
	if err != nil {
		http.Error(w, "Failed to open User Database", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	rows, err := sqlDB.Query("SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		http.Error(w, "Failed to get table data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tables []models.Table
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			http.Error(w, "Failed to get table data", http.StatusInternalServerError)
			return
		}

		tableInfo, err := sqlDB.Query(fmt.Sprintf("PRAGMA table_info(%s);", tableName))
		if err != nil {
			http.Error(w, "Failed to fetch table info", http.StatusInternalServerError)
			return
		}
		defer tableInfo.Close()

		var table models.Table
		for tableInfo.Next() {
			var row models.Row
			var cid int
			var notNull int
			var dfltValue sql.NullString

			if err := tableInfo.Scan(&cid, &row.Name, &row.DataType, &notNull, &dfltValue, &notNull); err != nil {
				http.Error(w, "Failed to read table structure", http.StatusInternalServerError)
				return
			}
			row.Null = notNull == 0
			table.Rows = append(table.Rows, row)
		}
		tables = append(tables, table)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tables)
}

func (ctrl CustomDatabaseController) GetTable(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User is not authorized", http.StatusUnauthorized)
		return
	}

	var request struct {
		TableName string `json:"table_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	var userDatabase models.Database
	if err := ctrl.DB.Where(&models.Database{UserID: userID}).First(&userDatabase).Error; err != nil {
		http.Error(w, "Failed to find user database", http.StatusInternalServerError)
		return
	}

	sqlDB, err := sql.Open("sqlite3", userDatabase.Path)
	if err != nil {
		http.Error(w, "Failed to open user database", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	rows, err := sqlDB.Query(fmt.Sprintf("SELECT * FROM %s;", request.TableName))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get data from table %s", request.TableName), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		http.Error(w, "Failed to fetch column names", http.StatusInternalServerError)
		return
	}

	var results []map[string]interface{}
	for rows.Next() {
		columnValues := make([]interface{}, len(columns))
		columnPointers := make([]interface{}, len(columns))

		for i := range columnValues {
			columnPointers[i] = &columnValues[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}

		row := make(map[string]interface{})
		for i, colName := range columns {
			val := columnValues[i]
			if b, ok := val.([]byte); ok {
				row[colName] = string(b)
			} else {
				row[colName] = val
			}
		}
		results = append(results, row)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func generateTable(table models.Table) string {
	if table.Name == "" || len(table.Rows) == 0 {
		return ""
	}
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString(fmt.Sprintf("CREATE TABLE %s (", table.Name))
	sqlBuilder.WriteString("\nID INTEGER PRIMARY KEY NOT NULL,)")
	for i, row := range table.Rows {
		sqlBuilder.WriteString(fmt.Sprintf("%s %s", row.Name, row.DataType))
		if !row.Null {
			sqlBuilder.WriteString(fmt.Sprintf(" NOT NULL"))
		}
		if i < len(table.Rows)-1 {
			sqlBuilder.WriteString(",")
		}
	}
	sqlBuilder.WriteString(");")
	return sqlBuilder.String()
}
