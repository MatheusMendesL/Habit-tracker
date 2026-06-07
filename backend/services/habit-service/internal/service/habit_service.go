package service

import (
	"context"
	"habit-service/db"
	AppErr "habit-service/internal/errors"
	"habit-service/internal/repository"
	pbUser "shared/pb/user"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HabitService struct {
	pbUser.UserServiceClient
	repo *repository.HabitRepository
}

func NewHabitService(r *repository.HabitRepository, userClient pbUser.UserServiceClient) *HabitService {
	return &HabitService{
		repo:              r,
		UserServiceClient: userClient,
	}
}

func (s *HabitService) GetHabitByID(ctx context.Context, habitId uuid.UUID) (db.Habit, error) {
	if habitId == uuid.Nil {
		return db.Habit{}, AppErr.ErrInvalidArgument
	}

	return s.repo.GetHabitByID(ctx, habitId)
}

func (s *HabitService) CreateHabit(ctx context.Context, arg repository.CreateHabitParams) (db.Habit, error) {
	if arg.UserID == uuid.Nil || arg.Name == "" {
		return db.Habit{}, AppErr.ErrInvalidArgument
	}

	_, err := s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: arg.UserID.String()})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return db.Habit{}, AppErr.ErrUserNotFound
		}
		return db.Habit{}, err
	}

	return s.repo.CreateHabit(ctx, arg)
}
