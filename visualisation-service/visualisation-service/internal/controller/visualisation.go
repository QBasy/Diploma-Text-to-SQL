package controller

import (
	"context"
	"log"
	pb "visualisation-service/generated/visualisationpb"
	"visualisation-service/internal/service"
)

type Server struct {
	pb.UnimplementedVisualisationServiceServer
}

func (s *Server) GenerateChart(ctx context.Context, req *pb.QueryResult) (*pb.SVGResponse, error) {
	log.Printf("1")
	svgBuffer, err := service.GenerateChart(req)
	if err != nil {
		log.Printf("error, %v", err)
		return nil, err
	}
	return &pb.SVGResponse{Svg: svgBuffer.String()}, nil
}
