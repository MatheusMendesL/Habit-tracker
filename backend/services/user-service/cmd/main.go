package main

import (
	"log"
	"net"
	pb "shared/pb/user"
	"user-service/handler"
	"user-service/internal/repository"
	"user-service/internal/service"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting server...")
	list, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, userHandler)

	if err := grpcServer.Serve(list); err != nil {
		log.Fatal(err)
	}
}
