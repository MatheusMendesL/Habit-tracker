package service

import (
	"context"
	"habit-service/db"
	AppErr "habit-service/internal/errors"
	"habit-service/internal/repository"
	pbHabit "shared/pb/habit"
	pbUser "shared/pb/user"

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
	if arg.UserID <= 0 || arg.Name == "" {
		return db.Routine{}, AppErr.ErrInvalidArgument
	}

	_, err := s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: arg.UserID})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return db.Routine{}, AppErr.ErrUserNotFound
		}
		return db.Routine{}, err
	}

	return s.repo.CreateRoutine(ctx, arg)
}

func (s *RoutineService) GetRoutineByID(ctx context.Context, routineID int32) (db.Routine, error) {
	if routineID <= 0 {
		return db.Routine{}, AppErr.ErrInvalidArgument
	}

	return s.repo.GetRoutineByID(ctx, routineID)
}
