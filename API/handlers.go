package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal user"})
		return
	}

	response, err := http.Post("http://localhost:5002/users", "application/json", bytes.NewBuffer(jsonUser))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to communicate with database service"})
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		c.JSON(response.StatusCode, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func registerHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal user"})
		return
	}

	response, err := http.Post("http://localhost:5002/users/create", "application/json", bytes.NewBuffer(jsonUser))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to communicate with database service"})
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		c.JSON(response.StatusCode, gin.H{"error": "User registration failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func queryToSQLHandler(c *gin.Context) {
	queryText := c.DefaultPostForm("query", "")
	if queryText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query text is required"})
		return
	}

	data := map[string]string{"text": queryText}
	jsonData, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare query data"})
		return
	}

	resp, err := http.Post("http://localhost:5003/convert", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to communicate with Text-to-SQL service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "Error processing query"})
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response from Text-to-SQL service"})
		return
	}

	// Отправка полученного SQL запроса в ответ
	c.JSON(http.StatusOK, result)
}
