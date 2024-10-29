package main

import (
	"database-service/controllers"
	"gorm.io/gorm"
	"net/http"
)

func routes(db *gorm.DB) {
	itemController := &controllers.ItemController{DB: db}
	userController := &controllers.UserController{DB: db}

	http.HandleFunc("/items", itemController.GetItems)
	http.HandleFunc("/items/create", itemController.CreateItem)
	http.HandleFunc("/users", userController.GetUsers)
	http.HandleFunc("/users/create", userController.CreateUser)
}
