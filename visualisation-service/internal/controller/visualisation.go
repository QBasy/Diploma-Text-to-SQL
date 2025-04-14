package controller

import (
	"context"
	pb "visualisation-service/generated/visualisationpb"
	"visualisation-service/internal/service"
)

type Server struct {
	pb.UnimplementedVisualisationServiceServer
}

func (s *Server) GenerateChart(ctx context.Context, req *pb.QueryResult) (*pb.SVGResponse, error) {
	svgBuffer, err := service.GenerateChart(req)
	if err != nil {
		return nil, err
	}
	return &pb.SVGResponse{Svg: svgBuffer.String()}, nil
}
