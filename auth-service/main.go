package main

import (
	"auth-service/controllers"
	"auth-service/models"
	"auth-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

var db gorm.DB

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}
	dsn := "host=" + os.Getenv("POSTGRES_HOST") +
		" user=" + os.Getenv("POSTGRES_USER") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_DB") +
		" port=" + os.Getenv("POSTGRES_PORT") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("AuthService: Failed to connect Database %v", err)
	}
	_ = db.AutoMigrate(&models.User{}, &models.PasswordResetToken{})
}

func main() {
	r := gin.Default()

	userController := controllers.NewAuthController(db)

	auth := r.Group("/auth")
	{
		auth.POST("/register", userController.Register)
		auth.POST("/login", userController.Login)
		auth.POST("/reset-password", userController.ResetPassword)
		auth.POST("/change-password", userController.ChangePassword)

		authGoogle := auth.Group("/google")
		{
			authGoogle.GET("/", userController.LoginWithGoogle)
			authGoogle.GET("/callback", userController.GoogleCallback)
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Run(utils.GetEnv("PORT", ":5003"))
}
