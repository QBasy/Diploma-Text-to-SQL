package controllers

import (
	"bytes"
	"database-service/models"
	"database/sql"
	"encoding/json"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
)

type TextToSQLController struct {
	DB          *gorm.DB
	RawDatabase *sql.DB
}

type SQLResult struct {
	SqlQuery      string  `json:"sql_query"`
	Confidence    float64 `json:"confidence,omitempty"`
	ExecutionTime float64 `json:"execution_time"`
}

type QueryInput struct {
	Text string `json:"text"`
}

func (ctrl TextToSQLController) ExecuteQuery(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User is not authorized", http.StatusUnauthorized)
		return
	}

	var request struct {
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	var user models.User
	if err := ctrl.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var userDatabase models.Database
	if err := ctrl.DB.Where("user_id = ? AND file_path = ?", userID, user.Databases).First(&userDatabase).Error; err != nil {
		http.Error(w, "Database not found or not associated with the user", http.StatusForbidden)
		return
	}

	sqlDB, err := sql.Open("sqlite3", userDatabase.Path)
	if err != nil {
		http.Error(w, "Failed to open user database", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	rows, err := sqlDB.Query(request.Query)
	if err != nil {
		http.Error(w, "Failed to execute query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		http.Error(w, "Failed to fetch columns", http.StatusInternalServerError)
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
			row[colName] = val
		}
		results = append(results, row)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(results)
}

func (ctrl TextToSQLController) ConvertTextToSQLAndExecute(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User is not authorized", http.StatusUnauthorized)
		return
	}

	request := &QueryInput{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	sqlResult, err := sendToTextToSql(request.Text)
	if err != nil {
		http.Error(w, "Failed to generate SQL", http.StatusInternalServerError)
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

	rows, err := sqlDB.Query(sqlResult.SqlQuery)
	if err != nil {
		http.Error(w, "Failed to execute SQL query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		http.Error(w, "Failed to fetch columns", http.StatusInternalServerError)
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
			row[colName] = val
		}
		results = append(results, row)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(results)
}

func sendToTextToSql(text string) (*SQLResult, error) {
	url := "http://localhost:5003/convert"
	queryInput := QueryInput{Text: text}
	jsonData, _ := json.Marshal(queryInput)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var result SQLResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
