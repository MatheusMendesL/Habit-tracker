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

func (s *RoutineHandler) EditRoutine(ctx context.Context, req *pbHabit.EditRoutineRequest) (*pbHabit.EditRoutineResponse, error) {
	ctx, cancel := WithTimeout(ctx)
	defer cancel()

	routineID, routineName := req.RoutineId, req.Name

	if err := s.Verification(routineID, "routine id", "routine_id"); err != nil {
		return nil, err
	}

	routineIDNew, err := uuid.Parse(routineID)

	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			AppErr.ErrInvalidArgument.Error(),
		)
	}

	params := db.UpdateRoutineParams{
		ID: routineIDNew,
	}

	name := ""

	if routineName != nil {
		params.Name = utils.ToNullString(*routineName)
		name = *routineName
	}

	routine, err := s.RoutineService.EditRoutine(ctx, params)

	if err != nil {
		s.logger.Error("Error to execute EditRoutine method",
			zap.String("Name", name),
			zap.String("routine_id", routineID),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("EditRoutine method was ok",
		zap.String("Name", name),
		zap.String("routine_id", routineID),
	)

	return &pbHabit.EditRoutineResponse{
		Routine: utils.ToProtoRoutine(routine),
	}, nil

}

func (s *RoutineHandler) DeleteRoutine(ctx context.Context, req *pbHabit.DeleteRoutineRequest) (*pbHabit.DeleteRoutineResponse, error) {
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

	err = s.RoutineService.DeleteRoutine(ctx, routineID)
	if err != nil {
		s.logger.Error("error to execute DeleteRoutine method",
			zap.String("routine_id", routineID.String()),
			zap.Error(err),
		)

		return &pbHabit.DeleteRoutineResponse{
			Success: false,
		}, ReceiveErrors(err)
	}

	s.logger.Info("DeleteRoutine method was ok",
		zap.String("routine_id", routineID.String()),
	)

	return &pbHabit.DeleteRoutineResponse{
		Success: true,
	}, nil

}

func (s *RoutineHandler) ListRoutinesByUser(ctx context.Context, req *pbHabit.ListRoutinesByUserRequest) (*pbHabit.ListRoutinesByUserResponse, error) {
	ctx, cancel := WithTimeout(ctx)
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

	routines, err := s.RoutineService.ListRoutinesByUser(ctx, userID)

	if err != nil {
		s.logger.Error("error to execute ListRoutinesByUser method",
			zap.String("user_id", req.UserId),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	pbRoutines := make([]*pbHabit.Routine, 0, len(routines))

	for _, routine := range routines {
		pbRoutines = append(pbRoutines, utils.ToProtoRoutine(routine))
	}

	// use this func to all returns that is necessary an array

	s.logger.Info("ListRoutinesByUser method was ok",
		zap.String("user_id", userID.String()),
		zap.Int("total", len(routines)),
	)

	return &pbHabit.ListRoutinesByUserResponse{
		Routines: pbRoutines,
	}, nil

}

/*
AddHabitToRoutine
RemoveHabitFromRoutine
*/

func (s *RoutineHandler) AddHabitToRoutine(ctx context.Context, req *pbHabit.AddHabitToRoutineRequest) (*pbHabit.AddHabitToRoutineResponse, error) {
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

	err = s.RoutineService.AddHabitToRoutine(ctx, db.AddHabitToRoutineParams{
		RoutineID: routineID,
		HabitID:   habitID,
	})

	if err != nil {
		s.logger.Error("error to execute AddHabitToRoutine method",
			zap.String("habit_id", habitID.String()),
			zap.String("routine_id", routineID.String()),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("AddHabitToRoutine method was ok",
		zap.String("habit_id", habitID.String()),
		zap.String("routine_id", routineID.String()),
	)

	return &pbHabit.AddHabitToRoutineResponse{
		Success: true,
	}, nil
}

func (s *RoutineHandler) RemoveHabitFromRoutine(ctx context.Context, req *pbHabit.RemoveHabitFromRoutineRequest) (*pbHabit.RemoveHabitFromRoutineResponse, error) {
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

	err = s.RoutineService.RemoveHabitFromRoutine(ctx, db.RemoveHabitFromRoutineParams{
		RoutineID: routineID,
		HabitID:   habitID,
	})

	if err != nil {
		s.logger.Error("error to execute RemoveHabitFromRoutine method",
			zap.String("habit_id", habitID.String()),
			zap.String("routine_id", routineID.String()),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("RemoveHabitFromRoutine method was ok",
		zap.String("habit_id", habitID.String()),
		zap.String("routine_id", routineID.String()),
	)

	return &pbHabit.RemoveHabitFromRoutineResponse{
		Success: true,
	}, nil
}
