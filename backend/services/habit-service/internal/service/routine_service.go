package service

import (
	"context"
	"habit-service/db"
	AppErr "habit-service/internal/errors"
	"habit-service/internal/repository"
	pbUser "shared/pb/user"
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

func (s *RoutineService) GetRoutineByID(ctx context.Context, routineID int32) (db.Routine, error) {
	if routineID <= 0 {
		return db.Routine{}, AppErr.ErrInvalidArgument
	}

	return s.repo.GetRoutineByID(ctx, routineID)
}
