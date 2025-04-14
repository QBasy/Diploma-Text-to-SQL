package main

import (
	"API/middleware"
	"API/utils"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

var (
	authServiceURL      = utils.GetEnv("AUTH_SERVICE_URL", "http://localhost:5003")
	databaseServiceURL  = utils.GetEnv("DATABASE_SERVICE_URL", "http://localhost:5004/api")
	textToSQLServiceURL = utils.GetEnv("TEXT_TO_SQL_SERVICE_URL", "http://127.0.0.1:5006")
	metadataServiceURL  = utils.GetEnv("METADATA_SERVICE_URL", "http://localhost:5005/api/metadata")
	historyServiceURL   = utils.GetEnv("HISTORY_SERVICE_URL", "http://localhost:5008/api/history")
)

func proxyHandler(serviceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetURL := serviceURL + c.FullPath()

		req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for k, v := range c.Request.Header {
			for _, vv := range v {
				req.Header.Add(k, vv)
			}
		}

		q := req.URL.Query()
		for k, v := range c.Request.URL.Query() {
			for _, vv := range v {
				q.Add(k, vv)
			}
		}
		req.URL.RawQuery = q.Encode()

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		for k, v := range resp.Header {
			for _, vv := range v {
				c.Writer.Header().Add(k, vv)
			}
		}

		c.Status(resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		_, err = c.Writer.Write(body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

func historyRoutes(r *gin.Engine) {
	historyGroup := r.Group("/history")
	historyGroup.Use(middleware.RateLimiter())
	{
		historyGroup.GET("", proxyHandler(historyServiceURL))
		historyGroup.POST("", proxyHandler(historyServiceURL))
		historyGroup.DELETE("", proxyHandler(historyServiceURL))
	}
}

func authRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	authGroup.Use(middleware.RateLimiter())
	{
		authGroup.POST("/register", proxyHandler(authServiceURL))
		authGroup.POST("/login", proxyHandler(authServiceURL))
		authGroup.POST("/reset-password", proxyHandler(authServiceURL))
		authGroup.POST("/change-password", proxyHandler(authServiceURL))
		authGroup.GET("/me", proxyHandler(authServiceURL))
		authGroup.GET("/google", proxyHandler(authServiceURL))
		authGroup.GET("/google/callback", proxyHandler(authServiceURL))
	}
}

func databaseRoutes(r *gin.Engine) {
	dbGroup := r.Group("/api/database")
	dbGroup.Use(middleware.RateLimiter())
	{
		dbGroup.POST("/create-database", proxyHandler(databaseServiceURL))
		dbGroup.POST("/execute-sql", proxyHandler(databaseServiceURL))
		dbGroup.GET("/schema", proxyHandler(databaseServiceURL))
		dbGroup.GET("/schema-complex", proxyHandler(databaseServiceURL))
		dbGroup.POST("/tables", proxyHandler(databaseServiceURL))
	}
}

func textToSQLRoutes(r *gin.Engine) {
	textToSQLGroup := r.Group("/text-to-sql")
	textToSQLGroup.Use(middleware.RateLimiter())
	{
		textToSQLGroup.POST("/gpt", proxyHandler(textToSQLServiceURL))
		textToSQLGroup.POST("/simple", proxyHandler(textToSQLServiceURL))
		textToSQLGroup.POST("/complex", proxyHandler(textToSQLServiceURL))
	}
}

func metadataRoutes(r *gin.Engine) {
	metadataGroup := r.Group("/metadata")
	metadataGroup.Use(middleware.RateLimiter())
	{
		metadataGroup.GET("/:user_id/:database_uuid", proxyHandler(metadataServiceURL))
		metadataGroup.POST("/:user_id/:database_uuid", proxyHandler(metadataServiceURL))
		metadataGroup.PUT("/:user_id/:database_uuid", proxyHandler(metadataServiceURL))
		metadataGroup.DELETE("/:user_id/:database_uuid", proxyHandler(metadataServiceURL))
	}
}

func handleHealthCheck(c *gin.Context) {
	services := map[string]string{
		"auth":     authServiceURL + "/health",
		"database": databaseServiceURL + "/health",
		"metadata": metadataServiceURL + "/health",
		"text2sql": textToSQLServiceURL + "/health",
		"history":  historyServiceURL + "/health",
	}

	status := make(map[string]string)

	for service, url := range services {
		resp, err := http.Get(url)
		if err != nil {
			status[service] = "down"
		} else {
			body, _ := io.ReadAll(resp.Body)
			status[service] = string(body)
			resp.Body.Close()
		}
	}

	c.JSON(http.StatusOK, status)
}
