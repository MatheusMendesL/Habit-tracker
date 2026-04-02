package repository

import (
	"context"
	"database/sql"
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
		if err == sql.ErrNoRows {
			return db.Habit{}, AppErr.ErrUserNotFound
		}
		return db.Habit{}, err
	}

	return db.Habit{
		ID:          res.ID,
		UserID:      res.UserID,
		Name:        res.Name,
		Description: res.Description,
		ImageUrl:    res.ImageUrl,
	}, nil
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
