package main

import (
	"database-service/controllers"
	"database/sql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

func routes(db *gorm.DB, rawDB *sql.DB) {
	userController := &controllers.UserController{DB: db}
	customController := &controllers.CustomController{DB: db, RawDatabase: rawDB}
	authController := &controllers.AuthController{DB: db}
	userDatabase := &controllers.CustomDatabaseController{DB: db, RawDatabase: rawDB}
	textToSQLController := &controllers.TextToSQLController{DB: db, RawDatabase: rawDB}

	authController.JwtSecret = []byte(os.Getenv("SECRET_KEY"))

	mux := http.NewServeMux()

	http.HandleFunc("/auth/login", authController.Login)
	http.HandleFunc("/auth/register", authController.Register)

	http.HandleFunc("/executeSQL", customController.TestConnection)
	http.HandleFunc("/execute-query", customController.ExecuteSQLCustomQuery)

	http.HandleFunc("/users", userController.GetUsers)
	http.HandleFunc("/users/create", userController.CreateUser)

	http.HandleFunc("/db/create-table", userDatabase.CreateTable)
	http.HandleFunc("/db/execute-query", textToSQLController.ExecuteQuery)
	http.HandleFunc("/db/get-database", userDatabase.GetDatabase)

	http.Handle("/", LogRequest(mux))
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rr := &ResponseRecorder{ResponseWriter: w, StatusCode: http.StatusOK}

		next.ServeHTTP(rr, r)

		log.Printf(
			"Method: %s, URL: %s, Status Code: %d, Duration: %s\n",
			r.Method,
			r.URL.Path,
			rr.StatusCode,
			time.Since(start),
		)
	})
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
}

func (rr *ResponseRecorder) WriteHeader(statusCode int) {
	rr.StatusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}
