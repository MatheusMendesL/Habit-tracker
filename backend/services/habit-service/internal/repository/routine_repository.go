package repository

import (
	"context"
	"database/sql"
	"errors"
	"habit-service/db"
	AppErr "habit-service/internal/errors"
	"strconv"

	"github.com/google/uuid"
)

type RoutineRepository struct {
	q *db.Queries
}

func NewRoutineRepository(q *db.Queries) *RoutineRepository {
	return &RoutineRepository{q: q}
}

type CreateRoutineParams struct {
	UserID uuid.UUID
	Name   string
}

func (r *RoutineRepository) CreateRoutine(ctx context.Context, arg CreateRoutineParams) (db.Routine, error) {
	params := db.CreateRoutineParams{
		UserID: arg.UserID,
		Name:   arg.Name,
	}

	res, err := r.q.CreateRoutine(ctx, params)

	if err != nil {
		return db.Routine{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return db.Routine{}, err
	}

	idNew, err := uuid.Parse(strconv.FormatInt(id, 10))

	if err != nil {
		return db.Routine{}, err
	}

	return r.GetRoutineByID(ctx, idNew)
}

func (r *RoutineRepository) GetRoutineByID(ctx context.Context, routineId uuid.UUID) (db.Routine, error) {
	res, err := r.q.GetRoutineByID(ctx, routineId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Routine{}, AppErr.ErrRoutineNotFound
		}
		return db.Routine{}, err
	}

	return res, nil
}
