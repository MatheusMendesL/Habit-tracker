package service

import (
	"context"
	"errors"
	"habit-service/db"
	AppErr "habit-service/internal/errors"
	"habit-service/internal/repository"
	pbUser "shared/pb/user"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RoutineService struct {
	pbUser.UserServiceClient
	repo *repository.RoutineRepository
}

func NewRoutineService(r *repository.RoutineRepository, userClient pbUser.UserServiceClient) *RoutineService {
	return &RoutineService{
		repo:              r,
		UserServiceClient: userClient,
	}
}

func (s *RoutineService) CreateRoutine(ctx context.Context, arg repository.CreateRoutineParams) (db.Routine, error) {
	if arg.UserID == uuid.Nil || arg.Name == "" {
		return db.Routine{}, AppErr.ErrInvalidArgument
	}

	_, err := s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: arg.UserID.String()})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return db.Routine{}, AppErr.ErrUserNotFound
		}
		return db.Routine{}, err
	}

	return s.repo.CreateRoutine(ctx, arg)
}

func (s *RoutineService) GetRoutineByID(ctx context.Context, routineID uuid.UUID) (db.Routine, error) {
	if routineID == uuid.Nil {
		return db.Routine{}, AppErr.ErrInvalidArgument
	}

	return s.repo.GetRoutineByID(ctx, routineID)
}

func (s *RoutineService) EditRoutine(ctx context.Context, routine db.UpdateRoutineParams) (db.Routine, error) {
	if routine.ID == uuid.Nil {
		return db.Routine{}, AppErr.ErrInvalidArgument
	}

	_, err := s.GetRoutineByID(ctx, routine.ID)

	if err != nil {
		if errors.Is(err, AppErr.ErrRoutineNotFound) {
			return db.Routine{}, AppErr.ErrRoutineNotFound
		}
		return db.Routine{}, err
	}

	return s.repo.EditRoutine(ctx, routine)
}

func (s *RoutineService) DeleteRoutine(ctx context.Context, routineID uuid.UUID) error {
	if routineID == uuid.Nil {
		return AppErr.ErrInvalidArgument
	}

	_, err := s.GetRoutineByID(ctx, routineID)

	if err != nil {
		if errors.Is(err, AppErr.ErrRoutineNotFound) {
			return AppErr.ErrRoutineNotFound
		}
		return err
	}

	return s.repo.DeleteRoutine(ctx, routineID)
}

func (s *RoutineService) ListRoutinesByUser(ctx context.Context, userID uuid.UUID) ([]db.Routine, error) {
	if userID == uuid.Nil {
		return []db.Routine{}, AppErr.ErrInvalidArgument
	}

	_, err := s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: userID.String()})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return []db.Routine{}, AppErr.ErrUserNotFound
		}
		return []db.Routine{}, err
	}

	return s.repo.ListRoutinesByUser(ctx, userID)
}
