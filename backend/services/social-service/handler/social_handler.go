package handler

import (
	"context"
	"database/sql"
	"errors"
	pbSocial "shared/pb/social"
	pbUser "shared/pb/user"

	"github.com/google/uuid"
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

func (s *SocialHandler) GetUsersByIDs(ctx context.Context, ids []string) (users []*pbUser.User, err error) {
	res, err := s.UserServiceClient.GetUsersByIDs(
		ctx,
		&pbUser.GetUsersByIDsRequest{
			UserIds: ids,
		},
	)

	if err != nil {
		return nil, err
	}

	return res.Users, nil
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

	FollowerID, err := uuid.Parse(req.FollowerId)
	FolloweeID, err := uuid.Parse(req.FolloweeId)

	if FollowerID == uuid.Nil || FolloweeID == uuid.Nil {
		s.logger.Warn("Invalid Users ID",
			zap.String("FollowerID", FollowerID.String()),
			zap.String("FolloweeID", FolloweeID.String()),
		)

		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	err = s.socialService.StartFollowing(ctx, FollowerID, FolloweeID)

	if err != nil {
		s.logger.Error("error to execute StartFollowing method",
			zap.String("FollowerID", FollowerID.String()),
			zap.String("FolloweeID", FolloweeID.String()),
			zap.Error(err),
		)

		return nil, ReceiveErrors(err)
	}

	s.logger.Info("StartFollowing method was ok",
		zap.String("followerID", FollowerID.String()),
		zap.String("followeeID", FolloweeID.String()),
	)

	return &pbSocial.StartFollowingResponse{
		Success: true,
	}, nil
}

func (s *SocialHandler) Unfollow(ctx context.Context, req *pbSocial.UnfollowRequest) (*pbSocial.UnfollowResponse, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	FollowerID, err := uuid.Parse(req.FollowerId)
	FolloweeID, err := uuid.Parse(req.FolloweeId)

	if FollowerID == uuid.Nil || FolloweeID == uuid.Nil {
		s.logger.Warn("Invalid Users ID",
			zap.String("FollowerID", FollowerID.String()),
			zap.String("FolloweeID", FolloweeID.String()),
		)

		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	err = s.socialService.Unfollow(ctx, FollowerID, FolloweeID)

	if err != nil {
		s.logger.Error("error to execute Unfollow method",
			zap.String("FollowerID", FollowerID.String()),
			zap.String("FolloweeID", FolloweeID.String()),
			zap.Error(err),
		)

		return nil, ReceiveErrors(err)
	}

	s.logger.Info("Unfollow method was ok",
		zap.String("FollowerID", FollowerID.String()),
		zap.String("FolloweeID", FolloweeID.String()),
	)

	return &pbSocial.UnfollowResponse{
		Success: true,
	}, nil
}

func (s *SocialHandler) ListFollowers(ctx context.Context, req *pbSocial.ListFollowersRequest) (*pbSocial.ListFollowersResponse, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	userID, err := uuid.Parse(req.UserId)

	if err != nil {
		s.logger.Warn("Invalid User ID",
			zap.String("userID", userID.String()),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	ids, err := s.socialService.ListFollowers(ctx, userID)

	if err != nil {
		s.logger.Error("error to execute ListFollowers method",
			zap.String("userID", userID.String()),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("ListFollowers method was ok",
		zap.String("userID", userID.String()),
	)

	userIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		userIDs = append(userIDs, id.String())
	}

	users, err := s.GetUsersByIDs(ctx, userIDs)

	if err != nil {
		return nil, err
	}

	var socialUsers []*pbSocial.User

	for _, u := range users {
		socialUsers = append(socialUsers, &pbSocial.User{
			Id:    u.Id,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	return &pbSocial.ListFollowersResponse{
		Users: socialUsers,
	}, nil
}

func (s *SocialHandler) ListFollowing(ctx context.Context, req *pbSocial.ListFollowingRequest) (*pbSocial.ListFollowingResponse, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	userID, err := uuid.Parse(req.UserId)

	if err != nil {
		s.logger.Warn("Invalid User ID",
			zap.String("userID", userID.String()),
		)
		return nil, status.Error(codes.InvalidArgument, AppErr.ErrInvalidArgument.Error())
	}

	ids, err := s.socialService.ListFollowing(ctx, userID)

	if err != nil {
		s.logger.Error("error to execute ListFollowing method",
			zap.String("userID", userID.String()),
			zap.Error(err),
		)
		return nil, ReceiveErrors(err)
	}

	s.logger.Info("ListFollowing method was ok",
		zap.String("userID", userID.String()),
	)

	userIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		userIDs = append(userIDs, id.String())
	}

	users, err := s.GetUsersByIDs(ctx, userIDs)

	if err != nil {
		return nil, err
	}

	var socialUsers []*pbSocial.User

	for _, u := range users {
		socialUsers = append(socialUsers, &pbSocial.User{
			Id:    u.Id,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	return &pbSocial.ListFollowingResponse{
		Users: socialUsers,
	}, nil
}
