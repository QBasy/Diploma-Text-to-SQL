package service

import (
	"context"
	"github.com/QBasy/Diploma-text-to-SQL/open-api-service/internal/config"
	pb "github.com/QBasy/Diploma-text-to-SQL/open-api-service/internal/proto/generated/visualisationpb"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

type VisualService struct {
	cfg *config.Config
}

func NewVisualService(cfg *config.Config) *VisualService {
	return &VisualService{cfg: cfg}
}

func (s *VisualService) GenerateSVG(w http.ResponseWriter, qr *pb.QueryResult) {
	conn, err := grpc.Dial(s.cfg.VisualisationService, grpc.WithInsecure())
	if err != nil {
		http.Error(w, "gRPC connection error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pb.NewVisualisationServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GenerateChart(ctx, qr)
	if err != nil {
		http.Error(w, "Chart generation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write([]byte(resp.Svg))
}
