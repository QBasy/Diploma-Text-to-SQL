package handler

import (
	"github.com/QBasy/Diploma-text-to-SQL/open-api-service/internal/config"
	"github.com/QBasy/Diploma-text-to-SQL/open-api-service/internal/service"
	"net/http"
)

type TextSQLHandler struct {
	service *service.TextSQLService
}

func NewTextSQLHandler(cfg *config.Config) *TextSQLHandler {
	return &TextSQLHandler{
		service: service.NewTextSQLService(cfg),
	}
}

func (h *TextSQLHandler) HandleSimple(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	h.service.ProxyRequest("simple", w, r)
}

func (h *TextSQLHandler) HandleComplex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	h.service.ProxyRequest("complex", w, r)
}
