package repository

import (
	"context"
	"database/sql"
	"errors"
	"habit-service/db"
	AppErr "habit-service/internal/errors"

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
	return r.q.CreateRoutine(ctx, db.CreateRoutineParams{
		UserID: arg.UserID,
		Name:   arg.Name,
	})
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

func (r *RoutineRepository) EditRoutine(ctx context.Context, req db.UpdateRoutineParams) (db.Routine, error) {
	routine, err := r.q.UpdateRoutine(ctx, req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Routine{}, AppErr.ErrRoutineNotFound
		}

		return db.Routine{}, err
	}

	return routine, nil

}

func (r *RoutineRepository) DeleteRoutine(ctx context.Context, routineID uuid.UUID) error {
	return r.q.DeleteRoutine(ctx, routineID)
}

func (r *RoutineRepository) ListRoutinesByUser(ctx context.Context, userID uuid.UUID) ([]db.Routine, error) {
	return r.q.ListRoutinesByUser(ctx, userID)
}

func (r *RoutineRepository) AddHabitToRoutine(ctx context.Context, arg db.AddHabitToRoutineParams) error {
	return r.q.AddHabitToRoutine(ctx, arg)
}

func (r *RoutineRepository) RemoveHabitFromRoutine(ctx context.Context, arg db.RemoveHabitFromRoutineParams) error {
	return r.q.RemoveHabitFromRoutine(ctx, arg)
}
