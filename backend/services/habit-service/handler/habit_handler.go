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

func (s *HabitHandler) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, defaultTimeout)
}

func ReceiveErrors(err error) error {
	switch {
	case errors.Is(err, AppErr.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, AppErr.ErrNullField):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, sql.ErrNoRows):
		return status.Error(codes.NotFound, AppErr.ErrUserNotFound.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
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

func (s *HabitHandler) GetHabitByID(ctx context.Context, req *pbHabit.GetHabitByIDRequest) (*pbHabit.GetHabitByIDResponse, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	habitID := req.HabitId

	if habitID == 0 {
		s.logger.Warn("invalid habit id",
			zap.Int32("habit_id", habitID),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	habit, err := s.HabitService.GetHabitByID(ctx, habitID)

	if err != nil {
		if errors.Is(err, AppErr.ErrHabitNotFound) {
			s.logger.Warn("Habit not found",
				zap.Int32("habit_id", habitID),
				zap.Error(err),
			)
			return nil, status.Error(codes.NotFound, AppErr.ErrHabitNotFound.Error())
		}
		s.logger.Error("error to execute GetHabitByID method",
			zap.Int32("habit_id", habitID),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("The method GetHabitByID was ok",
		zap.Int32("habit_id", habitID),
	)

	return &pbHabit.GetHabitByIDResponse{
		Habit: utils.ToProtoHabit(habit),
	}, nil

}

func (s *HabitHandler) CreateHabit(ctx context.Context, req *pbHabit.CreateHabitRequest) (*pbHabit.CreateHabitResponse, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	reqHabit := req.Habit

	if reqHabit.UserId <= 0 {
		s.logger.Warn("invalid user id",
			zap.Int32("user_id", reqHabit.UserId),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	if reqHabit.Name == "" {
		s.logger.Warn("invalid habit name",
			zap.String("habit_name", reqHabit.Name),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	args := db.CreateHabitParams{
		UserID:      reqHabit.UserId,
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
