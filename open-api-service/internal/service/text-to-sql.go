package service

import (
	"fmt"
	"github.com/QBasy/Diploma-text-to-SQL/open-api-service/internal/config"
	"io"
	"net/http"
)

type TextSQLService struct {
	cfg *config.Config
}

func NewTextSQLService(cfg *config.Config) *TextSQLService {
	return &TextSQLService{cfg: cfg}
}

func (s *TextSQLService) ProxyRequest(endpoint string, w http.ResponseWriter, r *http.Request) {
	fullURL := fmt.Sprintf("http://%s/text-to-sql/%s", s.cfg.ApiGatewayURL, endpoint)

	req, err := http.NewRequest("POST", fullURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header = r.Header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
