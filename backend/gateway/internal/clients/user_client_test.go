package clients

import (
	"context"
	"errors"
	"testing"

	"gateway/internal/dto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "shared/pb/user"
)

type userServiceClientMock struct {
	getUserByID   func(context.Context, *pb.GetUserByIDRequest, ...grpc.CallOption) (*pb.GetUserByIDResponse, error)
	getUsersByIDs func(context.Context, *pb.GetUsersByIDsRequest, ...grpc.CallOption) (*pb.GetUsersByIDsResponse, error)
	editPassword  func(context.Context, *pb.EditPasswordRequest, ...grpc.CallOption) (*pb.EditPasswordResponse, error)
	editUser      func(context.Context, *pb.EditUserRequest, ...grpc.CallOption) (*pb.EditUserResponse, error)
	deleteUser    func(context.Context, *pb.DeleteUserRequest, ...grpc.CallOption) (*pb.DeleteUserResponse, error)
	searchUser    func(context.Context, *pb.SearchUserRequest, ...grpc.CallOption) (*pb.SearchUserResponse, error)
}

func (m *userServiceClientMock) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest, opts ...grpc.CallOption) (*pb.GetUserByIDResponse, error) {
	return m.getUserByID(ctx, req, opts...)
}

func (m *userServiceClientMock) GetUsersByIDs(ctx context.Context, req *pb.GetUsersByIDsRequest, opts ...grpc.CallOption) (*pb.GetUsersByIDsResponse, error) {
	return m.getUsersByIDs(ctx, req, opts...)
}

func (m *userServiceClientMock) EditPassword(ctx context.Context, req *pb.EditPasswordRequest, opts ...grpc.CallOption) (*pb.EditPasswordResponse, error) {
	return m.editPassword(ctx, req, opts...)
}

func (m *userServiceClientMock) EditUser(ctx context.Context, req *pb.EditUserRequest, opts ...grpc.CallOption) (*pb.EditUserResponse, error) {
	return m.editUser(ctx, req, opts...)
}

func (m *userServiceClientMock) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest, opts ...grpc.CallOption) (*pb.DeleteUserResponse, error) {
	return m.deleteUser(ctx, req, opts...)
}

func (m *userServiceClientMock) SearchUser(ctx context.Context, req *pb.SearchUserRequest, opts ...grpc.CallOption) (*pb.SearchUserResponse, error) {
	return m.searchUser(ctx, req, opts...)
}

func TestUserClient_GetUserByID(t *testing.T) {
	client := NewUserClient(&userServiceClientMock{
		getUserByID: func(_ context.Context, req *pb.GetUserByIDRequest, _ ...grpc.CallOption) (*pb.GetUserByIDResponse, error) {
			if req.UserId != "user-id" {
				t.Errorf("GetUserByID() user id = %q, want %q", req.UserId, "user-id")
			}
			return &pb.GetUserByIDResponse{User: &pb.User{Id: "user-id", Name: "Matheus", Email: "matheus@email.com"}}, nil
		},
	})

	response, err := client.GetUserByID(context.Background(), "user-id")
	if err != nil {
		t.Fatalf("GetUserByID() returned unexpected error: %v", err)
	}
	if response == nil || response.User == nil {
		t.Fatal("GetUserByID() returned nil user")
	}
	if response.User.ID != "user-id" || response.User.Name != "Matheus" || response.User.Email != "matheus@email.com" {
		t.Fatalf("GetUserByID() returned wrong user: %#v", response.User)
	}
}

