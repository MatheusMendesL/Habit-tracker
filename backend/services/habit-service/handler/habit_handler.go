package handler

import (
	"context"
	"database/sql"
	"errors"
	AppErr "habit-service/internal/errors"
	"habit-service/internal/service"
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
	socialService *service.HabitService
	logger        *zap.Logger
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
		socialService:     s,
		logger:            logger,
		UserServiceClient: userClient,
	}
}
