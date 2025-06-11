package middleware

import (
	"database-service/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func VerifyAuthService() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("X-Auth-Service-Secret")

		if authHeader != utils.GetEnv("SECRET_KEY", "") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
			c.Abort()
			return
		}

		c.Next()
	}
}
