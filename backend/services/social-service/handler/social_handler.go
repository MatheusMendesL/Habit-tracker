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

	FollowerID := req.FollowerId
	FolloweeID := req.FolloweeId

	if FollowerID == 0 || FolloweeID == 0 {
		s.logger.Warn("Invalid Users ID",
			zap.Int32("FollowerID", FollowerID),
			zap.Int32("FolloweeID", FolloweeID),
		)

		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	err := s.socialService.StartFollowing(ctx, FollowerID, FolloweeID)

	if err != nil {
		s.logger.Error("error to execute StartFollowing method",
			zap.Int32("FollowerID", FollowerID),
			zap.Int32("FolloweeID", FolloweeID),
			zap.Error(err),
		)

		return nil, ReceiveErrors(err)
	}

	s.logger.Info("StartFollowing method was ok",
		zap.Int32("followerID", req.FollowerId),
		zap.Int32("followeeID", req.FolloweeId),
	)

	return &pbSocial.StartFollowingResponse{
		Success: true,
	}, nil
}

func (s *SocialHandler) Unfollow(ctx context.Context, req *pbSocial.UnfollowRequest) (*pbSocial.UnfollowResponse, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	FollowerID := req.FollowerId
	FolloweeID := req.FolloweeId

	if FollowerID == 0 || FolloweeID == 0 {
		s.logger.Warn("Invalid Users ID",
			zap.Int32("FollowerID", FollowerID),
			zap.Int32("FolloweeID", FolloweeID),
		)

		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	err := s.socialService.Unfollow(ctx, FollowerID, FolloweeID)

	if err != nil {
		s.logger.Error("error to execute Unfollow method",
			zap.Int32("FollowerID", FollowerID),
			zap.Int32("FolloweeID", FolloweeID),
			zap.Error(err),
		)

		return nil, ReceiveErrors(err)
	}

	s.logger.Info("Unfollow method was ok",
		zap.Int32("FollowerID", FollowerID),
		zap.Int32("FolloweeID", FolloweeID),
	)

	return &pbSocial.UnfollowResponse{
		Success: true,
	}, nil
}

func (s *SocialHandler) ListFollowers(ctx context.Context, req *pbSocial.ListFollowersRequest) (*pbSocial.ListFollowersResponse, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	userID := req.UserId

	if userID <= 0 {
		s.logger.Warn("Invalid User ID",
			zap.Int32("userID", userID),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	ids, err := s.socialService.ListFollowers(ctx, userID)

	if err != nil {
		s.logger.Error("error to execute ListFollowers method",
			zap.Int32("userID", userID),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("ListFollowers method was ok",
		zap.Int32("userID", userID),
	)

	return &pbSocial.ListFollowersResponse{
		FollowerIds: ids,
	}, nil
}

func (s *SocialHandler) ListFollowing(ctx context.Context, req *pbSocial.ListFollowingRequest) (*pbSocial.ListFollowingResponse, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	userID := req.UserId
	if userID <= 0 {
		s.logger.Warn("Invalid User ID",
			zap.Int32("userID", userID),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	ids, err := s.socialService.ListFollowing(ctx, userID)

	if err != nil {
		s.logger.Error("error to execute ListFollowing method",
			zap.Int32("userID", userID),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("ListFollowing method was ok",
		zap.Int32("userID", userID),
	)

	return &pbSocial.ListFollowingResponse{
		FollowingIds: ids,
	}, nil
}
