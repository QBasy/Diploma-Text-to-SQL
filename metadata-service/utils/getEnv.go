package utils

import (
	"log"
	"os"
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		if fallback == "" {
			log.Fatalf("Environment variable %s is not set", key)
		}
		return fallback
	}
	return value
}