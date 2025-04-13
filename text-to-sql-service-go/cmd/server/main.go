package main

import (
	"log"
	"net/http"
	Config "text-to-sql/internal/config"
	"text-to-sql/internal/handler"
)

func main() {
	config, err := Config.LoadConfig("config/local.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	log.Printf("Loaded config: PORT=%s, API_KEY=%s", config.TTSQL.PORT, config.TTSQL.APIKey)

	http.HandleFunc("/text-to-sql/gpt", handler.TextToSQLHandler)
	http.HandleFunc("/text-to-sql/groc", handler.TextToSQLHandlerWithGroc)

	log.Printf("Starting server on port %s...\n", config.TTSQL.PORT)
	err = http.ListenAndServe(":"+config.TTSQL.PORT, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
