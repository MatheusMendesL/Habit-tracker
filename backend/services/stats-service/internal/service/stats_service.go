package service

import (
	"context"
	"errors"
	pbHabit "shared/pb/habit"
	pbUser "shared/pb/user"
	"stats-service/db"
	AppErr "stats-service/internal/errors"
	"stats-service/internal/repository"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StatsService struct {
	pbHabit.HabitServiceClient
	pbUser.UserServiceClient
	repo *repository.StatsRepository
}

func ReturnError(err error, target error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, target) {
		return target
	}

	return err
}

func NewStatsService(r *repository.StatsRepository, userClient pbUser.UserServiceClient, habitClient pbHabit.HabitServiceClient) *StatsService {
	return &StatsService{
		repo:               r,
		UserServiceClient:  userClient,
		HabitServiceClient: habitClient,
	}
}

func (s *StatsService) CreateUserStats(ctx context.Context, userID uuid.UUID) (db.UserStats, error) {

	if userID == uuid.Nil {
		return db.UserStats{}, AppErr.ErrInvalidArgument
	}

	stats, err := s.repo.CreateUserStats(ctx, userID)
	return stats, ReturnError(err, AppErr.ErrUserNotFound)
}

func (s *StatsService) GetUserStats(ctx context.Context, userID uuid.UUID) (db.UserStats, error) {
	if userID == uuid.Nil {
		return db.UserStats{}, AppErr.ErrInvalidArgument
	}

	_, err := s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: userID.String()})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return db.UserStats{}, AppErr.ErrUserNotFound
		}
		return db.UserStats{}, err
	}

	stats, err := s.repo.GetUserStats(ctx, userID)
	return stats, ReturnError(err, AppErr.ErrUserStatsNotFound)
}

func (s *StatsService) DeleteUserStats(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return AppErr.ErrInvalidArgument
	}

	_, err := s.GetUserStats(ctx, userID)

	if err = ReturnError(err, AppErr.ErrUserStatsNotFound); err != nil {
		return err
	}

	return s.repo.DeleteUserStats(ctx, userID)
}
