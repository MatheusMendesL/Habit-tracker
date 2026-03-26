package handler

import (
	"context"
	"database/sql"
	"errors"
	pbSocial "shared/pb/social"
	pbUser "shared/pb/user"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	AppErr "social/internal/errors"
	"social/internal/service"
	"time"
)

type SocialHandler struct {
	pbUser.UserServiceClient
	pbSocial.UnimplementedSocialServiceServer
	socialService *service.SocialService
	logger        *zap.Logger
}

const defaultTimeout = 3 * time.Second

func (s *SocialHandler) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
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

func NewSocialHandler(
	s *service.SocialService,
	logger *zap.Logger,
	userClient pbUser.UserServiceClient,
) *SocialHandler {
	return &SocialHandler{
		socialService:     s,
		logger:            logger,
		UserServiceClient: userClient,
	}
}

func (s *SocialHandler) StartFollowing(ctx context.Context, req *pbSocial.StartFollowingRequest) (*pbSocial.StartFollowingResponse, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	_, err := s.GetUserByID(context.Background(), &pbUser.GetUserByIDRequest{UserId: 1})

	if err != nil {
		return nil, err
	}

	return &pbSocial.StartFollowingResponse{
		Success: true,
	}, nil
}
