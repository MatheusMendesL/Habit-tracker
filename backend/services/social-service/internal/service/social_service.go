package service

import (
	"context"
	AppErr "social/internal/errors"
	"social/internal/repository"
)

type SocialService struct {
	repo *repository.SocialRepository
}

func NewSocialService(r *repository.SocialRepository) *SocialService {
	return &SocialService{repo: r}
}

func (s *SocialService) StartFollowing(ctx context.Context, FollowerID, FolloweeID int32) error {
	if FolloweeID == FollowerID {
		return AppErr.ErrInvalidArgument
	}

	return s.repo.StartFollowing(ctx, FollowerID, FolloweeID)
}
