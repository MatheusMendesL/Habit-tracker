package handler

import (
	"context"
	"database/sql"
	"errors"
	pb "shared/pb/user"
	"user-service/db"
	AppErr "user-service/internal/errors"
	"user-service/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
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

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

func (s *UserHandler) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {

	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	user, err := s.userService.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, ReceiveErrors(err)
	}

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
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInformedIncorrect.Error())
	}

	users, err := s.userService.SearchUser(ctx, name, email)
	if err != nil {
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

	return &pb.SearchUserResponse{
		User: pbUsers,
	}, nil
}

func (s *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.userService.DeleteUser(ctx, req.UserId)
	if err != nil {
		return &pb.DeleteUserResponse{
			Success: false,
		}, ReceiveErrors(err)
	}

	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}

func (s *UserHandler) UpdateUser(ctx context.Context, req *pb.EditUserRequest) (*pb.EditUserResponse, error) {

	if req.Name == nil && req.Email == nil {
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInformedIncorrect.Error())
	}

	params := db.UpdateUserParams{
		ID:    req.UserId,
		Name:  req.Name,
		Email: req.Email,
	}

	user, err := s.userService.UpdateUser(ctx, params)
	if err != nil {
		return nil, ReceiveErrors(err)
	}

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
		return &pb.StartFollowingResponse{
			Success: false,
		}, status.Error(codes.InvalidArgument, AppErr.ErrNullField.Error())
	}

	err := s.userService.StartFollowing(ctx, req.FollowerId, req.FolloweeId)

	if err != nil {
		return &pb.StartFollowingResponse{
			Success: false,
		}, ReceiveErrors(err)
	}

	return &pb.StartFollowingResponse{
		Success: true,
	}, nil
}

func (s *UserHandler) StopFollowing(ctx context.Context, req *pb.UnfollowRequest) (*pb.UnfollowResponse, error) {
	if req.UnfollowingId == 0 || req.UnfollowedId == 0 {
		return &pb.UnfollowResponse{
			Success: false,
		}, status.Error(codes.InvalidArgument, AppErr.ErrNullField.Error())
	}

	err := s.userService.StopFollowing(ctx, req.UnfollowingId, req.UnfollowedId)

	if err != nil {
		return &pb.UnfollowResponse{
			Success: false,
		}, ReceiveErrors(err)
	}

	return &pb.UnfollowResponse{
		Success: true,
	}, nil
}

func (s *UserHandler) ListFollowers(ctx context.Context, req *pb.ListFollowersRequest) (*pb.ListFollowersResponse, error) {
	userID := req.GetUserId()
	if userID == 0 {
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	users, err := s.userService.ListFollowers(ctx, userID)
	if err != nil {
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

	return &pb.ListFollowersResponse{
		Followers: pbUsers,
	}, nil
}

func (s *UserHandler) ListFollowing(ctx context.Context, req *pb.ListFollowingRequest) (*pb.ListFollowingResponse, error) {
	userID := req.GetUserId()

	if userID == 0 {
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	users, err := s.userService.ListFollowing(ctx, userID)
	if err != nil {
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

	return &pb.ListFollowingResponse{
		Following: pbUsers,
	}, nil
}
