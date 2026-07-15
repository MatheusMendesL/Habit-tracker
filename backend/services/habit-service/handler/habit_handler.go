package handler

import (
	"context"
	"database/sql"
	"errors"
	"habit-service/db"
	AppErr "habit-service/internal/errors"
	"habit-service/internal/repository"
	"habit-service/internal/service"
	"habit-service/internal/utils"
	pbHabit "shared/pb/habit"
	pbUser "shared/pb/user"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HabitHandler struct {
	pbUser.UserServiceClient
	pbHabit.UnimplementedHabitServiceServer
	HabitService *service.HabitService
	logger       *zap.Logger
}

const defaultTimeout = 3 * time.Second

func WithTimeout(ctx context.Context, customTimeout ...time.Duration) (context.Context, context.CancelFunc) {
	timeout := defaultTimeout

	if len(customTimeout) > 0 {
		timeout = customTimeout[0]
	}

	return context.WithTimeout(ctx, timeout)
}

func ReceiveErrors(err error) error {
	switch {
	case errors.Is(err, AppErr.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, AppErr.ErrNullField):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, AppErr.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, AppErr.ErrRoutineNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, AppErr.ErrHabitNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, sql.ErrNoRows):
		return status.Error(codes.NotFound, err.Error())

	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func (s *HabitHandler) Verification(val any, name string, nameVal string) error {
	switch v := val.(type) {
	case string:
		if v == "" {
			s.logger.Warn("invalid "+name,
				zap.String(nameVal, v),
			)
			return status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
		}
	case int32:
		if v <= 0 {
			s.logger.Warn("invalid "+name,
				zap.Int32(nameVal, v),
			)
			return status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
		}
	case uuid.UUID:
		if v == uuid.Nil {
			s.logger.Warn("invalid "+name,
				zap.String(nameVal, v.String()),
			)
			return status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
		}
	}

	return nil
}

func NewHabitHandler(
	s *service.HabitService,
	logger *zap.Logger,
	userClient pbUser.UserServiceClient,
) *HabitHandler {
	return &HabitHandler{
		HabitService:      s,
		logger:            logger,
		UserServiceClient: userClient,
	}
}

func (s *HabitHandler) CreateHabit(ctx context.Context, req *pbHabit.CreateHabitRequest) (*pbHabit.CreateHabitResponse, error) {
	ctx, cancel := WithTimeout(ctx, 7*time.Second)
	defer cancel()

	reqHabit := req.Habit

	if err := s.Verification(reqHabit.UserId, "user id", "user_id"); err != nil {
		return nil, err
	}

	if err := s.Verification(reqHabit.Name, "habit name", "habit_name"); err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(reqHabit.UserId)

	if err != nil {
		s.logger.Error("error to transform to uuid",
			zap.Any("user_id", userID),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	args := db.CreateHabitParams{
		UserID:      userID,
		Name:        reqHabit.Name,
		Description: utils.ToNullString(reqHabit.Description),
		ImageUrl:    utils.ToNullString(reqHabit.ImageUrl),
	}

	habit, err := s.HabitService.CreateHabit(ctx, repository.CreateHabitParams(args))

	if err != nil {
		s.logger.Error("error to execute CreateHabit method",
			zap.Any("habit", reqHabit),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	return &pbHabit.CreateHabitResponse{
		Habit: utils.ToProtoHabit(habit),
	}, nil
}

func (s *HabitHandler) GetHabitByID(ctx context.Context, req *pbHabit.GetHabitByIDRequest) (*pbHabit.GetHabitByIDResponse, error) {
	ctx, cancel := WithTimeout(ctx)
	defer cancel()

	habitID := req.HabitId

	if err := s.Verification(habitID, "habit id", "habit_id"); err != nil {
		return nil, err
	}

	habitIDNew, err := uuid.Parse(habitID)

	if err != nil {
		s.logger.Error("error to transform to uuid",
			zap.Any("user_id", habitIDNew),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	habit, err := s.HabitService.GetHabitByID(ctx, habitIDNew)

	if err != nil {
		if errors.Is(err, AppErr.ErrHabitNotFound) {
			s.logger.Warn("Habit not found",
				zap.String("habit_id", habitIDNew.String()),
				zap.Error(err),
			)
			return nil, status.Error(codes.NotFound, AppErr.ErrHabitNotFound.Error())
		}
		s.logger.Error("error to execute GetHabitByID method",
			zap.String("habit_id", habitIDNew.String()),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("The method GetHabitByID was ok",
		zap.String("habit_id", habitIDNew.String()),
	)

	return &pbHabit.GetHabitByIDResponse{
		Habit: utils.ToProtoHabit(habit),
	}, nil

}

func (s *HabitHandler) EditHabit(ctx context.Context, req *pbHabit.EditHabitRequest) (*pbHabit.EditHabitResponse, error) {
	ctx, cancel := WithTimeout(ctx)
	defer cancel()

	habitID, HabitName, HabitDesc, HabitimgURL := req.HabitId, req.Name, req.Description, req.ImageUrl

	if err := s.Verification(habitID, "habit id", "habit_id"); err != nil {
		return nil, err
	}

	habitUUID, err := uuid.Parse(habitID)

	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			AppErr.ErrInvalidArgument.Error(),
		)
	}

	params := db.UpdateHabitParams{
		ID: habitUUID,
	}

	name := ""

	if HabitName != nil {
		params.Name = utils.ToNullString(*HabitName)
		name = *HabitName
	}

	desc := ""

	if HabitDesc != nil {
		params.Description = utils.ToNullString(*HabitDesc)
		desc = *HabitDesc
	}

	imgURL := ""
	if HabitimgURL != nil {
		params.ImageUrl = utils.ToNullString(*HabitimgURL)
		imgURL = *HabitimgURL
	}

	habit, err := s.HabitService.EditHabit(ctx, params)
	if err != nil {
		s.logger.Error("Error to execute EditHabit method",
			zap.String("Name", name),
			zap.String("description", desc),
			zap.String("imageURL", imgURL),
			zap.String("Habit_id", habitID),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("EditHabit method was ok",
		zap.String("Name", name),
		zap.String("description", desc),
		zap.String("imageURL", imgURL),
		zap.String("Habit_id", habitID),
	)

	return &pbHabit.EditHabitResponse{
		Habit: utils.ToProtoHabit(habit),
	}, nil
}

func (s *HabitHandler) DeleteHabit(ctx context.Context, req *pbHabit.DeleteHabitRequest) (*pbHabit.DeleteHabitResponse, error) {
	ctx, cancel := WithTimeout(ctx)
	defer cancel()

	habitID, err := uuid.Parse(req.HabitId)

	if err != nil {
		s.logger.Warn("invalid habit id",
			zap.String("habit_id", req.HabitId),
		)

		return nil, status.Error(
			codes.InvalidArgument,
			AppErr.ErrInvalidArgument.Error(),
		)
	}

	if err = s.Verification(habitID, "habit id", "habit_id"); err != nil {
		return nil, err
	}

	err = s.HabitService.DeleteHabit(ctx, habitID)
	if err != nil {
		s.logger.Error("error to execute DeleteHabit method",
			zap.String("habit_id", habitID.String()),
			zap.Error(err),
		)

		return &pbHabit.DeleteHabitResponse{
			Success: false,
		}, ReceiveErrors(err)
	}

	s.logger.Info("DeleteHabit method was ok",
		zap.String("habit_id", habitID.String()),
	)

	return &pbHabit.DeleteHabitResponse{
		Success: true,
	}, nil
}

func (s *HabitHandler) ListHabitsByUser(ctx context.Context, req *pbHabit.ListHabitsByUserRequest) (*pbHabit.ListHabitsByUserResponse, error) {
	ctx, cancel := WithTimeout(ctx, 7*time.Second)
	defer cancel()

	userID, err := uuid.Parse(req.UserId)

	if err != nil {
		s.logger.Warn("invalid user id",
			zap.String("user_id", req.UserId),
		)

		return nil, status.Error(
			codes.InvalidArgument,
			AppErr.ErrInvalidArgument.Error(),
		)
	}
	if err = s.Verification(userID, "user id", "user_id"); err != nil {
		return nil, err
	}

	habits, err := s.HabitService.ListHabitsByUser(ctx, userID)

	if err != nil {
		s.logger.Error("error to execute ListHabitsByUser method",
			zap.String("user_id", req.UserId),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	pbHabits := make([]*pbHabit.Habit, 0, len(habits))

	for _, habit := range habits {
		pbHabits = append(pbHabits, utils.ToProtoHabit(habit))
	}

	s.logger.Info("ListHabitsByUser method was ok",
		zap.String("user_id", userID.String()),
		zap.Int("total", len(habits)),
	)

	return &pbHabit.ListHabitsByUserResponse{
		Habits: pbHabits,
	}, nil
}

func (s *HabitHandler) ListHabitsByRoutine(ctx context.Context, req *pbHabit.ListHabitsByRoutineRequest) (*pbHabit.ListHabitsByRoutineResponse, error) {
	ctx, cancel := WithTimeout(ctx)
	defer cancel()

	routineID, err := uuid.Parse(req.RoutineId)

	if err != nil {
		s.logger.Warn("invalid routine id",
			zap.String("routine_id", req.RoutineId),
		)

		return nil, status.Error(
			codes.InvalidArgument,
			AppErr.ErrInvalidArgument.Error(),
		)
	}
	if err = s.Verification(routineID, "routine id", "routine_id"); err != nil {
		return nil, err
	}

	habits, err := s.HabitService.ListHabitsByRoutine(ctx, routineID)

	if err != nil {
		s.logger.Error("error to execute ListHabitsByRoutine method",
			zap.String("routine_id", req.RoutineId),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	pbHabits := make([]*pbHabit.Habit, 0, len(habits))

	for _, habit := range habits {
		pbHabits = append(pbHabits, utils.ToProtoHabit(habit))
	}

	s.logger.Info("ListHabitsByRoutine method was ok",
		zap.String("routine_id", routineID.String()),
		zap.Int("total", len(habits)),
	)

	return &pbHabit.ListHabitsByRoutineResponse{
		Habits: pbHabits,
	}, nil

}

func (s *HabitHandler) MarkHabitCompleted(ctx context.Context, req *pbHabit.MarkHabitCompletedRequest) (*pbHabit.MarkHabitCompletedResponse, error) {
	ctx, cancel := WithTimeout(ctx)
	defer cancel()

	habitID, err := uuid.Parse(req.HabitId)
	if err != nil {
		s.logger.Warn("invalid habit id",
			zap.String("habit_id", req.HabitId),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	var completedAt time.Time
	if req.CompletedAt != nil {
		completedAt = req.CompletedAt.AsTime().UTC()
	} else {
		s.logger.Warn("invalid completed_at date",
			zap.Any("completedAt", req.CompletedAt),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	params := db.MarkHabitCompletedParams{
		HabitID:     habitID,
		CompletedAt: completedAt,
	}

	err = s.HabitService.MarkHabitCompleted(ctx, params)
	if err != nil {
		s.logger.Error("error to execute MarkHabitCompleted method",
			zap.String("habit_id", habitID.String()),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("MarkHabitCompleted method was ok",
		zap.String("habit_id", habitID.String()),
	)

	return &pbHabit.MarkHabitCompletedResponse{Success: true}, nil
}

func (s *HabitHandler) UnmarkHabitCompleted(ctx context.Context, req *pbHabit.UnmarkHabitCompletedRequest) (*pbHabit.UnmarkHabitCompletedResponse, error) {
	ctx, cancel := WithTimeout(ctx)
	defer cancel()

	habitID, err := uuid.Parse(req.HabitId)
	if err != nil {
		s.logger.Warn("invalid habit id",
			zap.String("habit_id", req.HabitId),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	var completedAt time.Time
	if req.CompletedAt != nil {
		completedAt = req.CompletedAt.AsTime().UTC()
	} else {
		s.logger.Warn("invalid completed_at date",
			zap.Any("completedAt", req.CompletedAt),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	params := db.UnmarkHabitCompletedParams{
		HabitID:     habitID,
		CompletedAt: completedAt,
	}

	err = s.HabitService.UnmarkHabitCompleted(ctx, params)
	if err != nil {
		s.logger.Error("error to execute UnmarkHabitCompleted method",
			zap.String("habit_id", habitID.String()),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("UnmarkHabitCompleted method was ok",
		zap.String("habit_id", habitID.String()),
	)

	return &pbHabit.UnmarkHabitCompletedResponse{Success: true}, nil
}

/*

what to do in this order:

GetHabitLogs
*/
