package clients

import (
	"context"
	"fmt"

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
		User: &dto.User{
			ID:    res.User.Id,
			Name:  res.User.Name,
			Email: res.User.Email,
		},
	}, nil
}

func (c *UserClient) SearchUsers(ctx context.Context, request *dto.SearchUsersRequest) (*dto.SearchUsersResponse, error) {
	if request == nil {
		request = &dto.SearchUsersRequest{}
	}

	searchReq := &pb.SearchUserRequest{}
	if request.Name != nil && *request.Name != "" {
		name := *request.Name
		searchReq.Name = &name
	}
	if request.Email != nil && *request.Email != "" {
		email := *request.Email
		searchReq.Email = &email
	}

	res, err := c.client.SearchUser(ctx, searchReq)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return &dto.SearchUsersResponse{Users: []*dto.User{}}, nil
	}

	users := make([]*dto.User, 0, len(res.User))
	for _, user := range res.User {
		if user == nil {
			continue
		}
		users = append(users, &dto.User{
			ID:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return &dto.SearchUsersResponse{
		Users: users,
	}, nil
}

func (c *UserClient) DeleteUser(ctx context.Context, request *dto.DeleteUserRequest) (*dto.DeleteUserResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("invalid request")
	}

	res, err := c.client.DeleteUser(ctx, &pb.DeleteUserRequest{UserId: request.ID})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	return &dto.DeleteUserResponse{
		Success: res.Success,
	}, nil
}

func (c *UserClient) GetUsersByIDs(ctx context.Context, request *dto.GetUsersByIDsRequest) (*dto.GetUsersByIDsResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("invalid request")
	}

	res, err := c.client.GetUsersByIDs(ctx, &pb.GetUsersByIDsRequest{UserIds: request.IDs})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return &dto.GetUsersByIDsResponse{Users: []*dto.User{}}, nil
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
	return &dto.GetUsersByIDsResponse{Users: users}, nil
}

func (c *UserClient) EditUser(ctx context.Context, id string, request *dto.EditUserRequest) (*dto.EditUserResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("invalid request")
	}

	res, err := c.client.EditUser(ctx, &pb.EditUserRequest{
		UserId: id,
		Name:   request.Name,
		Email:  request.Email,
	})
	if err != nil {
		return nil, err
	}
	if res == nil || res.User == nil {
		return nil, nil
	}

	return &dto.EditUserResponse{
		User: &dto.User{
			ID:    res.User.Id,
			Name:  res.User.Name,
			Email: res.User.Email,
		},
	}, nil
}

func (c *UserClient) EditPassword(ctx context.Context, id string, request *dto.EditPasswordRequest) (*dto.EditPasswordResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("invalid request")
	}

	res, err := c.client.EditPassword(ctx, &pb.EditPasswordRequest{
		UserId:      id,
		NewPassword: request.NewPassword,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	return &dto.EditPasswordResponse{Success: res.Success}, nil
}
