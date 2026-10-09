package main

import (
	"log"
	"net/http"
	"os"
	"shared"

	"gateway/internal/clients"
	habitHandlers "gateway/internal/handlers/habit"
	socialHandler "gateway/internal/handlers/social"
	statsHandler "gateway/internal/handlers/stats"
	userHandler "gateway/internal/handlers/user"
	"gateway/internal/routes"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	pbHabit "shared/pb/habit"
	pbSocial "shared/pb/social"
	pbStats "shared/pb/stats"
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
	tlsCredentials, err := shared.LoadClientTLSCredentials()
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(typeServerUser, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		return err
	}
	defer conn.Close()

	typeServerSocial := os.Getenv("SOCIAL_SERVICE_ADDR")
	if typeServerSocial == "" {
		typeServerSocial = "localhost:8081"
	}
	socialConn, err := grpc.Dial(typeServerSocial, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		return err
	}
	defer socialConn.Close()

	typeServerHabit := os.Getenv("HABIT_SERVICE_ADDR")
	if typeServerHabit == "" {
		typeServerHabit = "localhost:8082"
	}
	habitConn, err := grpc.Dial(typeServerHabit, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		return err
	}
	defer habitConn.Close()

	typeServerStats := os.Getenv("STATS_SERVICE_ADDR")
	if typeServerStats == "" {
		typeServerStats = "localhost:8083"
	}
	statsConn, err := grpc.Dial(typeServerStats, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		return err
	}
	defer statsConn.Close()

	userClient := clients.NewUserClient(pb.NewUserServiceClient(conn))
	userHandler := userHandler.NewUserHandler(userClient)
	socialClient := clients.NewSocialClient(pbSocial.NewSocialServiceClient(socialConn))
	socialHandler := socialHandler.NewSocialHandler(socialClient)
	habitClient := clients.NewHabitClient(pbHabit.NewHabitServiceClient(habitConn))
	habitHandler := habitHandlers.NewHabitHandler(habitClient)
	routineClient := clients.NewRoutineClient(pbHabit.NewRoutineServiceClient(habitConn))
	routineHandler := habitHandlers.NewRoutineHandler(routineClient)
	statsClient := clients.NewStatsClient(pbStats.NewStatsServiceClient(statsConn))
	statsHTTPHandler := statsHandler.NewStatsHandler(statsClient)

	r := routes.ControlRoutes(userHandler, socialHandler, habitHandler, routineHandler, statsHTTPHandler)
	logger.Info("Gateway is running on :8081")

	portServer := os.Getenv("PORT")

	if err := http.ListenAndServe(portServer, r); err != nil {
		return err
	}

	log.Println("Gateway stopped")
	return nil
}
