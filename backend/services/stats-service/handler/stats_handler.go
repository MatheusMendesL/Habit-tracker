package handler

import (
	pbUser "shared/pb/user"
	"stats-service/internal/service"

	pbStats "shared/pb/stats"

	"go.uber.org/zap"
)

type StatsHandler struct {
	pbStats.UnimplementedStatsServiceServer
	StatsService      *service.StatsService
	logger            *zap.Logger
	UserServiceClient pbUser.UserServiceClient
}

func NewStatsHandler(
	s *service.StatsService,
	logger *zap.Logger,
	userClient pbUser.UserServiceClient,
) *StatsHandler {
	return &StatsHandler{
		StatsService:      s,
		logger:            logger,
		UserServiceClient: userClient,
	}
}
