package handler

import (
	"context"
	"database/sql"
	"errors"
	pb "shared/pb/user"
	"user-service/db"
	AppErr "user-service/internal/errors"
	"user-service/internal/service"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
	logger      *zap.Logger
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

	if req.UserId == 0 {
		s.logger.Warn("Invalid User id",
			zap.Int32("user_id", req.UserId),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	user, err := s.userService.GetUserByID(ctx, req.UserId)
	if err != nil {
		s.logger.Error("error to execute GetUserByID method",
			zap.Int32("user_id", req.UserId),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("The method GetUserByID was ok",
		zap.Int32("user_id", user.ID),
	)

	return &pb.GetUserByIDResponse{
		User: &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (s *UserHandler) SearchUser(ctx context.Context, req *pb.SearchUserRequest) (*pb.SearchUserResponse, error) {
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
		pbUsers = append(pbUsers, &pb.User{
			Id:    u.ID,
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
	userId := req.GetUserId()

	if userId == 0 {
		s.logger.Warn("invalid user_id",
			zap.Int32("user_id", req.UserId),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	err := s.userService.DeleteUser(ctx, req.UserId)
	if err != nil {
		s.logger.Error("error to execute DeleteUser method",
			zap.Int32("user_id", req.UserId),
			zap.Error(err),
		)

		return &pb.DeleteUserResponse{
			Success: false,
		}, ReceiveErrors(err)
	}

	s.logger.Info("DeleteUser method was ok",
		zap.Int32("user_id", req.UserId),
	)

	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}

func (s *UserHandler) UpdateUser(ctx context.Context, req *pb.EditUserRequest) (*pb.EditUserResponse, error) {

	if req.Name == nil && req.Email == nil {
		s.logger.Warn("empty data",
			zap.Int32("user_id", req.UserId),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInformedIncorrect.Error())
	}

	params := db.UpdateUserParams{
		ID:    req.UserId,
		Name:  req.Name,
		Email: req.Email,
	}

	user, err := s.userService.UpdateUser(ctx, params)
	if err != nil {
		s.logger.Error("error to execute UpdateUser method",
			zap.Int32("user_id", req.UserId),
			zap.String("name", req.GetName()),
			zap.String("email", req.GetEmail()),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("UpdateUser method was ok",
		zap.Int32("user_id", req.UserId),
		zap.String("name", req.GetName()),
		zap.String("email", req.GetEmail()),
	)

	return &pb.EditUserResponse{
		User: &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (s *UserHandler) StartFollowing(ctx context.Context, req *pb.StartFollowingRequest) (*pb.StartFollowingResponse, error) {
	if req.FollowerId == 0 || req.FolloweeId == 0 {
		s.logger.Warn("empty field",
			zap.Int32("followerID", req.FollowerId),
			zap.Int32("followeeID", req.FolloweeId),
		)
		return &pb.StartFollowingResponse{
			Success: false,
		}, status.Error(codes.InvalidArgument, AppErr.ErrNullField.Error())
	}

	err := s.userService.StartFollowing(ctx, req.FollowerId, req.FolloweeId)

	if err != nil {
		s.logger.Error("error to execute StartFollowing method",
			zap.Int32("followerID", req.FollowerId),
			zap.Int32("followeeID", req.FolloweeId),
			zap.Error(err),
		)
		return &pb.StartFollowingResponse{
			Success: false,
		}, ReceiveErrors(err)
	}

	s.logger.Info("StartFollowing method was ok",
		zap.Int32("followerID", req.FollowerId),
		zap.Int32("followeeID", req.FolloweeId),
	)

	return &pb.StartFollowingResponse{
		Success: true,
	}, nil
}

func (s *UserHandler) StopFollowing(ctx context.Context, req *pb.UnfollowRequest) (*pb.UnfollowResponse, error) {
	if req.UnfollowingId == 0 || req.UnfollowedId == 0 {
		s.logger.Warn("empty field",
			zap.Int32("UnfollowingId", req.UnfollowingId),
			zap.Int32("UnfollowedId", req.UnfollowedId),
		)
		return &pb.UnfollowResponse{
			Success: false,
		}, status.Error(codes.InvalidArgument, AppErr.ErrNullField.Error())
	}

	err := s.userService.StopFollowing(ctx, req.UnfollowingId, req.UnfollowedId)

	if err != nil {
		s.logger.Error("error to execute StopFollowing method",
			zap.Int32("UnfollowingId", req.UnfollowingId),
			zap.Int32("UnfollowedId", req.UnfollowedId),
			zap.Error(err),
		)
		return &pb.UnfollowResponse{
			Success: false,
		}, ReceiveErrors(err)
	}

	s.logger.Info("StopFollowing method was ok",
		zap.Int32("UnfollowingId", req.UnfollowingId),
		zap.Int32("UnfollowedId", req.UnfollowedId),
	)

	return &pb.UnfollowResponse{
		Success: true,
	}, nil
}

func (s *UserHandler) ListFollowers(ctx context.Context, req *pb.ListFollowersRequest) (*pb.ListFollowersResponse, error) {
	userID := req.GetUserId()
	if userID == 0 {
		s.logger.Warn("empty Field",
			zap.Int32("user_id", userID),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}
	users, err := s.userService.ListFollowers(ctx, userID)
	if err != nil {
		s.logger.Error("error to execute ListFollowers method",
			zap.Int32("user_id", userID),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	var pbUsers []*pb.User
	for _, u := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	s.logger.Info("ListFollowers method was ok",
		zap.Int32("user_id", userID),
		zap.Int("users_followed", len(pbUsers)),
	)

	return &pb.ListFollowersResponse{
		Followers: pbUsers,
	}, nil
}

func (s *UserHandler) ListFollowing(ctx context.Context, req *pb.ListFollowingRequest) (*pb.ListFollowingResponse, error) {
	userID := req.GetUserId()

	if userID == 0 {
		s.logger.Warn("empty Field",
			zap.Int32("user_id", userID),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	users, err := s.userService.ListFollowing(ctx, userID)
	if err != nil {
		s.logger.Error("error to execute ListFollowing method",
			zap.Int32("user_id", userID),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	var pbUsers []*pb.User
	for _, u := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	s.logger.Info("ListFollowing method was ok",
		zap.Int32("user_id", userID),
		zap.Int("users_following", len(pbUsers)),
	)

	return &pb.ListFollowingResponse{
		Following: pbUsers,
	}, nil
}

func (s *UserHandler) UpdatePassword(ctx context.Context, req *pb.EditPasswordRequest) (*pb.EditPasswordResponse, error) {
	newPassword := req.NewPassword
	if newPassword == "" {
		s.logger.Warn("newPassword is null")
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrNullField.Error())
	}
	userID := req.GetUserId()
	if userID == 0 {
		s.logger.Warn("user_id is null",
			zap.Int32("user_id", userID),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	err := s.userService.UpdatePassword(ctx, &db.UpdatePasswordParams{
		Password: newPassword,
		ID:       userID,
	})

	if err != nil {
		s.logger.Error("error to execute UpdatePassword method",
			zap.Int32("user_id", userID),
			zap.Error(err),
		)
		return &pb.EditPasswordResponse{
			Success: false,
		}, ReceiveErrors(err)
	}

	s.logger.Info("UpdatePassword method was ok",
		zap.Int32("user_id", userID),
	)

	return &pb.EditPasswordResponse{
		Success: true,
	}, nil
}
