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

func WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, defaultTimeout)
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
	ctx, cancel := WithTimeout(ctx)
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

/*

what to do in this order:

EditHabit
DeleteHabit
ListHabitsByUser
ListHabitsByRoutine
MarkHabitCompleted
UnmarkHabitCompleted
GetHabitLogs
*/
