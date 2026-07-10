package service

import (
	"context"
	"errors"
	"habit-service/db"
	AppErr "habit-service/internal/errors"
	"habit-service/internal/repository"
	pbUser "shared/pb/user"

	"github.com/google/uuid"
)

func ReturnError(err error, target error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, target) {
		return target
	}

	return err
}

type RoutineService struct {
	pbUser.UserServiceClient
	repo         *repository.RoutineRepository
	habitService *HabitService
}

func NewRoutineService(r *repository.RoutineRepository, userClient pbUser.UserServiceClient, habitService *HabitService) *RoutineService {
	return &RoutineService{
		repo:              r,
		UserServiceClient: userClient,
		habitService:      habitService,
	}
}

func (s *RoutineService) CreateRoutine(ctx context.Context, arg repository.CreateRoutineParams) (db.Routine, error) {
	if arg.UserID == uuid.Nil || arg.Name == "" {
		return db.Routine{}, AppErr.ErrInvalidArgument
	}

	_, err := s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: arg.UserID.String()})

	if err = ReturnError(err, AppErr.ErrUserNotFound); err != nil {
		return db.Routine{}, err
	}

	return s.repo.CreateRoutine(ctx, arg)
}

func (s *RoutineService) GetRoutineByID(ctx context.Context, routineID uuid.UUID) (db.Routine, error) {
	if routineID == uuid.Nil {
		return db.Routine{}, AppErr.ErrInvalidArgument
	}

	routine, err := s.repo.GetRoutineByID(ctx, routineID)

	return routine, ReturnError(err, AppErr.ErrRoutineNotFound)
}

func (s *RoutineService) EditRoutine(ctx context.Context, routine db.UpdateRoutineParams) (db.Routine, error) {
	if routine.ID == uuid.Nil {
		return db.Routine{}, AppErr.ErrInvalidArgument
	}

	_, err := s.GetRoutineByID(ctx, routine.ID)

	if err = ReturnError(err, AppErr.ErrRoutineNotFound); err != nil {
		return db.Routine{}, err
	}

	return s.repo.EditRoutine(ctx, routine)
}

func (s *RoutineService) DeleteRoutine(ctx context.Context, routineID uuid.UUID) error {
	if routineID == uuid.Nil {
		return AppErr.ErrInvalidArgument
	}

	_, err := s.GetRoutineByID(ctx, routineID)

	if err = ReturnError(err, AppErr.ErrRoutineNotFound); err != nil {
		return err
	}

	return s.repo.DeleteRoutine(ctx, routineID)
}

func (s *RoutineService) ListRoutinesByUser(ctx context.Context, userID uuid.UUID) ([]db.Routine, error) {
	if userID == uuid.Nil {
		return []db.Routine{}, AppErr.ErrInvalidArgument
	}

	_, err := s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: userID.String()})

	if err = ReturnError(err, AppErr.ErrUserNotFound); err != nil {
		return []db.Routine{}, err
	}

	return s.repo.ListRoutinesByUser(ctx, userID)
}

func (s *RoutineService) AddHabitToRoutine(ctx context.Context, params db.AddHabitToRoutineParams) error {
	if params.RoutineID == uuid.Nil || params.HabitID == uuid.Nil {
		return AppErr.ErrInvalidArgument
	}

	_, err := s.GetRoutineByID(ctx, params.RoutineID)

	if err = ReturnError(err, AppErr.ErrRoutineNotFound); err != nil {
		return err
	}

	_, err = s.habitService.GetHabitByID(ctx, params.HabitID)

	if err = ReturnError(err, AppErr.ErrHabitNotFound); err != nil {
		return err
	}

	return s.repo.AddHabitToRoutine(ctx, params)
}

func (s *RoutineService) RemoveHabitFromRoutine(ctx context.Context, params db.RemoveHabitFromRoutineParams) error {
	if params.RoutineID == uuid.Nil || params.HabitID == uuid.Nil {
		return AppErr.ErrInvalidArgument
	}

	_, err := s.GetRoutineByID(ctx, params.RoutineID)

	if err = ReturnError(err, AppErr.ErrRoutineNotFound); err != nil {
		return err
	}

	_, err = s.habitService.GetHabitByID(ctx, params.HabitID)

	if err = ReturnError(err, AppErr.ErrHabitNotFound); err != nil {
		return err
	}

	return s.repo.RemoveHabitFromRoutine(ctx, params)
}
