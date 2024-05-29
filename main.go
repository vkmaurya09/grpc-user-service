package main

import (
	"log"
	"net"

	"github.com/grpc-user-service/proto"
	"github.com/grpc-user-service/service"
	"google.golang.org/grpc"
)

func main() {
	userService := service.NewUserService()

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, userService)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("gRPC server running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
