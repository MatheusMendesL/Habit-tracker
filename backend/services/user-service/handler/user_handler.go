package handler

import (
	"context"
	"database/sql"
	"errors"
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
	// .GetName() e .GetEmail() retornam "" se o campo for nulo no JSON
	name := req.GetName()
	email := req.GetEmail()

	if name == "" && email == "" {
		return nil, status.Error(codes.InvalidArgument, "Informe nome ou email")
	}

	users, err := s.userService.SearchUser(ctx, name, email)
	if err != nil {
		return nil, err
	}

	// Mapeia para a resposta do Protobuf
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
