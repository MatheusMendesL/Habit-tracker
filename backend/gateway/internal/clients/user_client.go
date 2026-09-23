package clients

import (
	"context"

	"gateway/internal/dto"

	pb "shared/pb/user"
)

type UserClient struct {
	client pb.UserServiceClient
}

func NewUserClient(client pb.UserServiceClient) *UserClient {
	return &UserClient{client: client}
}

func (c *UserClient) GetUserByID(ctx context.Context, request string) (*dto.GetUserByIDResponse, error) {
	res, err := c.client.GetUserByID(ctx, &pb.GetUserByIDRequest{UserId: request})
	if err != nil {
		return nil, err
	}

	if res == nil || res.User == nil {
		return nil, nil
	}

	return &dto.GetUserByIDResponse{
		ID:    res.User.Id,
		Name:  res.User.Name,
		Email: res.User.Email,
	}, nil
}
