package handler

import (
	"context"
	"errors"
	AppErr "habit-service/internal/errors"
	"habit-service/internal/service"
	"habit-service/internal/utils"
	pbHabit "shared/pb/habit"
	pbUser "shared/pb/user"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RoutineHandler struct {
	pbUser.UserServiceClient
	pbHabit.UnimplementedRoutineServiceServer
	RoutineService *service.RoutineService
	logger         *zap.Logger
}

func (s *RoutineHandler) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, defaultTimeout)
}

func NewRoutineHandler(
	s *service.RoutineService,
	logger *zap.Logger,
	userClient pbUser.UserServiceClient,
) *RoutineHandler {
	return &RoutineHandler{
		RoutineService:    s,
		logger:            logger,
		UserServiceClient: userClient,
	}
}

func (s *RoutineHandler) GetRoutineByID(ctx context.Context, req *pbHabit.GetRoutineByIDRequest) (*pbHabit.GetRoutineByIDResponse, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	routineID := req.RoutineId

	if routineID == 0 {
		s.logger.Warn("invalid routine id",
			zap.Int32("routine_id", routineID),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	routine, err := s.RoutineService.GetRoutineByID(ctx, routineID)

	if err != nil {
		if errors.Is(err, AppErr.ErrRoutineNotFound) {
			s.logger.Warn("Routine not found",
				zap.Int32("routine_id", routineID),
				zap.Error(err),
			)
			return nil, status.Error(codes.NotFound, AppErr.ErrRoutineNotFound.Error())
		}
		s.logger.Error("error to execute GetRoutineByID method",
			zap.Int32("routine_id", routineID),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("The method GetRoutineByID was ok",
		zap.Int32("routine_id", routineID),
	)

	return &pbHabit.GetRoutineByIDResponse{
		Routine: utils.ToProtoRoutine(routine),
	}, nil
}
