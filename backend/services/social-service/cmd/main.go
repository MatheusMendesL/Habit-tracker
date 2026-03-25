package main

import (
	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
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

	userRepo := repository.NewUserRepository(queries)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, logger)

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

	pb.RegisterUserServiceServer(grpcServer, userHandler)

	if err := grpcServer.Serve(list); err != nil {
		logger.Fatal("The server is not running", zap.Error(err))
	}
}
