package main

import (
	"google.golang.org/grpc"
	"net"
	pb "visualisation-service/generated/visualisationpb"
	"visualisation-service/internal/controller"
	"visualisation-service/pkg/logger"
)

func main() {
	logger.Init()

	lis, err := net.Listen("tcp", ":5007")
	if err != nil {
		logger.ErrorLogger.Fatalf("Ошибка запуска сервера: %v", err)
	}

	Server := &controller.Server{}

	s := grpc.NewServer()
	pb.RegisterVisualisationServiceServer(s, Server)

	logger.InfoLogger.Println("gRPC-сервер Visualisation запущен на :5007")
	if err := s.Serve(lis); err != nil {
		logger.ErrorLogger.Fatalf("Ошибка работы сервера: %v", err)
	}
}
