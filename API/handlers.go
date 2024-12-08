package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	databaseServiceURL  = "http://localhost:5002"
	textToSQLServiceURL = "http://localhost:5003"
)

func convertToSQLHandler(c *gin.Context) {
	var request struct {
		Text string `json:"text"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/convert", textToSQLServiceURL), "application/json", bytes.NewBufferString(`{"text":"`+request.Text+`"}`))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Text-to-SQL service"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func textToSQLHealthHandler(c *gin.Context) {
	resp, err := http.Get(fmt.Sprintf("%s/health", textToSQLServiceURL))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Text-to-SQL service"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func executeCustomSQLHandler(c *gin.Context) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/execute-query", databaseServiceURL), "application/json", bytes.NewBuffer(toJSON(request)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Database service"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func getItemsHandler(c *gin.Context) {
	resp, err := http.Get(fmt.Sprintf("%s/items", databaseServiceURL))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Database service"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func createItemHandler(c *gin.Context) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/items/create", databaseServiceURL), "application/json", bytes.NewBuffer(toJSON(request)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Database service"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func getUsersHandler(c *gin.Context) {
	resp, err := http.Get(fmt.Sprintf("%s/users", databaseServiceURL))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Database service"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func createUserHandler(c *gin.Context) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/users/create", databaseServiceURL), "application/json", bytes.NewBuffer(toJSON(request)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Database service"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func toJSON(data interface{}) []byte {
	jsonData, _ := json.Marshal(data)
	return jsonData
}
