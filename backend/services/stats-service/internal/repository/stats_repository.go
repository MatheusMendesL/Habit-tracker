package repository

import (
	"context"
	"database/sql"
	"errors"
	"stats-service/db"
	AppErr "stats-service/internal/errors"

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

func (r *StatsRepository) GetUserStats(ctx context.Context, userID uuid.UUID) (db.UserStats, error) {
	userStats, err := r.q.GetUserStats(ctx, userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.UserStats{}, AppErr.ErrUserStatsNotFound
		}
		return db.UserStats{}, err
	}

	return userStats, nil
}

func (r *StatsRepository) DeleteUserStats(ctx context.Context, userID uuid.UUID) error {
	return r.q.DeleteUserStats(ctx, userID)
}

func (r *StatsRepository) RegisterHabitCompletion(ctx context.Context, userID uuid.UUID) error {
	return r.q.IncrementCompletedHabits(ctx, userID)
}
