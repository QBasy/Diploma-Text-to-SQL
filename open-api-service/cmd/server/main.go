package main

import (
	Config "github.com/QBasy/Diploma-text-to-SQL/open-api-service/internal/config"
	"github.com/QBasy/Diploma-text-to-SQL/open-api-service/internal/router"
	"log"
	"net/http"
)

func main() {
	cfg, err := Config.GetConfig("config/local.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	r := router.New(cfg)

	log.Printf("Server started on port %s\n", cfg.PORT)
	if err := http.ListenAndServe(":"+cfg.PORT, r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
