package controllers

import (
	"database/sql"
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type CustomController struct {
	DB *gorm.DB
}

func (cc *CustomController) ExecuteSQLCustomQuery(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SQL string `json:"sql"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("postgres", "user=postgres dbname=test sslmode=disable")
	if err != nil {
		log.Printf("Error connecting to the database: %v", err)
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Выполнение SQL-запроса
	rows, err := db.Query(request.SQL)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Failed to execute SQL query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	columns, err := rows.Columns()
	if err != nil {
		log.Printf("Error fetching columns: %v", err)
		http.Error(w, "Failed to fetch columns", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		row := make(map[string]interface{})
		columnPointers := make([]interface{}, len(columns))

		for i := range columnPointers {
			var value interface{}
			columnPointers[i] = &value
		}

		if err := rows.Scan(columnPointers...); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}

		for i, colName := range columns {
			row[colName] = *(columnPointers[i].(*interface{}))
		}

		results = append(results, row)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (cc *CustomController) TestConnection(w http.ResponseWriter, r *http.Request) {
	result := "connected"
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
