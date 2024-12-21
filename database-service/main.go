package main

import (
	"database-service/models"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

var database *gorm.DB

var rawDatabase *sql.DB

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

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database using GORM: %v", err)
	}

	if err := gormDB.AutoMigrate(&models.User{}, models.Database{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	database = gormDB

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database using sql.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	rawDatabase = sqlDB

	log.Println("Database connection established successfully")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}

	initDB()

	routes(database, rawDatabase)

	log.Println("Starting server on :5002")
	if err := http.ListenAndServe(":5002", nil); err != nil {
		log.Fatal(err)
	}
}
