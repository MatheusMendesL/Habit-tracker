package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	pb "shared/pb/user"
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
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, sql.ErrNoRows):
		return status.Error(codes.NotFound, AppErr.ErruUserNotFound.Error())
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
	if req.Name == nil || req.Email == nil {
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrNullField.Error())
	}

	users, err := s.userService.SearchUser(ctx, *req.Name, *req.Email)

	if err != nil {
		return nil, ReceiveErrors(err)
	}

	fmt.Print(users)
	return nil, nil
}
