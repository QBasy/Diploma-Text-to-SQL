package main

import (
	"database-service/controllers"
	"database-service/middleware"
	"database/sql"
	"gorm.io/gorm"
	"net/http"
	"os"
)

func routes(db *gorm.DB, rawDB *sql.DB) {
	itemController := &controllers.ItemController{DB: db}
	userController := &controllers.UserController{DB: db}
	customController := &controllers.CustomController{DB: db, RawDatabase: rawDB}
	authController := &controllers.AuthController{DB: db}
	userDatabase := &controllers.CustomDatabaseController{DB: db, RawDatabase: rawDB}

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

	http.HandleFunc("/db/create-table", userDatabase.CreateTable)
	http.HandleFunc("/db/execute-query", userDatabase.ExecuteQuery)
	http.HandleFunc("/db/get-database", userDatabase.GetCustomDatabase)

	http.Handle("/", middleware.LogRequest(mux))
}
