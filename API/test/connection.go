package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestApiService(t *testing.T) {
	// Ожидаем, что сервисы запустятся
	time.Sleep(5 * time.Second)

	// Проверяем доступность API сервиса
	resp, err := http.Get("http://localhost:8080/") // URL API сервиса
	if err != nil {
		t.Fatalf("Error connecting to API service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %v", resp.StatusCode)
	}

	dbRequest := map[string]string{"sql": "SELECT * FROM users"}
	reqBodyBytes, _ := json.Marshal(dbRequest)
	resp, err = http.Post("http://localhost:5002/executeSQL", "application/json", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		t.Fatalf("Error connecting to Database Service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %v", resp.StatusCode)
	}

	// Аналогичные тесты можно написать для text-to-sql-service
}
