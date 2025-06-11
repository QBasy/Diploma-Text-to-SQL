package controller

import (
	"context"
	pb "visualisation-service/generated/visualisationpb"
	"visualisation-service/internal/service"
	"visualisation-service/pkg/logger"
)

type Server struct {
	pb.UnimplementedVisualisationServiceServer
}

func (s *Server) GenerateChart(ctx context.Context, req *pb.QueryResult) (*pb.SVGResponse, error) {
	logger.InfoLogger.Printf("Запрос на генерацию графика получен")
	svgBuffer, err := service.GenerateChart(req)
	if err != nil {
		logger.ErrorLogger.Printf("Ошибка генерации графика: %v", err)
		return nil, err
	}
	return &pb.SVGResponse{Svg: svgBuffer.String()}, nil
}
