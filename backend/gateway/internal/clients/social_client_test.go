package clients

import (
	"context"
	"errors"
	"testing"

	"gateway/internal/dto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "shared/pb/social"
)

type socialServiceClientMock struct {
	startFollowing func(context.Context, *pb.StartFollowingRequest, ...grpc.CallOption) (*pb.StartFollowingResponse, error)
	unfollow       func(context.Context, *pb.UnfollowRequest, ...grpc.CallOption) (*pb.UnfollowResponse, error)
	listFollowers  func(context.Context, *pb.ListFollowersRequest, ...grpc.CallOption) (*pb.ListFollowersResponse, error)
	listFollowing  func(context.Context, *pb.ListFollowingRequest, ...grpc.CallOption) (*pb.ListFollowingResponse, error)
}

func (m *socialServiceClientMock) StartFollowing(ctx context.Context, req *pb.StartFollowingRequest, opts ...grpc.CallOption) (*pb.StartFollowingResponse, error) {
	return m.startFollowing(ctx, req, opts...)
}

func (m *socialServiceClientMock) Unfollow(ctx context.Context, req *pb.UnfollowRequest, opts ...grpc.CallOption) (*pb.UnfollowResponse, error) {
	return m.unfollow(ctx, req, opts...)
}

func (m *socialServiceClientMock) ListFollowers(ctx context.Context, req *pb.ListFollowersRequest, opts ...grpc.CallOption) (*pb.ListFollowersResponse, error) {
	return m.listFollowers(ctx, req, opts...)
}

func (m *socialServiceClientMock) ListFollowing(ctx context.Context, req *pb.ListFollowingRequest, opts ...grpc.CallOption) (*pb.ListFollowingResponse, error) {
	return m.listFollowing(ctx, req, opts...)
}

func TestSocialClient_StartFollowing(t *testing.T) {
	request := &dto.FollowRequest{FollowerID: "follower-id", FolloweeID: "followee-id"}
	client := NewSocialClient(&socialServiceClientMock{
		startFollowing: func(_ context.Context, req *pb.StartFollowingRequest, _ ...grpc.CallOption) (*pb.StartFollowingResponse, error) {
			if req.FollowerId != request.FollowerID || req.FolloweeId != request.FolloweeID {
				t.Errorf("StartFollowing() request = %#v", req)
			}
			return &pb.StartFollowingResponse{Success: true}, nil
		},
	})

	response, err := client.StartFollowing(context.Background(), request)
	if err != nil {
		t.Fatalf("StartFollowing() returned unexpected error: %v", err)
	}
	if response == nil || !response.Success {
		t.Fatalf("StartFollowing() returned %#v, want success", response)
	}
}

func TestSocialClient_StartFollowing_SameUser(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "users cannot follow themselves")
	client := NewSocialClient(&socialServiceClientMock{
		startFollowing: func(_ context.Context, req *pb.StartFollowingRequest, _ ...grpc.CallOption) (*pb.StartFollowingResponse, error) {
			if req.FollowerId != req.FolloweeId {
				t.Errorf("StartFollowing() expected same ids, got %q and %q", req.FollowerId, req.FolloweeId)
			}
			return nil, expectedErr
		},
	})

	response, err := client.StartFollowing(context.Background(), &dto.FollowRequest{FollowerID: "same-id", FolloweeID: "same-id"})
	if response != nil {
		t.Fatalf("StartFollowing() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("StartFollowing() error = %v, want %v", err, expectedErr)
	}
}

