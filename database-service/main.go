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

func initDB() {
	host, username, password, dbName, port := os.Getenv("HOST"), os.Getenv("USERNAME"), os.Getenv("PASSWORD"), os.Getenv("DATABASE_NAME"), os.Getenv("PORT")

	username = "postgres"

	log.Printf("Host: %s, Username: %s, Password: %s, Database: %s, Port: %s\n", host, username, password, dbName, port)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, username, password, dbName, port)
	log.Println("DSN:", dsn)

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
