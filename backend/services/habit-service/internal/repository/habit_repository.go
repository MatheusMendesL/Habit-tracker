package repository

import (
	"context"
	"database/sql"
	"errors"
	"habit-service/db"
	AppErr "habit-service/internal/errors"
)

type HabitRepository struct {
	q *db.Queries
}

func NewHabitRepository(q *db.Queries) *HabitRepository {
	return &HabitRepository{q: q}
}

func (r *HabitRepository) GetHabitByID(ctx context.Context, habitId int32) (db.Habit, error) {
	res, err := r.q.GetHabitByID(ctx, habitId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Habit{}, AppErr.ErrHabitNotFound
		}
		return db.Habit{}, err
	}

	return res, nil
}

type CreateHabitParams struct {
	UserID      int32
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

	res, err := r.q.CreateHabit(ctx, params)

	if err != nil {
		return db.Habit{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return db.Habit{}, err
	}

	return r.GetHabitByID(ctx, int32(id))
}

func (r *HabitRepository) GetRoutineByID(ctx context.Context, routineId int32) (db.Routine, error) {
	res, err := r.q.GetRoutineByID(ctx, routineId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Routine{}, AppErr.ErrRoutineNotFound
		}
		return db.Routine{}, err
	}

	return res, nil
}

type CreateRoutineParams struct {
	UserID int32
	Name   string
}

func (r *HabitRepository) CreateRoutine(ctx context.Context, arg CreateRoutineParams) (db.Routine, error) {
	return db.Routine{}, nil
}
