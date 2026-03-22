package main

import (
	"log"
	"net"
	pb "shared/pb/user"
	"user-service/db"
	"user-service/handler"
	"user-service/internal/repository"
	"user-service/internal/service"

	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting server...")

	if err := startServer(); err != nil {
		panic(err)
	}
}

func startServer() (err error) {
	list, err := net.Listen("tcp", ":8080")

	if err != nil {
		return err
	}
	dbConn, queries, err := db.Conn()
	if err != nil {
		return err
	}
	defer dbConn.Close()

	userRepo := repository.NewUserRepository(queries)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcZap.UnaryServerInterceptor(logger),
		),
	)

	pb.RegisterUserServiceServer(grpcServer, userHandler)

	if err := grpcServer.Serve(list); err != nil {
		return err
	}

	return nil
}
