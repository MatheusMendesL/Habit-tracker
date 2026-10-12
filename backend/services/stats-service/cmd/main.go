package main

import (
	"net"
	"os"
	"shared"
	pbHabit "shared/pb/habit"
	pb "shared/pb/stats"
	pbUser "shared/pb/user"
	"stats-service/db"
	"stats-service/handler"
	"stats-service/internal/repository"
	"stats-service/internal/service"

	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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

	logger.Info("Starting Stats service server")

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

	statsRepo := repository.NewStatsRepository(queries)

	clientTLS, err := shared.LoadClientTLSCredentials()
	if err != nil {
		logger.Fatal("failed to load client TLS credentials", zap.Error(err))
	}

	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	if userServiceAddr == "" {
		userServiceAddr = "localhost:8080"
	}
	userConn, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(clientTLS))
	if err != nil {
		logger.Fatal("failed to create user service client", zap.Error(err))
	}
	defer userConn.Close()

	habitServiceAddr := os.Getenv("HABIT_SERVICE_ADDR")
	if habitServiceAddr == "" {
		habitServiceAddr = "localhost:8082"
	}
	habitConn, err := grpc.NewClient(habitServiceAddr, grpc.WithTransportCredentials(clientTLS))
	if err != nil {
		logger.Fatal("failed to create habit service client", zap.Error(err))
	}
	defer habitConn.Close()

	userServiceClient := pbUser.NewUserServiceClient(userConn)
	habitServiceClient := pbHabit.NewHabitServiceClient(habitConn)
	statsService := service.NewStatsService(statsRepo, userServiceClient, habitServiceClient)
	statsHandler := handler.NewStatsHandler(statsService, logger, userServiceClient)

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

	pb.RegisterStatsServiceServer(grpcServer, statsHandler)

	if err := grpcServer.Serve(list); err != nil {
		logger.Fatal("The server is not running", zap.Error(err))
	}
}
