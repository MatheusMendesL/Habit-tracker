package service

import (
	"context"
	pbUser "shared/pb/user"
	"stats-service/db"
	"stats-service/internal/repository"

	"github.com/google/uuid"
)

type StatsService struct {
	pbUser.UserServiceClient
	repo *repository.StatsRepository
}

// dps eu add o ser e o habit aqui como modulo
func NewStatsService(r *repository.StatsRepository, userClient pbUser.UserServiceClient) *StatsService {
	return &StatsService{
		repo:              r,
		UserServiceClient: userClient,
	}
}

func (s *StatsService) CreateUserStats(ctx context.Context, userID uuid.UUID) (db.UserStats, error) {

	return db.UserStats{}, nil
}
