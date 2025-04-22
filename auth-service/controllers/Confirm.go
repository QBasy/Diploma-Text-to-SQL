package controllers

import (
	"auth-service/models"
	"auth-service/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

func CleanExpiredUsers(db *gorm.DB, interval time.Duration) {
	for {
		time.Sleep(interval)

		var expiredTokens []models.EmailVerificationToken
		now := time.Now()

		if err := db.Where("expiry < ?", now).Find(&expiredTokens).Error; err != nil {
			log.Printf("Failed to fetch expired tokens: %v", err)
			continue
		}

		for _, token := range expiredTokens {
			var user models.User
			if err := db.First(&user, "uuid = ?", token.UserUUID).Error; err == nil {
				log.Printf("Deleting unconfirmed user: %s", user.Email)
				db.Delete(&user)
			}
			db.Delete(&token)
		}
	}
}

func (uc *AuthController) ConfirmEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	var vToken models.EmailVerificationToken
	if err := uc.db.Where("token = ?", token).First(&vToken).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	if time.Now().After(vToken.Expiry) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token expired"})
		return
	}

	var user models.User
	if err := uc.db.First(&user, "uuid = ?", vToken.UserUUID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	uc.db.Delete(&vToken)

	dbRequest := CreateDatabaseRequest{
		UserUUID: user.UUID,
		Name:     "default",
	}

	jsonBody, err := json.Marshal(dbRequest)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare database creation request"})
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/database/create-database", utils.GetEnv("API_GATEWAY_URL", "")), bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Request creation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request to database service"})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Service-Secret", utils.GetEnv("SECRET_KEY", ""))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Failed to create user DB: %v | Status: %v", err, resp.StatusCode)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email confirmed and database created"})
}
