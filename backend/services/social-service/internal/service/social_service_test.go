package service

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbUser "shared/pb/user"
	"social/db"
	AppErr "social/internal/errors"
	"social/internal/repository"
)

type mockSocialUserServiceClient struct {
	userByIDErrs []error
	userByIDIdx  int
}

func (m *mockSocialUserServiceClient) nextUserByIDErr() error {
	if len(m.userByIDErrs) == 0 {
		return nil
	}
	if m.userByIDIdx >= len(m.userByIDErrs) {
		return nil
	}
	err := m.userByIDErrs[m.userByIDIdx]
	m.userByIDIdx++
	return err
}

func (m *mockSocialUserServiceClient) GetUserByID(ctx context.Context, in *pbUser.GetUserByIDRequest, opts ...grpc.CallOption) (*pbUser.GetUserByIDResponse, error) {
	if err := m.nextUserByIDErr(); err != nil {
		return nil, err
	}
	return &pbUser.GetUserByIDResponse{User: &pbUser.User{Id: in.UserId, Name: "Teste", Email: "teste@email.com"}}, nil
}

func (m *mockSocialUserServiceClient) GetUsersByIDs(ctx context.Context, in *pbUser.GetUsersByIDsRequest, opts ...grpc.CallOption) (*pbUser.GetUsersByIDsResponse, error) {
	return &pbUser.GetUsersByIDsResponse{}, nil
}

func (m *mockSocialUserServiceClient) EditPassword(ctx context.Context, in *pbUser.EditPasswordRequest, opts ...grpc.CallOption) (*pbUser.EditPasswordResponse, error) {
	return &pbUser.EditPasswordResponse{}, nil
}

func (m *mockSocialUserServiceClient) EditUser(ctx context.Context, in *pbUser.EditUserRequest, opts ...grpc.CallOption) (*pbUser.EditUserResponse, error) {
	return &pbUser.EditUserResponse{}, nil
}

func (m *mockSocialUserServiceClient) DeleteUser(ctx context.Context, in *pbUser.DeleteUserRequest, opts ...grpc.CallOption) (*pbUser.DeleteUserResponse, error) {
	return &pbUser.DeleteUserResponse{}, nil
}

func (m *mockSocialUserServiceClient) SearchUser(ctx context.Context, in *pbUser.SearchUserRequest, opts ...grpc.CallOption) (*pbUser.SearchUserResponse, error) {
	return &pbUser.SearchUserResponse{}, nil
}

