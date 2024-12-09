package test

import (
	"API/handlers"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	textToSQLServiceURL = "http://localhost:5003"
	databaseServiceURL  = "http://localhost:5002"
	handler             = handlers.Handlers{}
)

// Mock endpoints for external services
func mockTextToSQLService() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/convert" && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"query": "SELECT * FROM users;"}`))
		} else if r.URL.Path == "/health" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func mockDatabaseService() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/items" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id": 1, "name": "Item1"}]`))
		} else if r.URL.Path == "/users/create" && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id": 1, "name": "John Doe"}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestConvertToSQLHandler(t *testing.T) {
	textToSQLServer := mockTextToSQLService()
	defer textToSQLServer.Close()

	textToSQLServiceURL = textToSQLServer.URL

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/text-to-sql/convert", handler.ConvertToSQLHandler)

	requestBody := `{"text": "Get all users"}`
	req, _ := http.NewRequest(http.MethodPost, "/text-to-sql/convert", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedBody := `{"query": "SELECT * FROM users;"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestGetItemsHandler(t *testing.T) {
	databaseServer := mockDatabaseService()
	defer databaseServer.Close()

	// Replace database service URL with mock URL
	databaseServiceURL = databaseServer.URL

	// Set up Gin engine with the route
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/database/items", handler.GetItemsHandler)

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/database/items", nil)

	// Perform the test request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	expectedBody := `[{"id": 1, "name": "Item1"}]`
	assert.JSONEq(t, expectedBody, w.Body.String())
}