func TestUserClient_GetUserByID_NotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	client := NewUserClient(&userServiceClientMock{
		getUserByID: func(context.Context, *pb.GetUserByIDRequest, ...grpc.CallOption) (*pb.GetUserByIDResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.GetUserByID(context.Background(), "user-id")
	if response != nil {
		t.Fatalf("GetUserByID() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("GetUserByID() error = %v, want %v", err, expectedErr)
	}
}

func TestUserClient_GetUserByID_DatabaseError(t *testing.T) {
	expectedErr := errors.New("database unavailable")
	client := NewUserClient(&userServiceClientMock{
		getUserByID: func(context.Context, *pb.GetUserByIDRequest, ...grpc.CallOption) (*pb.GetUserByIDResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.GetUserByID(context.Background(), "user-id")
	if response != nil {
		t.Fatalf("GetUserByID() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("GetUserByID() error = %v, want %v", err, expectedErr)
	}
}

func TestUserClient_SearchUsers(t *testing.T) {
	name, email := "Matheus", ""
	client := NewUserClient(&userServiceClientMock{
		searchUser: func(_ context.Context, req *pb.SearchUserRequest, _ ...grpc.CallOption) (*pb.SearchUserResponse, error) {
			if req.GetName() != name || req.GetEmail() != "" {
				t.Errorf("SearchUser() filters = (%q, %q), want (%q, %q)", req.GetName(), req.GetEmail(), name, "")
			}
			return &pb.SearchUserResponse{User: []*pb.User{
				{Id: "one", Name: "Matheus", Email: "matheus@email.com"},
				{Id: "two", Name: "Matheus Mendes", Email: "matheus.mendes@email.com"},
			}}, nil
		},
	})

	response, err := client.SearchUsers(context.Background(), &dto.SearchUsersRequest{Name: &name, Email: &email})
	if err != nil {
		t.Fatalf("SearchUsers() returned unexpected error: %v", err)
	}
	if len(response.Users) != 2 {
		t.Fatalf("SearchUsers() returned %d users, want 2", len(response.Users))
	}
	if response.Users[0].Name != "Matheus" {
		t.Fatalf("SearchUsers() first user name = %q, want %q", response.Users[0].Name, "Matheus")
	}
	if response.Users[1].Email != "matheus.mendes@email.com" {
		t.Fatalf("SearchUsers() second user email = %q, want %q", response.Users[1].Email, "matheus.mendes@email.com")
	}
}

func TestUserClient_SearchUsers_EmptyResult(t *testing.T) {
	name := "usuario_inexistente"
	client := NewUserClient(&userServiceClientMock{
		searchUser: func(_ context.Context, req *pb.SearchUserRequest, _ ...grpc.CallOption) (*pb.SearchUserResponse, error) {
			if req.GetName() != name {
				t.Errorf("SearchUser() name = %q, want %q", req.GetName(), name)
			}
			return &pb.SearchUserResponse{User: []*pb.User{}}, nil
		},
	})

	response, err := client.SearchUsers(context.Background(), &dto.SearchUsersRequest{Name: &name})
	if err != nil {
		t.Fatalf("SearchUsers() returned unexpected error: %v", err)
	}
	if response == nil || len(response.Users) != 0 {
		t.Fatalf("SearchUsers() returned %#v, want empty users", response)
	}
}

func TestUserClient_SearchUsers_DatabaseError(t *testing.T) {
	expectedErr := errors.New("database unavailable")
	name := "Matheus"
	client := NewUserClient(&userServiceClientMock{
		searchUser: func(context.Context, *pb.SearchUserRequest, ...grpc.CallOption) (*pb.SearchUserResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.SearchUsers(context.Background(), &dto.SearchUsersRequest{Name: &name})
	if response != nil {
		t.Fatalf("SearchUsers() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("SearchUsers() error = %v, want %v", err, expectedErr)
	}
}

func TestUserClient_EditUser(t *testing.T) {
	userID := "user-id"
	name, email := "Matheus Updated", "matheus.updated@email.com"
	client := NewUserClient(&userServiceClientMock{
		editUser: func(_ context.Context, req *pb.EditUserRequest, _ ...grpc.CallOption) (*pb.EditUserResponse, error) {
			if req.UserId != userID || req.GetName() != name || req.GetEmail() != email {
				t.Errorf("EditUser() request = %#v", req)
			}
			return &pb.EditUserResponse{User: &pb.User{Id: userID, Name: name, Email: email}}, nil
		},
	})

	response, err := client.EditUser(context.Background(), userID, &dto.EditUserRequest{Name: &name, Email: &email})
	if err != nil {
		t.Fatalf("EditUser() returned unexpected error: %v", err)
	}
	if response == nil || response.User == nil {
		t.Fatal("EditUser() returned nil user")
	}
	if response.User.ID != userID || response.User.Name != name || response.User.Email != email {
		t.Fatalf("EditUser() returned wrong user: %#v", response.User)
	}
}

func TestUserClient_EditUser_NotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	client := NewUserClient(&userServiceClientMock{
		editUser: func(context.Context, *pb.EditUserRequest, ...grpc.CallOption) (*pb.EditUserResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.EditUser(context.Background(), "user-id", &dto.EditUserRequest{})
	if response != nil {
		t.Fatalf("EditUser() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("EditUser() error = %v, want %v", err, expectedErr)
	}
}

func TestUserClient_EditUser_DatabaseError(t *testing.T) {
	expectedErr := errors.New("database unavailable")
	client := NewUserClient(&userServiceClientMock{
		editUser: func(context.Context, *pb.EditUserRequest, ...grpc.CallOption) (*pb.EditUserResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.EditUser(context.Background(), "user-id", &dto.EditUserRequest{})
	if response != nil {
		t.Fatalf("EditUser() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("EditUser() error = %v, want %v", err, expectedErr)
	}
}

func TestUserClient_DeleteUser(t *testing.T) {
	userID := "user-id"
	client := NewUserClient(&userServiceClientMock{
		deleteUser: func(_ context.Context, req *pb.DeleteUserRequest, _ ...grpc.CallOption) (*pb.DeleteUserResponse, error) {
			if req.UserId != userID {
				t.Errorf("DeleteUser() user id = %q, want %q", req.UserId, userID)
			}
			return &pb.DeleteUserResponse{Success: true}, nil
		},
	})

	response, err := client.DeleteUser(context.Background(), &dto.DeleteUserRequest{ID: userID})
	if err != nil {
		t.Fatalf("DeleteUser() returned unexpected error: %v", err)
	}
	if response == nil || !response.Success {
		t.Fatalf("DeleteUser() returned %#v, want success", response)
	}
}

func TestUserClient_DeleteUser_NotFound(t *testing.T) {
	userID := "missing-user"
	client := NewUserClient(&userServiceClientMock{
		deleteUser: func(_ context.Context, req *pb.DeleteUserRequest, _ ...grpc.CallOption) (*pb.DeleteUserResponse, error) {
			if req.UserId != userID {
				t.Errorf("DeleteUser() user id = %q, want %q", req.UserId, userID)
			}
			return &pb.DeleteUserResponse{Success: false}, nil
		},
	})

	response, err := client.DeleteUser(context.Background(), &dto.DeleteUserRequest{ID: userID})
	if err != nil {
		t.Fatalf("DeleteUser() returned unexpected error: %v", err)
	}
	if response == nil || response.Success {
		t.Fatalf("DeleteUser() returned %#v, want unsuccessful response", response)
	}
}

func TestUserClient_DeleteUser_DatabaseError(t *testing.T) {
	expectedErr := errors.New("database unavailable")
	client := NewUserClient(&userServiceClientMock{
		deleteUser: func(context.Context, *pb.DeleteUserRequest, ...grpc.CallOption) (*pb.DeleteUserResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.DeleteUser(context.Background(), &dto.DeleteUserRequest{ID: "user-id"})
	if response != nil {
		t.Fatalf("DeleteUser() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("DeleteUser() error = %v, want %v", err, expectedErr)
	}
}

func TestUserClient_EditPassword(t *testing.T) {
	userID, password := "user-id", "new-secure-password"
	client := NewUserClient(&userServiceClientMock{
		editPassword: func(_ context.Context, req *pb.EditPasswordRequest, _ ...grpc.CallOption) (*pb.EditPasswordResponse, error) {
			if req.UserId != userID || req.NewPassword != password {
				t.Errorf("EditPassword() request = %#v", req)
			}
			return &pb.EditPasswordResponse{Success: true}, nil
		},
	})

	response, err := client.EditPassword(context.Background(), userID, &dto.EditPasswordRequest{NewPassword: password})
	if err != nil {
		t.Fatalf("EditPassword() returned unexpected error: %v", err)
	}
	if response == nil || !response.Success {
		t.Fatalf("EditPassword() returned %#v, want success", response)
	}
}

func TestUserClient_EditPassword_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database unavailable")
	client := NewUserClient(&userServiceClientMock{
		editPassword: func(context.Context, *pb.EditPasswordRequest, ...grpc.CallOption) (*pb.EditPasswordResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.EditPassword(context.Background(), "user-id", &dto.EditPasswordRequest{NewPassword: "new-secure-password"})
	if response != nil {
		t.Fatalf("EditPassword() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("EditPassword() error = %v, want %v", err, expectedErr)
	}
}

func TestUserClient_GetUsersByIDs(t *testing.T) {
	ids := []string{"one", "two"}
	client := NewUserClient(&userServiceClientMock{
		getUsersByIDs: func(_ context.Context, req *pb.GetUsersByIDsRequest, _ ...grpc.CallOption) (*pb.GetUsersByIDsResponse, error) {
			if len(req.UserIds) != len(ids) || req.UserIds[0] != ids[0] || req.UserIds[1] != ids[1] {
				t.Errorf("GetUsersByIDs() ids = %#v, want %#v", req.UserIds, ids)
			}
			return &pb.GetUsersByIDsResponse{Users: []*pb.User{
				{Id: ids[0], Name: "Matheus", Email: "matheus@email.com"},
				{Id: ids[1], Name: "Mendes", Email: "mendes@email.com"},
			}}, nil
		},
	})

	response, err := client.GetUsersByIDs(context.Background(), &dto.GetUsersByIDsRequest{IDs: ids})
	if err != nil {
		t.Fatalf("GetUsersByIDs() returned unexpected error: %v", err)
	}
	if len(response.Users) != 2 {
		t.Fatalf("GetUsersByIDs() returned %d users, want 2", len(response.Users))
	}
	if response.Users[0].Name != "Matheus" {
		t.Fatalf("GetUsersByIDs() first user name = %q, want %q", response.Users[0].Name, "Matheus")
	}
	if response.Users[1].Email != "mendes@email.com" {
		t.Fatalf("GetUsersByIDs() second user email = %q, want %q", response.Users[1].Email, "mendes@email.com")
	}
}

func TestUserClient_GetUsersByIDs_NotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	client := NewUserClient(&userServiceClientMock{
		getUsersByIDs: func(context.Context, *pb.GetUsersByIDsRequest, ...grpc.CallOption) (*pb.GetUsersByIDsResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.GetUsersByIDs(context.Background(), &dto.GetUsersByIDsRequest{IDs: []string{"missing-user"}})
	if response != nil {
		t.Fatalf("GetUsersByIDs() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("GetUsersByIDs() error = %v, want %v", err, expectedErr)
	}
}

func TestUserClient_GetUsersByIDs_DatabaseError(t *testing.T) {
	expectedErr := errors.New("database unavailable")
	client := NewUserClient(&userServiceClientMock{
		getUsersByIDs: func(context.Context, *pb.GetUsersByIDsRequest, ...grpc.CallOption) (*pb.GetUsersByIDsResponse, error) {
			return nil, expectedErr
		},
	})

	response, err := client.GetUsersByIDs(context.Background(), &dto.GetUsersByIDsRequest{IDs: []string{"user-id"}})
	if response != nil {
		t.Fatalf("GetUsersByIDs() returned %#v, want nil", response)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("GetUsersByIDs() error = %v, want %v", err, expectedErr)
	}
}
