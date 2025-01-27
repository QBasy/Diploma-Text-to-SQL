package main

import (
	"database-service/models"
	"database-service/routes"
	"database-service/utils"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var database *gorm.DB

var rawDatabase *sql.DB

func initDB() {
	host := utils.GetEnv("POSTGRES_HOST", "")
	user := utils.GetEnv("POSTGRES_USER", "")
	password := utils.GetEnv("POSTGRES_PASSWORD", "")
	dbname := utils.GetEnv("POSTGRES_DB", "")
	port := utils.GetEnv("POSTGRES_PORT", "")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database using GORM: %v", err)
	}

	if err := gormDB.AutoMigrate(&models.User{}, models.UserDatabase{}); err != nil {
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

	r := routes.SetupRouter()

	log.Println("Starting server on :5002")
	port := os.Getenv("SERVER_PORT")
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
