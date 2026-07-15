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
	repo        *repository.HabitRepository
	routineRepo *repository.RoutineRepository
}

func NewHabitService(
	habitRepo *repository.HabitRepository,
	routineRepo *repository.RoutineRepository,
	userClient pbUser.UserServiceClient,
) *HabitService {
	return &HabitService{
		repo:              habitRepo,
		routineRepo:       routineRepo,
		UserServiceClient: userClient,
	}
}

func (s *HabitService) GetHabitByID(ctx context.Context, habitId uuid.UUID) (db.Habit, error) {
	if habitId == uuid.Nil {
		return db.Habit{}, AppErr.ErrInvalidArgument
	}

	habit, err := s.repo.GetHabitByID(ctx, habitId)
	return habit, ReturnError(err, AppErr.ErrHabitNotFound)
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

func (s *HabitService) EditHabit(ctx context.Context, habit db.UpdateHabitParams) (db.Habit, error) {
	if habit.ID == uuid.Nil {
		return db.Habit{}, AppErr.ErrInvalidArgument
	}

	_, err := s.GetHabitByID(ctx, habit.ID)
	if err = ReturnError(err, AppErr.ErrHabitNotFound); err != nil {
		return db.Habit{}, err
	}

	return s.repo.EditHabit(ctx, habit)
}

func (s *HabitService) DeleteHabit(ctx context.Context, habitID uuid.UUID) error {
	if habitID == uuid.Nil {
		return AppErr.ErrInvalidArgument
	}

	_, err := s.GetHabitByID(ctx, habitID)
	if err = ReturnError(err, AppErr.ErrHabitNotFound); err != nil {
		return err
	}

	return s.repo.DeleteHabit(ctx, habitID)
}

func (s *HabitService) ListHabitsByUser(ctx context.Context, userID uuid.UUID) ([]db.Habit, error) {
	if userID == uuid.Nil {
		return []db.Habit{}, AppErr.ErrInvalidArgument
	}

	_, err := s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: userID.String()})

	if err = ReturnError(err, AppErr.ErrUserNotFound); err != nil {
		return []db.Habit{}, err
	}

	return s.repo.ListHabitsByUser(ctx, userID)
}

func (s *HabitService) ListHabitsByRoutine(ctx context.Context, routineID uuid.UUID) ([]db.Habit, error) {
	if routineID == uuid.Nil {
		return []db.Habit{}, AppErr.ErrInvalidArgument
	}

	_, err := s.routineRepo.GetRoutineByID(ctx, routineID)

	if err = ReturnError(err, AppErr.ErrRoutineNotFound); err != nil {
		return []db.Habit{}, err
	}

	return s.repo.ListHabitsByRoutine(ctx, routineID)
}

func (s *HabitService) MarkHabitCompleted(ctx context.Context, markHabit db.MarkHabitCompletedParams) error {
	if markHabit.HabitID == uuid.Nil || markHabit.CompletedAt.IsZero() {
		return AppErr.ErrInvalidArgument
	}

	_, err := s.GetHabitByID(ctx, markHabit.HabitID)

	if err = ReturnError(err, AppErr.ErrHabitNotFound); err != nil {
		return err
	}

	return s.repo.MarkHabitCompleted(ctx, markHabit)
}

func (s *HabitService) UnmarkHabitCompleted(ctx context.Context, unmarkHabit db.UnmarkHabitCompletedParams) error {
	if unmarkHabit.HabitID == uuid.Nil || unmarkHabit.CompletedAt.IsZero() {
		return AppErr.ErrInvalidArgument
	}

	_, err := s.GetHabitByID(ctx, unmarkHabit.HabitID)

	if err = ReturnError(err, AppErr.ErrHabitNotFound); err != nil {
		return err
	}

	return s.repo.UnmarkHabitCompleted(ctx, unmarkHabit)
}

func (s *HabitService) GetHabitLogs(ctx context.Context, arg db.GetHabitLogsParams) ([]db.GetHabitLogsRow, error) {
	if arg.HabitID == uuid.Nil || arg.StartDate.IsZero() || arg.EndDate.IsZero() {
		return []db.GetHabitLogsRow{}, AppErr.ErrInvalidArgument
	}

	_, err := s.GetHabitByID(ctx, arg.HabitID)
	if err = ReturnError(err, AppErr.ErrHabitNotFound); err != nil {
		return []db.GetHabitLogsRow{}, err
	}

	return s.repo.GetHabitsLog(ctx, arg)
}
