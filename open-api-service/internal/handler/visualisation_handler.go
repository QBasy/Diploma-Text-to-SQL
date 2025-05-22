package handler

import (
	"encoding/json"
	"github.com/QBasy/Diploma-text-to-SQL/open-api-service/internal/config"
	pb "github.com/QBasy/Diploma-text-to-SQL/open-api-service/internal/proto/generated/visualisationpb"
	"github.com/QBasy/Diploma-text-to-SQL/open-api-service/internal/service"
	"net/http"
)

type VisualHandler struct {
	service *service.VisualService
}

func NewVisualHandler(cfg *config.Config) *VisualHandler {
	return &VisualHandler{
		service: service.NewVisualService(cfg),
	}
}

func (h *VisualHandler) HandleVisual(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var qr pb.QueryResult
	if err := json.NewDecoder(r.Body).Decode(&qr); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	h.service.GenerateSVG(w, &qr)
}
