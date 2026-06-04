package repository

import (
	"context"
	"social/db"

	"github.com/google/uuid"
)

type SocialRepository struct {
	q *db.Queries
}

func NewSocialRepository(q *db.Queries) *SocialRepository {
	return &SocialRepository{q: q}
}

func (r *SocialRepository) StartFollowing(ctx context.Context, FollowerID, FolloweeID uuid.UUID) error {
	params := db.StartFollowingParams{
		FollowerID: FollowerID,
		FolloweeID: FolloweeID,
	}
	err := r.q.StartFollowing(ctx, params)

	if err != nil {
		return err
	}

	return nil
}

func (r *SocialRepository) Unfollow(ctx context.Context, FollowerID, FolloweeID uuid.UUID) error {
	params := db.UnfollowParams{
		FollowerID: FollowerID,
		FolloweeID: FolloweeID,
	}
	err := r.q.Unfollow(ctx, params)

	if err != nil {
		return err
	}

	return nil
}

func (r *SocialRepository) ListFollowers(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	followersIDs, err := r.q.ListFollowers(ctx, userID)

	if err != nil {
		return nil, err
	}

	return followersIDs, nil
}

func (r *SocialRepository) ListFollowing(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	followingIDs, err := r.q.ListFollowing(ctx, userID)
	if err != nil {
		return nil, err
	}

	return followingIDs, nil
}
