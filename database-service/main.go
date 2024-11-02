package main

import (
	"database-service/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

var database *gorm.DB

func initDB() {
	host, user, password, dbName, port := os.Getenv("HOST"), os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DATABASE_NAME"), os.Getenv("PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)

	if db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		panic("failed to connect database")
	} else {
		_ = db.AutoMigrate(&models.User{}, &models.Item{})
		database = db
	}
}

func main() {
	initDB()

	routes(database)

	log.Println("Starting server on :5002")
	if err := http.ListenAndServe(":5002", nil); err != nil {
		log.Fatal(err)
	}
}
