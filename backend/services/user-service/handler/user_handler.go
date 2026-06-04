package handler

import (
	"context"
	"database/sql"
	"errors"
	pb "shared/pb/user"
	"time"
	"user-service/db"
	AppErr "user-service/internal/errors"
	"user-service/internal/service"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
	logger      *zap.Logger
}

const defaultTimeout = 3 * time.Second

func (s *UserHandler) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
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

func NewUserHandler(s *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: s,
		logger:      logger,
	}
}

func (s *UserHandler) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {

	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	userID, err := uuid.Parse(req.UserId)

	if err != nil {
		s.logger.Warn("Invalid User id",
			zap.String("user_id", userID.String()),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	id := userID.String()

	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error("error to execute GetUserByID method",
			zap.String("user_id", id),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("The method GetUserByID was ok",
		zap.String("user_id", id),
	)

	return &pb.GetUserByIDResponse{
		User: &pb.User{
			Id:    id,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (s *UserHandler) SearchUser(ctx context.Context, req *pb.SearchUserRequest) (*pb.SearchUserResponse, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	name := req.GetName()
	email := req.GetEmail()

	if name == "" && email == "" {
		s.logger.Warn("Empty Search User Name or Email")
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInformedIncorrect.Error())
	}

	users, err := s.userService.SearchUser(ctx, name, email)
	if err != nil {
		s.logger.Error("error to execute SearchUser method",
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	var pbUsers []*pb.User
	for _, u := range users {

		id := u.ID.String()

		pbUsers = append(pbUsers, &pb.User{
			Id:    id,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	s.logger.Info("The method SearchUser was ok",
		zap.String("name", name),
		zap.String("email", email),
	)

	return &pb.SearchUserResponse{
		User: pbUsers,
	}, nil
}

func (s *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	userId, err := uuid.Parse(req.UserId)

	if err != nil {
		s.logger.Warn("invalid user_id",
			zap.String("user_id", req.UserId),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	err = s.userService.DeleteUser(ctx, userId)
	if err != nil {
		s.logger.Error("error to execute DeleteUser method",
			zap.String("user_id", req.UserId),
			zap.Error(err),
		)

		return &pb.DeleteUserResponse{
			Success: false,
		}, ReceiveErrors(err)
	}

	s.logger.Info("DeleteUser method was ok",
		zap.String("user_id", req.UserId),
	)

	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}

func (s *UserHandler) EditUser(ctx context.Context, req *pb.EditUserRequest) (*pb.EditUserResponse, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	if req.Name == nil && req.Email == nil {
		s.logger.Warn("empty data",
			zap.String("user_id", req.UserId),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInformedIncorrect.Error())
	}

	NameUser := req.GetName()
	EmailUser := req.GetEmail()

	userId, err := uuid.Parse(req.UserId)

	if err != nil {
		s.logger.Warn("invalid user_id",
			zap.String("user_id", req.UserId),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	user, err := s.userService.GetUserByID(ctx, userId)
	if err != nil {
		s.logger.Error("error to execute EditUser method",
			zap.String("user_id", req.UserId),
			zap.String("name", req.GetName()),
			zap.String("email", req.GetEmail()),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	if NameUser == "" {
		NameUser = user.Name
	}

	if EmailUser == "" {
		EmailUser = user.Email
	}

	params := db.UpdateUserParams{
		ID:    userId,
		Name:  NameUser,
		Email: EmailUser,
	}

	user, err = s.userService.EditUser(ctx, params)
	if err != nil {
		s.logger.Error("error to execute EditUser method",
			zap.String("user_id", req.UserId),
			zap.String("name", req.GetName()),
			zap.String("email", req.GetEmail()),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("EditUser method was ok",
		zap.String("user_id", req.UserId),
		zap.String("name", req.GetName()),
		zap.String("email", req.GetEmail()),
	)

	return &pb.EditUserResponse{
		User: &pb.User{
			Id:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (s *UserHandler) EditPassword(ctx context.Context, req *pb.EditPasswordRequest) (*pb.EditPasswordResponse, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	newPassword := req.NewPassword
	if newPassword == "" {
		s.logger.Warn("newPassword is null")
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrNullField.Error())
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		s.logger.Warn("user_id is null",
			zap.String("user_id", req.UserId),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	err = s.userService.EditPassword(ctx, &db.UpdatePasswordParams{
		Password: newPassword,
		ID:       userID,
	})

	if err != nil {
		s.logger.Error("error to execute UpdatePassword method",
			zap.String("user_id", req.UserId),
			zap.Error(err),
		)
		return &pb.EditPasswordResponse{
			Success: false,
		}, ReceiveErrors(err)
	}

	s.logger.Info("UpdatePassword method was ok",
		zap.String("user_id", req.UserId),
	)

	return &pb.EditPasswordResponse{
		Success: true,
	}, nil
}

func (s *UserHandler) GetUsersByIDs(ctx context.Context, req *pb.GetUsersByIDsRequest) (*pb.GetUsersByIDsResponse, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	var userIDs []uuid.UUID

	for _, id := range req.UserIds {
		parsedID, err := uuid.Parse(id)
		if err != nil {
			s.logger.Warn("user_id is null",
				zap.Any("user_id", id),
			)
			return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
		}

		userIDs = append(userIDs, parsedID)
	}

	users, err := s.userService.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		s.logger.Error("error to execute GetUsersByIds method",
			zap.Any("user_ids", req.UserIds),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	var pbUsers []*pb.User
	for _, u := range users {
		id := u.ID.String()
		pbUsers = append(pbUsers, &pb.User{
			Id:    id,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	s.logger.Info("GetUsersByIds method was ok",
		zap.Any("users", pbUsers),
	)

	return &pb.GetUsersByIDsResponse{
		Users: pbUsers,
	}, nil

}
