package service

import (
	"github.com/QBasy/Diploma-text-to-SQL/open-api-service/internal/config"
	"net/http"
)

type AuthService interface {
	config.Config
}

func verifyToken(token string) bool {
	req, err := http.NewRequest("GET", "http://localhost:5001/auth/me", nil)
	if err != nil {
		return false
	}
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