func TestSocialClient_StartFollowing_FollowerUserNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	client := NewSocialClient(&socialServiceClientMock{
		startFollowing: func(context.Context, *pb.StartFollowingRequest, ...grpc.CallOption) (*pb.StartFollowingResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.StartFollowing(context.Background(), &dto.FollowRequest{FollowerID: "missing-follower", FolloweeID: "followee-id"})
	if response != nil {
		t.Fatalf("StartFollowing() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("StartFollowing() error = %v, want %v", err, expectedErr)
	}
}

func TestSocialClient_StartFollowing_FolloweeUserNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	client := NewSocialClient(&socialServiceClientMock{
		startFollowing: func(context.Context, *pb.StartFollowingRequest, ...grpc.CallOption) (*pb.StartFollowingResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.StartFollowing(context.Background(), &dto.FollowRequest{FollowerID: "follower-id", FolloweeID: "missing-followee"})
	if response != nil {
		t.Fatalf("StartFollowing() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("StartFollowing() error = %v, want %v", err, expectedErr)
	}
}

func TestSocialClient_StartFollowing_DatabaseError(t *testing.T) {
	expectedErr := errors.New("database failure")
	client := NewSocialClient(&socialServiceClientMock{
		startFollowing: func(context.Context, *pb.StartFollowingRequest, ...grpc.CallOption) (*pb.StartFollowingResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.StartFollowing(context.Background(), &dto.FollowRequest{FollowerID: "follower-id", FolloweeID: "followee-id"})
	if response != nil {
		t.Fatalf("StartFollowing() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("StartFollowing() error = %v, want %v", err, expectedErr)
	}
}

func TestSocialClient_Unfollow(t *testing.T) {
	request := &dto.FollowRequest{FollowerID: "follower-id", FolloweeID: "followee-id"}
	client := NewSocialClient(&socialServiceClientMock{
		unfollow: func(_ context.Context, req *pb.UnfollowRequest, _ ...grpc.CallOption) (*pb.UnfollowResponse, error) {
			if req.FollowerId != request.FollowerID || req.FolloweeId != request.FolloweeID {
				t.Errorf("Unfollow() request = %#v", req)
			}
			return &pb.UnfollowResponse{Success: true}, nil
		},
	})

	response, err := client.Unfollow(context.Background(), request)
	if err != nil {
		t.Fatalf("Unfollow() returned unexpected error: %v", err)
	}
	if response == nil || !response.Success {
		t.Fatalf("Unfollow() returned %#v, want success", response)
	}
}

func TestSocialClient_Unfollow_SameUser(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "users cannot unfollow themselves")
	client := NewSocialClient(&socialServiceClientMock{
		unfollow: func(_ context.Context, req *pb.UnfollowRequest, _ ...grpc.CallOption) (*pb.UnfollowResponse, error) {
			if req.FollowerId != req.FolloweeId {
				t.Errorf("Unfollow() expected same ids, got %q and %q", req.FollowerId, req.FolloweeId)
			}
			return nil, expectedErr
		},
	})

	response, err := client.Unfollow(context.Background(), &dto.FollowRequest{FollowerID: "same-id", FolloweeID: "same-id"})
	if response != nil {
		t.Fatalf("Unfollow() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Unfollow() error = %v, want %v", err, expectedErr)
	}
}

func TestSocialClient_Unfollow_UserNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	client := NewSocialClient(&socialServiceClientMock{
		unfollow: func(context.Context, *pb.UnfollowRequest, ...grpc.CallOption) (*pb.UnfollowResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.Unfollow(context.Background(), &dto.FollowRequest{FollowerID: "missing-id", FolloweeID: "followee-id"})
	if response != nil {
		t.Fatalf("Unfollow() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Unfollow() error = %v, want %v", err, expectedErr)
	}
}

func TestSocialClient_Unfollow_DatabaseError(t *testing.T) {
	expectedErr := errors.New("database failure")
	client := NewSocialClient(&socialServiceClientMock{
		unfollow: func(context.Context, *pb.UnfollowRequest, ...grpc.CallOption) (*pb.UnfollowResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.Unfollow(context.Background(), &dto.FollowRequest{FollowerID: "follower-id", FolloweeID: "followee-id"})
	if response != nil {
		t.Fatalf("Unfollow() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Unfollow() error = %v, want %v", err, expectedErr)
	}
}

func TestSocialClient_ListFollowers(t *testing.T) {
	userID := "user-id"
	client := NewSocialClient(&socialServiceClientMock{
		listFollowers: func(_ context.Context, req *pb.ListFollowersRequest, _ ...grpc.CallOption) (*pb.ListFollowersResponse, error) {
			if req.UserId != userID {
				t.Errorf("ListFollowers() user id = %q, want %q", req.UserId, userID)
			}
			return &pb.ListFollowersResponse{Users: []*pb.User{
				{Id: "one", Name: "Matheus", Email: "matheus@email.com"},
				{Id: "two", Name: "Mendes", Email: "mendes@email.com"},
			}}, nil
		},
	})

	response, err := client.ListFollowers(context.Background(), userID)
	if err != nil {
		t.Fatalf("ListFollowers() returned unexpected error: %v", err)
	}
	if len(response.Users) != 2 {
		t.Fatalf("ListFollowers() returned %d users, want 2", len(response.Users))
	}
	if response.Users[0].ID != "one" || response.Users[0].Name != "Matheus" {
		t.Fatalf("ListFollowers() first user = %#v, want Matheus", response.Users[0])
	}
	if response.Users[1].ID != "two" || response.Users[1].Email != "mendes@email.com" {
		t.Fatalf("ListFollowers() second user = %#v, want Mendes", response.Users[1])
	}
}

func TestSocialClient_ListFollowers_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid user id")
	client := NewSocialClient(&socialServiceClientMock{
		listFollowers: func(_ context.Context, req *pb.ListFollowersRequest, _ ...grpc.CallOption) (*pb.ListFollowersResponse, error) {
			if req.UserId != "" {
				t.Errorf("ListFollowers() user id = %q, want empty", req.UserId)
			}
			return nil, expectedErr
		},
	})

	response, err := client.ListFollowers(context.Background(), "")
	if response != nil {
		t.Fatalf("ListFollowers() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("ListFollowers() error = %v, want %v", err, expectedErr)
	}
}

func TestSocialClient_ListFollowers_UserNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	client := NewSocialClient(&socialServiceClientMock{
		listFollowers: func(context.Context, *pb.ListFollowersRequest, ...grpc.CallOption) (*pb.ListFollowersResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.ListFollowers(context.Background(), "missing-user")
	if response != nil {
		t.Fatalf("ListFollowers() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("ListFollowers() error = %v, want %v", err, expectedErr)
	}
}

func TestSocialClient_ListFollowers_DatabaseError(t *testing.T) {
	expectedErr := errors.New("database failure")
	client := NewSocialClient(&socialServiceClientMock{
		listFollowers: func(context.Context, *pb.ListFollowersRequest, ...grpc.CallOption) (*pb.ListFollowersResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.ListFollowers(context.Background(), "user-id")
	if response != nil {
		t.Fatalf("ListFollowers() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("ListFollowers() error = %v, want %v", err, expectedErr)
	}
}

func TestSocialClient_ListFollowing(t *testing.T) {
	userID := "user-id"
	client := NewSocialClient(&socialServiceClientMock{
		listFollowing: func(_ context.Context, req *pb.ListFollowingRequest, _ ...grpc.CallOption) (*pb.ListFollowingResponse, error) {
			if req.UserId != userID {
				t.Errorf("ListFollowing() user id = %q, want %q", req.UserId, userID)
			}
			return &pb.ListFollowingResponse{Users: []*pb.User{
				{Id: "one", Name: "Matheus", Email: "matheus@email.com"},
				{Id: "two", Name: "Mendes", Email: "mendes@email.com"},
			}}, nil
		},
	})

	response, err := client.ListFollowing(context.Background(), userID)
	if err != nil {
		t.Fatalf("ListFollowing() returned unexpected error: %v", err)
	}
	if len(response.Users) != 2 {
		t.Fatalf("ListFollowing() returned %d users, want 2", len(response.Users))
	}
	if response.Users[0].ID != "one" || response.Users[0].Name != "Matheus" {
		t.Fatalf("ListFollowing() first user = %#v, want Matheus", response.Users[0])
	}
	if response.Users[1].ID != "two" || response.Users[1].Email != "mendes@email.com" {
		t.Fatalf("ListFollowing() second user = %#v, want Mendes", response.Users[1])
	}
}

func TestSocialClient_ListFollowing_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid user id")
	client := NewSocialClient(&socialServiceClientMock{
		listFollowing: func(_ context.Context, req *pb.ListFollowingRequest, _ ...grpc.CallOption) (*pb.ListFollowingResponse, error) {
			if req.UserId != "" {
				t.Errorf("ListFollowing() user id = %q, want empty", req.UserId)
			}
			return nil, expectedErr
		},
	})

	response, err := client.ListFollowing(context.Background(), "")
	if response != nil {
		t.Fatalf("ListFollowing() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("ListFollowing() error = %v, want %v", err, expectedErr)
	}
}

func TestSocialClient_ListFollowing_UserNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	client := NewSocialClient(&socialServiceClientMock{
		listFollowing: func(context.Context, *pb.ListFollowingRequest, ...grpc.CallOption) (*pb.ListFollowingResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.ListFollowing(context.Background(), "missing-user")
	if response != nil {
		t.Fatalf("ListFollowing() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("ListFollowing() error = %v, want %v", err, expectedErr)
	}
}

func TestSocialClient_ListFollowing_DatabaseError(t *testing.T) {
	expectedErr := errors.New("database failure")
	client := NewSocialClient(&socialServiceClientMock{
		listFollowing: func(context.Context, *pb.ListFollowingRequest, ...grpc.CallOption) (*pb.ListFollowingResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.ListFollowing(context.Background(), "user-id")
	if response != nil {
		t.Fatalf("ListFollowing() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("ListFollowing() error = %v, want %v", err, expectedErr)
	}
}
