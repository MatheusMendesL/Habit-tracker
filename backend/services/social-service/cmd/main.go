package main

import (
	"log"
	"net"
	"os"
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

	socialRepo := repository.NewSocialRepository(queries)

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Não foi possível conectar: %v", err)
	}
	defer conn.Close()

	userServiceClient := pbUser.NewUserServiceClient(conn)
	socialService := service.NewSocialService(socialRepo)
	socialHandler := handler.NewSocialHandler(socialService, logger, userServiceClient)

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

	pb.RegisterSocialServiceServer(grpcServer, socialHandler)

	if err := grpcServer.Serve(list); err != nil {
		logger.Fatal("The server is not running", zap.Error(err))
	}
}
