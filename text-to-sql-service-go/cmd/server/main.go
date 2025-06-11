package main

import (
	"net/http"
	Config "text-to-sql/internal/config"
	"text-to-sql/internal/handler"
	"text-to-sql/internal/middleware"
	"text-to-sql/pkg/logger"
)

func main() {
	logger.Init()

	config, err := Config.LoadConfig("config/local.yaml")
	if err != nil {
		logger.ErrorLogger.Fatalf("Error loading config: %v", err)
	}
	logger.InfoLogger.Printf("Loaded config: PORT=%s, API_KEY=%s", config.TTSQL.PORT, config.TTSQL.APIKey)

	mux := http.NewServeMux()

	mux.HandleFunc("/text-to-sql/gpt", handler.TextToSQLHandler)
	mux.HandleFunc("/text-to-sql/groc", handler.TextToSQLHandlerWithGroc)
	mux.HandleFunc("/text-to-sql/simple", handler.TextToSQLHandlerWithGroc)
	mux.HandleFunc("/text-to-sql/complex", handler.TextToSQLHandlerWithGroc)

	loggedMux := middleware.LoggingMiddleware(mux)

	logger.InfoLogger.Printf("Starting server on port %s...\n", config.TTSQL.PORT)
	err = http.ListenAndServe(":"+config.TTSQL.PORT, loggedMux)
	if err != nil {
		logger.ErrorLogger.Fatalf("Error starting server: %v", err)
	}
}
