package main

import (
	"log"
	"net/http"
	"os"

	"gateway/internal/clients"
	userHandler "gateway/internal/handlers/user"
	"gateway/internal/routes"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "shared/pb/user"

	"go.uber.org/zap"
)

func main() {
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
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	userClient := clients.NewUserClient(pb.NewUserServiceClient(conn))
	userHandler := userHandler.NewUserHandler(userClient)

	r := routes.ControlRoutes(userHandler)
	logger.Info("Gateway is running on :8081")

	if err := http.ListenAndServe(":8081", r); err != nil {
		return err
	}

	log.Println("Gateway stopped")
	return nil
}
