package main

import (
	"auth-service/controllers"
	"auth-service/models"
	"auth-service/utils"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db := setupDatabase()

	r := gin.Default()
	userController := controllers.NewAuthController(db)

	auth := r.Group("/auth")
	{
		auth.POST("/register", userController.Register)
		auth.POST("/login", userController.Login)
		auth.POST("/reset-password", userController.ResetPassword)
		auth.POST("/change-password", userController.ChangePassword)
		auth.GET("/me", userController.GetMe)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "5003"
	}
	r.Run(":" + port)
}

func setupDatabase() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		utils.GetEnv("DB_HOST", "localhost"),
		utils.GetEnv("DB_USER", "postgres"),
		utils.GetEnv("DB_PASSWORD", ""),
		utils.GetEnv("DB_NAME", "mvpdiploma"),
		utils.GetEnv("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("AuthService: Failed to connect Database %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.UserDatabase{}, &models.PasswordResetToken{}, &models.EmailVerificationToken{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	return db
}
