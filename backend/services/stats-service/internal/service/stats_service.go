package service

import (
	"context"
	"stats-service/db"
	"stats-service/internal/repository"

	"github.com/google/uuid"
)

type StatsService struct {
	repo *repository.StatsRepository
}

// dps eu add o ser e o habit aqui como modulo
func NewSocialService(r *repository.StatsRepository) *StatsService {
	return &StatsService{
		repo: r,
	}
}

func (s *StatsService) CreateUserStats(ctx context.Context, userID uuid.UUID ) (db.UserStats, error){


	return db.UserStats{}, nil
}