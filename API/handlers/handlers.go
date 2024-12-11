package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
}

const (
	databaseServiceURL  = "http://database-service-container:5002"
	textToSQLServiceURL = "http://text_to_sql_service-container:5003"
)

func (h *Handlers) ConvertToSQLHandler(c *gin.Context) {
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

func (h *Handlers) TextToSQLHealthHandler(c *gin.Context) {
	resp, err := http.Get(fmt.Sprintf("%s/health", textToSQLServiceURL))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Text-to-SQL service"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func (h *Handlers) ExecuteCustomTextToSQLHandler(c *gin.Context) {
	var request struct {
		Text string `json:"text"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	textToSQLResp, err := http.Post(fmt.Sprintf("%s/convert", textToSQLServiceURL), "application/json", bytes.NewBufferString(`{"text":"`+request.Text+`"}`))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Text-to-SQL service"})
		return
	}
	defer textToSQLResp.Body.Close()

	if textToSQLResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(textToSQLResp.Body)
		c.JSON(textToSQLResp.StatusCode, gin.H{"error": "Text-to-SQL service error", "details": string(body)})
		return
	}

	var sqlResponse struct {
		Query string `json:"query"`
	}
	if err := json.NewDecoder(textToSQLResp.Body).Decode(&sqlResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode Text-to-SQL response"})
		return
	}

	dbResp, err := http.Post(fmt.Sprintf("%s/execute-query", databaseServiceURL), "application/json", bytes.NewBuffer(toJSON(sqlResponse)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Database service"})
		return
	}
	defer dbResp.Body.Close()

	body, _ := io.ReadAll(dbResp.Body)
	c.Data(dbResp.StatusCode, "application/json", body)
}

func (h *Handlers) ExecuteCustomSQLHandler(c *gin.Context) {
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

func (h *Handlers) GetItemsHandler(c *gin.Context) {
	resp, err := http.Get(fmt.Sprintf("%s/items", databaseServiceURL))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Database service"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func (h *Handlers) CreateItemHandler(c *gin.Context) {
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

func (h *Handlers) GetUsersHandler(c *gin.Context) {
	resp, err := http.Get(fmt.Sprintf("%s/users", databaseServiceURL))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Database service"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func (h *Handlers) CreateUserHandler(c *gin.Context) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println(err)
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
