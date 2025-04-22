package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text-to-sql/internal/model"
)

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ParseRequest(r *http.Request) (*model.Request, error) {
	var req model.Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return &req, fmt.Errorf("error parsing request: %v", err)
	}
	return &req, nil
}

func ParseResponse(w http.ResponseWriter, response model.Response) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	HandleError(err)
}
