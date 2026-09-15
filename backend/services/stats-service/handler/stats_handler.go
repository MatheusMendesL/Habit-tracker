package handler

import (
	"context"
	"database/sql"
	"errors"
	pbUser "shared/pb/user"
	AppErr "stats-service/internal/errors"
	"stats-service/internal/service"
	"stats-service/internal/utils"
	"time"

	pbStats "shared/pb/stats"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StatsHandler struct {
	pbStats.UnimplementedStatsServiceServer
	StatsService      *service.StatsService
	logger            *zap.Logger
	UserServiceClient pbUser.UserServiceClient
}

func NewStatsHandler(
	s *service.StatsService,
	logger *zap.Logger,
	userClient pbUser.UserServiceClient,
) *StatsHandler {
	return &StatsHandler{
		StatsService:      s,
		logger:            logger,
		UserServiceClient: userClient,
	}
}

const defaultTimeout = 3 * time.Second

func WithTimeout(ctx context.Context, customTimeout ...time.Duration) (context.Context, context.CancelFunc) {
	timeout := defaultTimeout

	if len(customTimeout) > 0 {
		timeout = customTimeout[0]
	}

	return context.WithTimeout(ctx, timeout)
}

func (s *StatsHandler) Verification(val any, name string, nameVal string) error {
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

func ReceiveErrors(err error) error {
	switch {
	case errors.Is(err, AppErr.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, AppErr.ErrNullField):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, AppErr.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, AppErr.ErrUserStatsNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, sql.ErrNoRows):
		return status.Error(codes.NotFound, err.Error())

	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func (s *StatsHandler) CreateUserStats(ctx context.Context, req *pbStats.CreateUserStatsRequest) (*pbStats.CreateUserStatsResponse, error) {
	ctx, cancel := WithTimeout(ctx, 7*time.Second)
	defer cancel()

	userIDnew, err := uuid.Parse(req.UserId)

	if err != nil {
		s.logger.Error("error to transform to uuid",
			zap.Any("user_id", userIDnew),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	idUser := req.UserId

	if err := s.Verification(idUser, "user id", "user_id"); err != nil {
		return nil, err
	}

	userStats, err := s.StatsService.CreateUserStats(ctx, userIDnew)

	if err != nil {
		s.logger.Error("error to execute CreateStats method",
			zap.String("userid", userIDnew.String()),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("CreateUserStats method was ok",
		zap.String("user_id", userIDnew.String()),
	)

	return &pbStats.CreateUserStatsResponse{
		Stats: utils.ToProtoStats(userStats),
	}, nil
}

func (s *StatsHandler) GetUserStats(ctx context.Context, req *pbStats.GetUserStatsRequest) (*pbStats.GetUserStatsResponse, error) {
	ctx, cancel := WithTimeout(ctx)
	defer cancel()

	userIDnew, err := uuid.Parse(req.UserId)

	if err != nil {
		s.logger.Error("error to transform to uuid",
			zap.Any("user_id", userIDnew),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	idUser := req.UserId

	if err := s.Verification(idUser, "user id", "user_id"); err != nil {
		return nil, err
	}

	userStats, err := s.StatsService.GetUserStats(ctx, userIDnew)

	if err != nil {
		s.logger.Error("error to execute GetUserStats method",
			zap.String("user_id", userIDnew.String()),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("GetUserStats method was ok",
		zap.String("user_id", userIDnew.String()),
	)

	return &pbStats.GetUserStatsResponse{Stats: utils.ToProtoStats(userStats)}, nil
}

func (s *StatsHandler) DeleteUserStats(ctx context.Context, req *pbStats.DeleteUserStatsRequest) (*pbStats.DeleteUserStatsResponse, error) {
	ctx, cancel := WithTimeout(ctx)
	defer cancel()

	userIDnew, err := uuid.Parse(req.UserId)

	if err != nil {
		s.logger.Error("error to transform to uuid",
			zap.Any("user_id", userIDnew),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	idUser := req.UserId

	if err := s.Verification(idUser, "user id", "user_id"); err != nil {
		return nil, err
	}

	err = s.StatsService.DeleteUserStats(ctx, userIDnew)

	if err != nil {
		s.logger.Error("error to execute DeleteUserStats method",
			zap.String("user_id", userIDnew.String()),
			zap.Error(err),
		)
	}

	s.logger.Info("DeleteUserStats method was ok",
		zap.String("user_id", userIDnew.String()),
	)

	return &pbStats.DeleteUserStatsResponse{
		Success: false,
	}, ReceiveErrors(err)
}

func (s *StatsHandler) RegisterHabitCompletion(ctx context.Context, req *pbStats.RegisterHabitCompletionRequest) (*pbStats.RegisterHabitCompletionResponse, error) {
	ctx, cancel := WithTimeout(ctx)
	defer cancel()

	userIDnew, err := uuid.Parse(req.UserId)
	if err != nil {
		s.logger.Error("error to transform to uuid", zap.Error(err))
		return nil, ReceiveErrors(err)
	}

	if err := s.Verification(req.UserId, "user id", "user_id"); err != nil {
		return nil, err
	}
	if err := s.Verification(req.HabitId, "habit id", "habit_id"); err != nil {
		return nil, err
	}

	var completedAt time.Time
	if req.CompletedAt != nil {
		completedAt = req.CompletedAt.AsTime()
	} else {
		completedAt = time.Now()
	}

	err = s.StatsService.RegisterHabitCompletion(ctx, userIDnew, req.HabitId, completedAt)
	if err != nil {
		s.logger.Error("error to execute RegisterHabitCompletion method", zap.Error(err))
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("RegisterHabitCompletion method was ok",
		zap.String("user_id", userIDnew.String()),
		zap.String("habit_id", req.HabitId),
	)

	return &pbStats.RegisterHabitCompletionResponse{
		Success: true,
	}, nil
}
