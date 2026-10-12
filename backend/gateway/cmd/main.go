package main

import (
	"log"
	"net/http"
	"os"

	"gateway/internal/clients"
	habitHandlers "gateway/internal/handlers/habit"
	socialHandler "gateway/internal/handlers/social"
	statsHandler "gateway/internal/handlers/stats"
	userHandler "gateway/internal/handlers/user"
	"gateway/internal/routes"
	"shared"
	pbHabit "shared/pb/habit"
	pbSocial "shared/pb/social"
	pbStats "shared/pb/stats"
	pbUser "shared/pb/user"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
		logger.Error("Failed to execute code", zap.Error(err))
		os.Exit(1)
	}
}

func serviceAddress(envName, defaultAddress string) string {
	if address := os.Getenv(envName); address != "" {
		return address
	}
	return defaultAddress
}

func runServer(logger *zap.Logger) error {
	tlsCredentials, err := shared.LoadClientTLSCredentials()
	if err != nil {
		return err
	}

	addresses := []string{
		serviceAddress("USER_SERVICE_ADDR", "localhost:8080"),
		serviceAddress("SOCIAL_SERVICE_ADDR", "localhost:8081"),
		serviceAddress("HABIT_SERVICE_ADDR", "localhost:8082"),
		serviceAddress("STATS_SERVICE_ADDR", "localhost:8083"),
	}
	connections := make([]*grpc.ClientConn, 0, len(addresses))
	for _, address := range addresses {
		conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(tlsCredentials))
		if err != nil {
			for _, openedConn := range connections {
				_ = openedConn.Close()
			}
			return err
		}
		connections = append(connections, conn)
	}
	defer func() {
		for _, conn := range connections {
			_ = conn.Close()
		}
	}()

	userClient := clients.NewUserClient(pbUser.NewUserServiceClient(connections[0]))
	socialClient := clients.NewSocialClient(pbSocial.NewSocialServiceClient(connections[1]))
	habitClient := clients.NewHabitClient(pbHabit.NewHabitServiceClient(connections[2]))
	routineClient := clients.NewRoutineClient(pbHabit.NewRoutineServiceClient(connections[2]))
	statsClient := clients.NewStatsClient(pbStats.NewStatsServiceClient(connections[3]))

	r := routes.ControlRoutes(
		userHandler.NewUserHandler(userClient),
		socialHandler.NewSocialHandler(socialClient),
		habitHandlers.NewHabitHandler(habitClient, routineClient),
		habitHandlers.NewRoutineHandler(routineClient, habitClient),
		statsHandler.NewStatsHandler(statsClient),
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8084"
	}
	logger.Info("Gateway is running", zap.String("port", port))
	if err := http.ListenAndServe(port, r); err != nil {
		return err
	}

	return nil
}
