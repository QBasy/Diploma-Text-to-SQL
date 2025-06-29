package main

import (
	"database-service/internal/routes"
	"database-service/pkg/utils"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var database *gorm.DB

var rawDatabase *sql.DB

func initDB() {
	host := utils.GetEnv("DB_HOST", "")
	user := utils.GetEnv("DB_USER", "")
	password := utils.GetEnv("DB_PASSWORD", "")
	dbname := utils.GetEnv("DB_NAME", "")
	port := utils.GetEnv("DB_PORT", "")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database using GORM: %v", err)
	}

	if err := utils.Migrate(gormDB); err != nil {
		log.Fatal("Migration failed: ", err)
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

	port := os.Getenv("SERVER_PORT")
	r.Run(":" + port)
}
