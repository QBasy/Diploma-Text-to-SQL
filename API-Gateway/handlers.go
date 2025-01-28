package main

import (
	"API/middleware"
	"API/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
)

var (
	authServiceURL      = utils.GetEnv("AUTH_SERVICE_URL", "http://localhost:5003")
	databaseServiceURL  = utils.GetEnv("DATABASE_SERVICE_URL", "http://localhost:5004")
	textToSQLServiceURL = utils.GetEnv("TEXT_TO_SQL_SERVICE_URL", "http://localhost:5006")
	metadataServiceURL  = utils.GetEnv("METADATA_SERVICE_URL", "http://localhost:5005")
)

func authRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	authGroup.Use(middleware.RateLimiter())
	{
		authGroup.POST("/register", handleAuthRegister)
		authGroup.POST("/login", handleAuthLogin)
		authGroup.POST("/reset-password", handleAuthResetPassword)
		authGroup.POST("/change-password", handleAuthChangePassword)

		authGroup.GET("/google", handleAuthGoogleLogin)
		authGroup.GET("/google/callback", handleAuthGoogleCallback)
	}
}

func databaseRoutes(r *gin.Engine) {
	dbGroup := r.Group("/api/database")
	dbGroup.Use(middleware.RateLimiter())
	{
		dbGroup.POST("/create-database", handleCreateDatabase)
		dbGroup.POST("/execute-sql", handleExecuteSQL)
		dbGroup.GET("/schema", handleGetSchema)
		dbGroup.POST("/tables", handleCreateTable)
	}
}

func textToSQLRoutes(r *gin.Engine) {
	textToSQLGroup := r.Group("/api/text-to-sql")
	textToSQLGroup.Use(middleware.RateLimiter())
	{
		textToSQLGroup.POST("/simple", handleTextToSQLSimple)
		textToSQLGroup.POST("/complex", handleTextToSQLComplex)
	}
}

func metadataRoutes(r *gin.Engine) {
	metadataGroup := r.Group("/api/metadata")
	metadataGroup.Use(middleware.RateLimiter())
	{
		metadataGroup.GET("/user", handleGetUserMetadata)
		metadataGroup.POST("/user", handleUpdateUserMetadata)
		metadataGroup.GET("/database", handleGetDatabaseMetadata)
		metadataGroup.POST("/database", handleUpdateDatabaseMetadata)
	}
}

func handleAuthRegister(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().
		SetBody(c.Request.Body).
		Post(authServiceURL + "/auth/register")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Body())
}

func handleAuthLogin(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().
		SetBody(c.Request.Body).
		Post(authServiceURL + "/auth/login")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Body())
}

func handleCreateDatabase(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().
		SetBody(c.Request.Body).
		Post(databaseServiceURL + "/api/create-database")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Body())
}

func handleExecuteSQL(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().
		SetBody(c.Request.Body).
		Post(databaseServiceURL + "/api/execute-sql")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Body())
}

func handleTextToSQLSimple(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().
		SetBody(c.Request.Body).
		Post(textToSQLServiceURL + "/text-to-sql/simple")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Body())
}

func handleTextToSQLComplex(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().
		SetBody(c.Request.Body).
		Post(textToSQLServiceURL + "/text-to-sql/complex")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Body())
}

func handleGetUserMetadata(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().
		Get(metadataServiceURL + "/api/metadata/user")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Body())
}

func handleAuthResetPassword(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().
		SetBody(c.Request.Body).
		Post(authServiceURL + "/auth/reset-password")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Body())
}

func handleAuthChangePassword(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().
		SetBody(c.Request.Body).
		Post(authServiceURL + "/auth/change-password")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Body())
}

func handleAuthGoogleLogin(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().
		Get(authServiceURL + "/auth/google")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Body())
}

func handleAuthGoogleCallback(c *gin.Context) {
	client := resty.New()

	queryParams := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0]
		}
	}

	resp, err := client.R().
		SetQueryParams(queryParams).
		Get(authServiceURL + "/auth/google/callback")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Body())
}

func handleGetSchema(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().
		Get(databaseServiceURL + "/api/schema")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Body())
}

func handleCreateTable(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().
		SetBody(c.Request.Body).
		Post(databaseServiceURL + "/api/tables")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Body())
}

func handleUpdateUserMetadata(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().
		SetBody(c.Request.Body).
		Post(metadataServiceURL + "/api/metadata/user")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Body())
}

func handleGetDatabaseMetadata(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().
		Get(metadataServiceURL + "/api/metadata/database")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Body())
}

func handleUpdateDatabaseMetadata(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().
		SetBody(c.Request.Body).
		Post(metadataServiceURL + "/api/metadata/database")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Body())
}

func handleHealthCheck(c *gin.Context) {
	services := []string{
		authServiceURL + "/health",
		databaseServiceURL + "/health",
		metadataServiceURL + "/health",
		textToSQLServiceURL + "/health",
	}

	client := resty.New()
	status := make(map[string]string)

	for _, service := range services {
		resp, err := client.R().Get(service)
		if err != nil {
			status[service] = "down"
		} else {
			status[service] = resp.String()
		}
	}

	c.JSON(http.StatusOK, status)
}
