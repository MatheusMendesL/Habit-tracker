package repository

import (
	"context"
	"database/sql"
	"errors"
	"habit-service/db"
	AppErr "habit-service/internal/errors"

	"github.com/google/uuid"
)

type HabitRepository struct {
	q *db.Queries
}

func NewHabitRepository(q *db.Queries) *HabitRepository {
	return &HabitRepository{q: q}
}

type CreateHabitParams struct {
	UserID      uuid.UUID
	Name        string
	Description sql.NullString
	ImageUrl    sql.NullString
}

func (r *HabitRepository) CreateHabit(ctx context.Context, arg CreateHabitParams) (db.Habit, error) {
	params := db.CreateHabitParams{
		UserID:      arg.UserID,
		Name:        arg.Name,
		Description: arg.Description,
		ImageUrl:    arg.ImageUrl,
	}

	habit, err := r.q.CreateHabit(ctx, params)
	if err != nil {
		return db.Habit{}, err
	}

	return habit, nil
}

func (r *HabitRepository) GetHabitByID(ctx context.Context, habitId uuid.UUID) (db.Habit, error) {
	res, err := r.q.GetHabitByID(ctx, habitId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Habit{}, AppErr.ErrHabitNotFound
		}
		return db.Habit{}, err
	}

	return res, nil
}

func (r *HabitRepository) EditHabit(ctx context.Context, req db.UpdateHabitParams) (db.Habit, error) {
	habit, err := r.q.UpdateHabit(ctx, req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Habit{}, AppErr.ErrHabitNotFound
		}

		return db.Habit{}, err
	}

	return habit, nil
}

func (r *HabitRepository) DeleteHabit(ctx context.Context, habitID uuid.UUID) error {
	return r.q.DeleteHabit(ctx, habitID)
}

func (r *HabitRepository) ListHabitsByUser(ctx context.Context, userID uuid.UUID) ([]db.Habit, error) {
	habits, err := r.q.ListHabitsByUser(ctx, userID)

	if err != nil {
		return []db.Habit{}, err
	}

	return habits, nil
}

func (r *HabitRepository) ListHabitsByRoutine(ctx context.Context, routineID uuid.UUID) ([]db.Habit, error) {
	habits, err := r.q.ListHabitsByRoutine(ctx, routineID)

	if err != nil {
		return []db.Habit{}, err
	}

	return habits, nil
}

func (r *HabitRepository) MarkHabitCompleted(ctx context.Context, req db.MarkHabitCompletedParams) error {
	return r.q.MarkHabitCompleted(ctx, req)
}

func (r *HabitRepository) UnmarkHabitCompleted(ctx context.Context, req db.UnmarkHabitCompletedParams) error {
	return r.q.UnmarkHabitCompleted(ctx, req)
}

func (r *HabitRepository) GetHabitsLog(ctx context.Context, req db.GetHabitLogsParams) ([]db.GetHabitLogsRow, error) {
	logs, err := r.q.GetHabitLogs(ctx, req)

	if err != nil {
		return []db.GetHabitLogsRow{}, err
	}

	return logs, nil
}
