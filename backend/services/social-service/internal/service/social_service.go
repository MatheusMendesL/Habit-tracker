package service

import (
	"context"
	pbUser "shared/pb/user"
	AppErr "social/internal/errors"
	"social/internal/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SocialService struct {
	pbUser.UserServiceClient
	repo *repository.SocialRepository
}

func NewSocialService(r *repository.SocialRepository, userClient pbUser.UserServiceClient) *SocialService {
	return &SocialService{
		repo:              r,
		UserServiceClient: userClient,
	}
}

func (s *SocialService) StartFollowing(ctx context.Context, FollowerID, FolloweeID int32) error {
	if FollowerID == FolloweeID {
		return AppErr.ErrInvalidArgument
	}

	_, err := s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: FollowerID})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return AppErr.ErrUserNotFound
		}
		return err
	}

	_, err = s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: FolloweeID})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return AppErr.ErrUserNotFound
		}
		return err
	}

	return s.repo.StartFollowing(ctx, FollowerID, FolloweeID)
}

func (s *SocialService) Unfollow(ctx context.Context, FollowerID, FolloweeID int32) error {
	if FolloweeID == FollowerID {
		return AppErr.ErrInvalidArgument
	}

	_, err := s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: FollowerID})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return AppErr.ErrUserNotFound
		}
		return err
	}

	_, err = s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: FolloweeID})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return AppErr.ErrUserNotFound
		}
		return err
	}

	return s.repo.Unfollow(ctx, FollowerID, FolloweeID)
}

func (s *SocialService) ListFollowers(ctx context.Context, userID int32) ([]int32, error) {
	if userID <= 0 {
		return nil, AppErr.ErrInvalidArgument
	}

	_, err := s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: userID})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, AppErr.ErrUserNotFound
		}
		return nil, err
	}

	return s.repo.ListFollowers(ctx, userID)
}

func (s *SocialService) ListFollowing(ctx context.Context, userID int32) ([]int32, error) {
	if userID <= 0 {
		return nil, AppErr.ErrInvalidArgument
	}

	_, err := s.GetUserByID(ctx, &pbUser.GetUserByIDRequest{UserId: userID})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, AppErr.ErrUserNotFound
		}
		return nil, err
	}

	return s.repo.ListFollowing(ctx, userID)
}
