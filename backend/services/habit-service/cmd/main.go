package main

import (
	"habit-service/db"
	"habit-service/handler"
	"habit-service/internal/repository"
	"habit-service/internal/service"
	"log"
	"net"
	"os"
	pb "shared/pb/habit"
	pbUser "shared/pb/user"

	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	_ = godotenv.Load()

	logger.Info("Starting server")

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

	HabitRepo := repository.NewHabitRepository(queries)

	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Não foi possível conectar: %v", err)
	}
	defer conn.Close()

	userServiceClient := pbUser.NewUserServiceClient(conn)
	habitService := service.NewHabitService(HabitRepo, userServiceClient)
	habitHandler := handler.NewHabitHandler(habitService, logger, userServiceClient)

	/*tlsCredentials, err := loadTLCredentials()

	if err != nil {
		logger.Fatal("failed to load TLS credentials", zap.Error(err))
	}*/

	grpcServer := grpc.NewServer(
		/*grpc.Creds(tlsCredentials),*/
		grpc.UnaryInterceptor(
			grpcZap.UnaryServerInterceptor(logger),
		),
	)

	pb.RegisterHabitServiceServer(grpcServer, habitHandler)

	if err := grpcServer.Serve(list); err != nil {
		logger.Fatal("The server is not running", zap.Error(err))
	}
}
