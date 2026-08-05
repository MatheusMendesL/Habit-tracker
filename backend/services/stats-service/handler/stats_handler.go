package handler

import (
	"stats-service/internal/service"

	pbStats "shared/pb/stats"

	"go.uber.org/zap"
)

type RoutineHandler struct {
	pbStats.UnimplementedStatsServiceServer
	StatsService *service.StatsService
	logger       *zap.Logger
}