package service

import (
	"context"
	"errors"
	pbHabit "shared/pb/habit"
	pbUser "shared/pb/user"
	"stats-service/db"
	AppErr "stats-service/internal/errors"
	"stats-service/internal/repository"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StatsService struct {
	pbHabit.HabitServiceClient
	pbUser.UserServiceClient
	repo          *repository.StatsRepository
	routineClient pbHabit.RoutineServiceClient
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

func NewStatsService(r *repository.StatsRepository, userClient pbUser.UserServiceClient, habitClient pbHabit.HabitServiceClient, routineClients ...pbHabit.RoutineServiceClient) *StatsService {
	service := &StatsService{
		repo:               r,
		UserServiceClient:  userClient,
		HabitServiceClient: habitClient,
	}
	if len(routineClients) > 0 {
		service.routineClient = routineClients[0]
	}
	return service
}

func (s *StatsService) verifyHabitOwnership(ctx context.Context, userID uuid.UUID, habitID string) error {
	habit, err := s.HabitServiceClient.GetHabitByID(ctx, &pbHabit.GetHabitByIDRequest{HabitId: habitID})
	if err != nil {
		return err
	}
	if habit == nil || habit.Habit == nil {
		return status.Error(codes.Internal, "habit service returned an empty response")
	}
	if habit.Habit.UserId != userID.String() {
		return status.Error(codes.PermissionDenied, "habit does not belong to the user")
	}
	return nil
}

func (s *StatsService) verifyRoutineOwnership(ctx context.Context, userID uuid.UUID, routineID string) error {
	if s.routineClient == nil {
		return status.Error(codes.Unavailable, "routine service client is unavailable")
	}
	routine, err := s.routineClient.GetRoutineByID(ctx, &pbHabit.GetRoutineByIDRequest{RoutineId: routineID})
	if err != nil {
		return err
	}
	if routine == nil || routine.Routine == nil {
		return status.Error(codes.Internal, "routine service returned an empty response")
	}
	if routine.Routine.UserId != userID.String() {
		return status.Error(codes.PermissionDenied, "routine does not belong to the user")
	}
	return nil
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

func (s *StatsService) RegisterHabitCompletion(ctx context.Context, userID uuid.UUID, habitID string, completedAt time.Time) error {
	if userID == uuid.Nil || habitID == "" {
		return AppErr.ErrInvalidArgument
	}
	completedAt = completedAt.UTC().Truncate(time.Microsecond)

	_, err := s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: userID.String()})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return AppErr.ErrUserNotFound
		}
		return err
	}

	if err := s.verifyHabitOwnership(ctx, userID, habitID); err != nil {
		return err
	}
	return s.repo.RegisterHabitCompletion(ctx, userID, completedAt)
}

func (s *StatsService) UndoHabitCompletion(ctx context.Context, userID uuid.UUID, habitID string, completedAt time.Time) error {
	if userID == uuid.Nil || habitID == "" {
		return AppErr.ErrInvalidArgument
	}
	completedAt = completedAt.UTC().Truncate(time.Microsecond)

	_, err := s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: userID.String()})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return AppErr.ErrUserNotFound
		}
		return err
	}

	if err := s.verifyHabitOwnership(ctx, userID, habitID); err != nil {
		return err
	}
	return s.repo.UndoHabitCompletion(ctx, userID, completedAt)
}

func (s *StatsService) RegisterRoutineCompletion(ctx context.Context, userID uuid.UUID, routineID string, completedAt time.Time) error {
	if userID == uuid.Nil || routineID == "" {
		return AppErr.ErrInvalidArgument
	}
	completedAt = completedAt.UTC().Truncate(time.Microsecond)

	_, err := s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: userID.String()})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return AppErr.ErrUserNotFound
		}
		return err
	}

	if err := s.verifyRoutineOwnership(ctx, userID, routineID); err != nil {
		return err
	}
	return s.repo.RegisterRoutineCompletion(ctx, userID, completedAt)
}

func (s *StatsService) UndoRoutineCompletion(ctx context.Context, userID uuid.UUID, routineID string, completedAt time.Time) error {
	if userID == uuid.Nil || routineID == "" {
		return AppErr.ErrInvalidArgument
	}
	completedAt = completedAt.UTC().Truncate(time.Microsecond)

	_, err := s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: userID.String()})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return AppErr.ErrUserNotFound
		}
		return err
	}

	if err := s.verifyRoutineOwnership(ctx, userID, routineID); err != nil {
		return err
	}
	return s.repo.UndoRoutineCompletion(ctx, userID, completedAt)
}
