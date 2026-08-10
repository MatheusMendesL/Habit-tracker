package repository

import (
	"context"
	"stats-service/db"

	"github.com/google/uuid"
)

type StatsRepository struct {
	q *db.Queries
}

func NewStatsRepository(q *db.Queries) *StatsRepository {
	return &StatsRepository{q: q}
}

func (r *StatsRepository) CreateUserStats(ctx context.Context, userID uuid.UUID) (db.UserStats, error) {
	userStats, err := r.q.CreateUserStats(ctx, userID)

	if err != nil {
		return db.UserStats{}, err
	}

	return userStats, nil
}
