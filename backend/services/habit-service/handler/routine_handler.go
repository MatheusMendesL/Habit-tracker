package handler

import (
	"context"
	"errors"
	"habit-service/db"
	AppErr "habit-service/internal/errors"
	"habit-service/internal/repository"
	"habit-service/internal/service"
	"habit-service/internal/utils"
	pbHabit "shared/pb/habit"
	pbUser "shared/pb/user"

	"github.com/google/uuid"
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

/*
UpdateRoutine
DeleteRoutine
ListRoutinesByUser
*/

func (s *RoutineHandler) Verification(val any, name string, nameVal string) error {
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

func (s *RoutineHandler) CreateRoutine(ctx context.Context, req *pbHabit.CreateRoutineRequest) (*pbHabit.CreateRoutineResponse, error) {
	ctx, cancel := WithTimeout(ctx)
	defer cancel()

	reqRoutine := req.Routine

	if err := s.Verification(reqRoutine.UserId, "user id", "user_id"); err != nil {
		return nil, err
	}

	if err := s.Verification(reqRoutine.Name, "routine name", "routine_name"); err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(reqRoutine.UserId)

	if err != nil {
		s.logger.Error("error to transform to uuid",
			zap.Any("routine", reqRoutine),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	args := db.CreateRoutineParams{
		UserID: userID,
		Name:   reqRoutine.Name,
	}

	routine, err := s.RoutineService.CreateRoutine(ctx, repository.CreateRoutineParams(args))

	if err != nil {
		s.logger.Error("error to execute CreateRoutine method",
			zap.Any("routine", reqRoutine),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("The method CreateRoutine was ok")

	return &pbHabit.CreateRoutineResponse{Routine: utils.ToProtoRoutine(routine)}, nil

}

func (s *RoutineHandler) GetRoutineByID(ctx context.Context, req *pbHabit.GetRoutineByIDRequest) (*pbHabit.GetRoutineByIDResponse, error) {
	ctx, cancel := WithTimeout(ctx)
	defer cancel()

	routineID := req.RoutineId

	if err := s.Verification(routineID, "routine id", "routine_id"); err != nil {
		return nil, err
	}

	routineIDNew, err := uuid.Parse(routineID)

	if err != nil {
		s.logger.Error("error to transform to uuid",
			zap.Any("routine_id", routineIDNew),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	routine, err := s.RoutineService.GetRoutineByID(ctx, routineIDNew)

	if err != nil {
		if errors.Is(err, AppErr.ErrRoutineNotFound) {
			s.logger.Warn("Routine not found",
				zap.String("routine_id", routineIDNew.String()),
				zap.Error(err),
			)
			return nil, status.Error(codes.NotFound, AppErr.ErrRoutineNotFound.Error())
		}
		s.logger.Error("error to execute GetRoutineByID method",
			zap.String("routine_id", routineIDNew.String()),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("The method GetRoutineByID was ok",
		zap.String("routine_id", routineIDNew.String()),
	)

	return &pbHabit.GetRoutineByIDResponse{
		Routine: utils.ToProtoRoutine(routine),
	}, nil
}
