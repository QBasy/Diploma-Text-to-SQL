package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	pb "visualisation-service/generated/visualisationpb"
	"visualisation-service/internal/controller"
)

func main() {
	lis, err := net.Listen("tcp", ":5007")
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	Server := &controller.Server{}

	s := grpc.NewServer()
	pb.RegisterVisualisationServiceServer(s, Server)

	log.Println("gRPC-сервер Visualisation запущен на :5007")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Ошибка работы сервера: %v", err)
	}
}
