package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userUUID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_uuid": userUUID,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user_uuid"].(string), nil
	}
	return "", jwt.ErrInvalidKey
}

func GetUserUUIDFromRequest(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
		return "", errors.New("authorization token missing or invalid")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	userUUID, err := ValidateJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return "", err
	}

	return userUUID, nil
}