func TestSocialService_StartFollowing(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	followerID := uuid.New()
	followeeID := uuid.New()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO follows (follower_id, followee_id)\nVALUES ($1, $2)\n")).
		WithArgs(followerID, followeeID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	queries := db.New(mockDB)
	repo := repository.NewSocialRepository(queries)
	service := NewSocialService(repo, &mockSocialUserServiceClient{})

	if err := service.StartFollowing(context.Background(), followerID, followeeID); err != nil {
		t.Fatalf("StartFollowing() returned unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestSocialService_StartFollowing_SameUser(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	queries := db.New(mockDB)
	repo := repository.NewSocialRepository(queries)
	service := NewSocialService(repo, &mockSocialUserServiceClient{})

	err = service.StartFollowing(context.Background(), userID, userID)
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("StartFollowing() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestSocialService_StartFollowing_FollowerUserNotFound(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	followerID := uuid.New()
	followeeID := uuid.New()
	queries := db.New(mockDB)
	repo := repository.NewSocialRepository(queries)
	userClient := &mockSocialUserServiceClient{userByIDErrs: []error{status.Error(codes.NotFound, AppErr.ErrUserNotFound.Error())}}
	service := NewSocialService(repo, userClient)

	err = service.StartFollowing(context.Background(), followerID, followeeID)
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("StartFollowing() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}
}

func TestSocialService_StartFollowing_FolloweeUserNotFound(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	followerID := uuid.New()
	followeeID := uuid.New()
	queries := db.New(mockDB)
	repo := repository.NewSocialRepository(queries)
	userClient := &mockSocialUserServiceClient{userByIDErrs: []error{nil, status.Error(codes.NotFound, AppErr.ErrUserNotFound.Error())}}
	service := NewSocialService(repo, userClient)

	err = service.StartFollowing(context.Background(), followerID, followeeID)
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("StartFollowing() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}
}

func TestSocialService_StartFollowing_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	followerID := uuid.New()
	followeeID := uuid.New()
	expectedErr := errors.New("database failure")
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO follows (follower_id, followee_id)\nVALUES ($1, $2)\n")).
		WithArgs(followerID, followeeID).
		WillReturnError(expectedErr)

	queries := db.New(mockDB)
	repo := repository.NewSocialRepository(queries)
	service := NewSocialService(repo, &mockSocialUserServiceClient{})

	err = service.StartFollowing(context.Background(), followerID, followeeID)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("StartFollowing() error = %v, want %v", err, expectedErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestSocialService_Unfollow(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	followerID := uuid.New()
	followeeID := uuid.New()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM follows\nWHERE follower_id = $1\n  AND followee_id = $2\n")).
		WithArgs(followerID, followeeID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	queries := db.New(mockDB)
	repo := repository.NewSocialRepository(queries)
	service := NewSocialService(repo, &mockSocialUserServiceClient{})

	if err := service.Unfollow(context.Background(), followerID, followeeID); err != nil {
		t.Fatalf("Unfollow() returned unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestSocialService_Unfollow_SameUser(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	queries := db.New(mockDB)
	repo := repository.NewSocialRepository(queries)
	service := NewSocialService(repo, &mockSocialUserServiceClient{})

	err = service.Unfollow(context.Background(), userID, userID)
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("Unfollow() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestSocialService_Unfollow_UserNotFound(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	followerID := uuid.New()
	followeeID := uuid.New()
	queries := db.New(mockDB)
	repo := repository.NewSocialRepository(queries)
	userClient := &mockSocialUserServiceClient{userByIDErrs: []error{status.Error(codes.NotFound, AppErr.ErrUserNotFound.Error())}}
	service := NewSocialService(repo, userClient)

	err = service.Unfollow(context.Background(), followerID, followeeID)
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("Unfollow() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}
}

func TestSocialService_Unfollow_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	followerID := uuid.New()
	followeeID := uuid.New()
	expectedErr := errors.New("database failure")
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM follows\nWHERE follower_id = $1\n  AND followee_id = $2\n")).
		WithArgs(followerID, followeeID).
		WillReturnError(expectedErr)

	queries := db.New(mockDB)
	repo := repository.NewSocialRepository(queries)
	service := NewSocialService(repo, &mockSocialUserServiceClient{})

	err = service.Unfollow(context.Background(), followerID, followeeID)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Unfollow() error = %v, want %v", err, expectedErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestSocialService_ListFollowers(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	followerOne := uuid.New()
	followerTwo := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT follower_id\nFROM follows\nWHERE followee_id = $1\n")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"follower_id"}).
				AddRow(followerOne).
				AddRow(followerTwo),
		)

	queries := db.New(mockDB)
	repo := repository.NewSocialRepository(queries)
	service := NewSocialService(repo, &mockSocialUserServiceClient{})

	followers, err := service.ListFollowers(context.Background(), userID)
	if err != nil {
		t.Fatalf("ListFollowers() returned unexpected error: %v", err)
	}
	if len(followers) != 2 {
		t.Fatalf("ListFollowers() returned %d followers, want 2", len(followers))
	}
	if followers[0] != followerOne || followers[1] != followerTwo {
		t.Fatalf("ListFollowers() response mismatch: got %v, want %v and %v", followers, followerOne, followerTwo)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestSocialService_ListFollowers_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	queries := db.New(mockDB)
	repo := repository.NewSocialRepository(queries)
	service := NewSocialService(repo, &mockSocialUserServiceClient{})

	followers, err := service.ListFollowers(context.Background(), uuid.Nil)
	if followers != nil {
		t.Fatalf("ListFollowers() expected nil slice, got %#v", followers)
	}
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("ListFollowers() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestSocialService_ListFollowers_UserNotFound(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	queries := db.New(mockDB)
	repo := repository.NewSocialRepository(queries)
	userClient := &mockSocialUserServiceClient{userByIDErrs: []error{status.Error(codes.NotFound, AppErr.ErrUserNotFound.Error())}}
	service := NewSocialService(repo, userClient)

	followers, err := service.ListFollowers(context.Background(), userID)
	if followers != nil {
		t.Fatalf("ListFollowers() expected nil slice, got %#v", followers)
	}
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("ListFollowers() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}
}

func TestSocialService_ListFollowers_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	expectedErr := errors.New("database failure")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT follower_id\nFROM follows\nWHERE followee_id = $1\n")).
		WithArgs(userID).
		WillReturnError(expectedErr)

	queries := db.New(mockDB)
	repo := repository.NewSocialRepository(queries)
	service := NewSocialService(repo, &mockSocialUserServiceClient{})

	followers, err := service.ListFollowers(context.Background(), userID)
	if followers != nil {
		t.Fatalf("ListFollowers() expected nil slice, got %#v", followers)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("ListFollowers() error = %v, want %v", err, expectedErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestSocialService_ListFollowing(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	followeeOne := uuid.New()
	followeeTwo := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT followee_id\nFROM follows\nWHERE follower_id = $1\n")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"followee_id"}).
				AddRow(followeeOne).
				AddRow(followeeTwo),
		)

	queries := db.New(mockDB)
	repo := repository.NewSocialRepository(queries)
	service := NewSocialService(repo, &mockSocialUserServiceClient{})

	following, err := service.ListFollowing(context.Background(), userID)
	if err != nil {
		t.Fatalf("ListFollowing() returned unexpected error: %v", err)
	}
	if len(following) != 2 {
		t.Fatalf("ListFollowing() returned %d following, want 2", len(following))
	}
	if following[0] != followeeOne || following[1] != followeeTwo {
		t.Fatalf("ListFollowing() response mismatch: got %v, want %v and %v", following, followeeOne, followeeTwo)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestSocialService_ListFollowing_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	queries := db.New(mockDB)
	repo := repository.NewSocialRepository(queries)
	service := NewSocialService(repo, &mockSocialUserServiceClient{})

	following, err := service.ListFollowing(context.Background(), uuid.Nil)
	if following != nil {
		t.Fatalf("ListFollowing() expected nil slice, got %#v", following)
	}
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("ListFollowing() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestSocialService_ListFollowing_UserNotFound(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	queries := db.New(mockDB)
	repo := repository.NewSocialRepository(queries)
	userClient := &mockSocialUserServiceClient{userByIDErrs: []error{status.Error(codes.NotFound, AppErr.ErrUserNotFound.Error())}}
	service := NewSocialService(repo, userClient)

	following, err := service.ListFollowing(context.Background(), userID)
	if following != nil {
		t.Fatalf("ListFollowing() expected nil slice, got %#v", following)
	}
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("ListFollowing() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}
}

func TestSocialService_ListFollowing_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	expectedErr := errors.New("database failure")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT followee_id\nFROM follows\nWHERE follower_id = $1\n")).
		WithArgs(userID).
		WillReturnError(expectedErr)

	queries := db.New(mockDB)
	repo := repository.NewSocialRepository(queries)
	service := NewSocialService(repo, &mockSocialUserServiceClient{})

	following, err := service.ListFollowing(context.Background(), userID)
	if following != nil {
		t.Fatalf("ListFollowing() expected nil slice, got %#v", following)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("ListFollowing() error = %v, want %v", err, expectedErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}
