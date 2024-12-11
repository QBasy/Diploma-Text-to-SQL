package handlers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handlers) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/auth/login", databaseServiceURL), "application/json", bytes.NewBuffer(toJSON(request)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Database service"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func (h *Handlers) Register(c *gin.Context) {
	var request map[string]interface{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	const maxRetries = 5
	var resp *http.Response
	var err error
	for i := 0; i < maxRetries; i++ {
		resp, err = http.Post(fmt.Sprintf("%s/auth/register", databaseServiceURL), "application/json", bytes.NewBuffer(toJSON(request)))
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Database service"})
		return
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}
