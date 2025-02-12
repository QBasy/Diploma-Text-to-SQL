package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

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
