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

func initDB() *gorm.DB {
	host, user, password, dbName, port := os.Getenv("HOST"), os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DATABASE_NAME"), os.Getenv("PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)
	if db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		panic("failed to connect database")
	} else {
		_ = db.AutoMigrate(&models.User{}, &models.Item{})
		return db
	}
}

func main() {
	var db *gorm.DB
	db = initDB()

	routes(db)

	log.Println("Starting server on :5002")
	if err := http.ListenAndServe(":5002", nil); err != nil {
		log.Fatal(err)
	}
}
