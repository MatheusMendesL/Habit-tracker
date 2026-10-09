package clients

import (
	"context"
	"fmt"

	"gateway/internal/dto"

	pb "shared/pb/social"
)

type SocialClient struct {
	client pb.SocialServiceClient
}

func NewSocialClient(client pb.SocialServiceClient) *SocialClient {
	return &SocialClient{client: client}
}

func (c *SocialClient) StartFollowing(ctx context.Context, request *dto.FollowRequest) (*dto.FollowResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("invalid request")
	}

	res, err := c.client.StartFollowing(ctx, &pb.StartFollowingRequest{
		FollowerId: request.FollowerID,
		FolloweeId: request.FolloweeID,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from social service")
	}
	return &dto.FollowResponse{Success: res.Success}, nil
}

func (c *SocialClient) Unfollow(ctx context.Context, request *dto.FollowRequest) (*dto.FollowResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("invalid request")
	}

	res, err := c.client.Unfollow(ctx, &pb.UnfollowRequest{
		FollowerId: request.FollowerID,
		FolloweeId: request.FolloweeID,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from social service")
	}
	return &dto.FollowResponse{Success: res.Success}, nil
}

func (c *SocialClient) ListFollowers(ctx context.Context, userID string) (*dto.SocialUsersResponse, error) {
	res, err := c.client.ListFollowers(ctx, &pb.ListFollowersRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from social service")
	}

	users := make([]*dto.User, 0, len(res.Users))
	for _, user := range res.Users {
		if user == nil {
			continue
		}
		users = append(users, &dto.User{
			ID:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		})
	}
	return &dto.SocialUsersResponse{Users: users}, nil
}

func (c *SocialClient) ListFollowing(ctx context.Context, userID string) (*dto.SocialUsersResponse, error) {
	res, err := c.client.ListFollowing(ctx, &pb.ListFollowingRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from social service")
	}

	users := make([]*dto.User, 0, len(res.Users))
	for _, user := range res.Users {
		if user == nil {
			continue
		}
		users = append(users, &dto.User{
			ID:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		})
	}
	return &dto.SocialUsersResponse{Users: users}, nil
}
