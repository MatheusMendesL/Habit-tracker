package repository

import (
	"context"
	"database/sql"
	"habit-service/db"
)

func toNullString(param string) sql.NullString {
	return sql.NullString{
		String: param,
		Valid:  true,
	}
}

type HabitRepository struct {
	q *db.Queries
}

func NewHabitRepository(q *db.Queries) *HabitRepository {
	return &HabitRepository{q: q}
}

func (r *HabitRepository) GetHabitByID(ctx context.Context, habitId int32) (db.Habit, error) {
	res, err := r.q.GetHabitByID(ctx, habitId)

	if err != nil {
		return db.Habit{}, err
	}

	return db.Habit{
		ID:          res.ID,
		UserID:      res.UserID,
		Name:        res.Name,
		Description: res.Description,
		ImageUrl:    res.ImageUrl,
	}, err
}

func (r *HabitRepository) CreateHabit(ctx context.Context, userId int32, name, desc, imageUrl string) (db.Habit, error) {
	params := db.CreateHabitParams{
		UserID:      userId,
		Name:        name,
		Description: toNullString(desc),
		ImageUrl:    toNullString(imageUrl),
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
