package handler

import (
	"context"
	pb "shared/pb/user"
	"user-service/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

func (s *UserHandler) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {

	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	user, err := s.userService.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserByIDResponse{
		User: &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
