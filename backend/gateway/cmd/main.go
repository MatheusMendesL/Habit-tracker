package main

import (
	"log"
	"net/http"
	"os"
	"shared"

	"gateway/internal/clients"
	userHandler "gateway/internal/handlers/user"
	"gateway/internal/routes"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	pb "shared/pb/user"

	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found for gateway; relying on environment variables")
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	if err := runServer(logger); err != nil {
		logger.Error("Failed to execute code",
			zap.Error(err),
		)
		os.Exit(1)
	}

	logger.Info("All systems offline")
}

func runServer(logger *zap.Logger) error {
	typeServerUser := os.Getenv("USER_SERVICE_ADDR")
	if typeServerUser == "" {
		typeServerUser = "localhost:8080"
	}

	tlsCredentials, err := shared.LoadClientTLSCredentials()
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(typeServerUser, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		return err
	}
	defer conn.Close()

	userClient := clients.NewUserClient(pb.NewUserServiceClient(conn))
	userHandler := userHandler.NewUserHandler(userClient)

	r := routes.ControlRoutes(userHandler)
	logger.Info("Gateway is running on :8081")

	portServer := os.Getenv("PORT")

	if err := http.ListenAndServe(portServer, r); err != nil {
		return err
	}

	log.Println("Gateway stopped")
	return nil
}
