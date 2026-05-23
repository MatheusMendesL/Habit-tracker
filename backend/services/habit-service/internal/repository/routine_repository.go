package repository

import (
	"context"
	"database/sql"
	"errors"
	"habit-service/db"
	AppErr "habit-service/internal/errors"
)

type RoutineRepository struct {
	q *db.Queries
}

func NewRoutineRepository(q *db.Queries) *RoutineRepository {
	return &RoutineRepository{q: q}
}

func (r *RoutineRepository) GetRoutineByID(ctx context.Context, routineId int32) (db.Routine, error) {
	res, err := r.q.GetRoutineByID(ctx, routineId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Routine{}, AppErr.ErrRoutineNotFound
		}
		return db.Routine{}, err
	}

	return res, nil
}
