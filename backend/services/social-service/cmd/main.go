package main

import (
	"log"
	"net"
	"os"
	"shared"
	pb "shared/pb/social"
	pbUser "shared/pb/user"
	"social/db"
	"social/handler"
	"social/internal/repository"
	"social/internal/service"

	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	startServer()
}

func startServer() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	if err = godotenv.Load(".env"); err != nil {
		logger.Warn("No .env file found, relying on environment variables", zap.Error(err))
	}

	logger.Info("Starting Social service server")

	typeServer := os.Getenv("TYPE")
	portServer := os.Getenv("PORT")

	list, err := net.Listen(typeServer, portServer)

	if err != nil {
		logger.Fatal("The server is not listening", zap.Error(err))
	}
	dbConn, queries, err := db.Conn()
	if err != nil {
		logger.Fatal("Error to connect with de DB", zap.Error(err))
	}
	defer dbConn.Close()

	socialRepo := repository.NewSocialRepository(queries)

	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	if userServiceAddr == "" {
		userServiceAddr = "localhost:8080"
	}

	clientTLS, err := shared.LoadClientTLSCredentials()
	if err != nil {
		logger.Fatal("failed to load client TLS credentials", zap.Error(err))
	}

	conn, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(clientTLS))
	if err != nil {
		log.Fatalf("Não foi possível conectar: %v", err)
	}
	defer conn.Close()

	userServiceClient := pbUser.NewUserServiceClient(conn)
	socialService := service.NewSocialService(socialRepo, userServiceClient)
	socialHandler := handler.NewSocialHandler(socialService, logger, userServiceClient)

	serverTLS, err := shared.LoadServerTLSCredentials()
	if err != nil {
		logger.Fatal("failed to load server TLS credentials", zap.Error(err))
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(serverTLS),
		grpc.UnaryInterceptor(
			grpcZap.UnaryServerInterceptor(logger),
		),
	)

	pb.RegisterSocialServiceServer(grpcServer, socialHandler)

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(list); err != nil {
		logger.Fatal("The server is not running", zap.Error(err))
	}
}
