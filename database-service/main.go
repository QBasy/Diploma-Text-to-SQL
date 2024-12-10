package main

import (
	"database-service/models"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

var database *gorm.DB

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		if fallback == "" {
			log.Fatalf("Environment variable %s is not set", key)
		}
		return fallback
	}
	return value
}

func initDB() {
	host := getEnv("DB_HOST", "")
	user := getEnv("DB_USERNAME", "")
	password := getEnv("DB_PASSWORD", "")
	dbname := getEnv("DB_NAME", "")
	port := getEnv("DB_PORT", "")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Item{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	database = db
	log.Println("Database connection established successfully")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}

	initDB()

	routes(database)

	log.Println("Starting server on :5002")
	if err := http.ListenAndServe(":5002", nil); err != nil {
		log.Fatal(err)
	}
}
