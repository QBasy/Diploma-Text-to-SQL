package controllers

import (
	"auth-service/models"
	"auth-service/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

type AuthController struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{db: db}
}

type CreateDatabaseRequest struct {
	UserUUID string `json:"user_uuid"`
	Name     string `json:"name"`
}

type RegisterRequest struct {
	Username string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ChangePasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (uc *AuthController) Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Registering user: %s, %s", request.Username, request.Email)

	var existingUser models.User
	if err := uc.db.Where("username = ? OR email = ?", request.Username, request.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		UUID:         uuid.New().String(),
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: string(hashedPassword),
	}
	if err := uc.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	dbRequest := CreateDatabaseRequest{
		UserUUID: user.UUID,
		Name:     "default",
	}

	apiGatewayURL := utils.GetEnv("API_GATEWAY_URL", "http://localhost:5001")

	jsonBody, err := json.Marshal(dbRequest)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		uc.db.Delete(&user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare database creation request"})
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/database/create-database", apiGatewayURL), bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Request creation error: %v", err)
		uc.db.Delete(&user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request to database service"})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Service-Secret", utils.GetEnv("SECRET_KEY", ""))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Failed to create database: %v | Status: %v", err, resp.StatusCode)
		uc.db.Delete(&user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "uuid": user.UUID})
}

func (uc *AuthController) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := uc.db.Where("email = ?", request.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateJWT(user.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (uc *AuthController) ResetPassword(c *gin.Context) {
	var request ResetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := uc.db.Where("email = ?", request.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	token := uuid.New().String()
	expiry := time.Now().Add(1 * time.Hour)

	resetToken := models.PasswordResetToken{
		UserUUID: user.UUID,
		Token:    token,
		Expiry:   expiry,
	}
	if err := uc.db.Create(&resetToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reset token"})
		return
	}

	if err := utils.SendEmail(user.Email, "Password Reset", "Your reset token: "+token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset token sent to your email"})
}

func (uc *AuthController) ChangePassword(c *gin.Context) {
	var request ChangePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var resetToken models.PasswordResetToken
	if err := uc.db.Where("token = ?", request.Token).First(&resetToken).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired token"})
		return
	}

	if time.Now().After(resetToken.Expiry) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token has expired"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if err := uc.db.Model(&models.User{}).Where("id = ?", resetToken.UserUUID).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	uc.db.Delete(&resetToken)

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func (uc *AuthController) GetMe(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	userUUID, err := utils.ValidateJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var user models.User
	if err := uc.db.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uuid":  user.UUID,
		"name":  user.Username,
		"email": user.Email,
	})
}
