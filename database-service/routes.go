package main

import (
	"database-service/controllers"
	"database-service/middleware"
	"gorm.io/gorm"
	"net/http"
	"os"
)

func routes(db *gorm.DB) {
	itemController := &controllers.ItemController{DB: db}
	userController := &controllers.UserController{DB: db}
	customController := &controllers.CustomController{DB: db}
	authController := &controllers.AuthController{DB: db}

	authController.JwtSecret = []byte(os.Getenv("SECRET_KEY"))

	mux := http.NewServeMux()

	http.HandleFunc("/auth/login", authController.Login)
	http.HandleFunc("/auth/register", authController.Register)

	http.HandleFunc("/executeSQL", customController.TestConnection)
	http.HandleFunc("/execute-query", customController.ExecuteSQLCustomQuery)

	http.HandleFunc("/items", itemController.GetItems)
	http.HandleFunc("/items/create", itemController.CreateItem)

	http.HandleFunc("/users", userController.GetUsers)
	http.HandleFunc("/users/create", userController.CreateUser)

	http.Handle("/", middleware.LogRequest(mux))
}
