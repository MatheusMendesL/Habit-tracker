package service

import (
	"context"
	"database/sql"
	"errors"
	"habit-service/db"
	AppErr "habit-service/internal/errors"
	"habit-service/internal/repository"
	"regexp"
	pbUser "shared/pb/user"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type mockUserServiceClient struct{}

func (m *mockUserServiceClient) GetUserByID(ctx context.Context, in *pbUser.GetUserByIDRequest, opts ...grpc.CallOption) (*pbUser.GetUserByIDResponse, error) {
	return &pbUser.GetUserByIDResponse{User: &pbUser.User{Id: in.UserId, Name: "Teste", Email: "teste@email.com"}}, nil
}

func (m *mockUserServiceClient) GetUsersByIDs(ctx context.Context, in *pbUser.GetUsersByIDsRequest, opts ...grpc.CallOption) (*pbUser.GetUsersByIDsResponse, error) {
	return &pbUser.GetUsersByIDsResponse{}, nil
}

func (m *mockUserServiceClient) EditPassword(ctx context.Context, in *pbUser.EditPasswordRequest, opts ...grpc.CallOption) (*pbUser.EditPasswordResponse, error) {
	return &pbUser.EditPasswordResponse{}, nil
}

func (m *mockUserServiceClient) EditUser(ctx context.Context, in *pbUser.EditUserRequest, opts ...grpc.CallOption) (*pbUser.EditUserResponse, error) {
	return &pbUser.EditUserResponse{}, nil
}

func (m *mockUserServiceClient) DeleteUser(ctx context.Context, in *pbUser.DeleteUserRequest, opts ...grpc.CallOption) (*pbUser.DeleteUserResponse, error) {
	return &pbUser.DeleteUserResponse{}, nil
}

func (m *mockUserServiceClient) SearchUser(ctx context.Context, in *pbUser.SearchUserRequest, opts ...grpc.CallOption) (*pbUser.SearchUserResponse, error) {
	return &pbUser.SearchUserResponse{}, nil
}

func TestHabitService_GetHabitByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).
				AddRow(habitID, userID, "habito 1", "descrição", "teste.png", createdAt),
		)

	queries := db.New(mockDB)
	habitRepo := repository.NewHabitRepository(queries)
	routineRepo := repository.NewRoutineRepository(queries)
	userClient := &mockUserServiceClient{}
	service := NewHabitService(habitRepo, routineRepo, userClient)

	habit, err := service.GetHabitByID(context.Background(), habitID)
	if err != nil {
		t.Fatalf("GetHabitByID() returned unexpected error: %v", err)
	}
	if habit.ID != habitID {
		t.Fatalf("GetHabitByID() returned wrong habit ID: got %s, want %s", habit.ID, habitID)
	}
	if habit.UserID != userID {
		t.Fatalf("GetHabitByID() returned wrong user ID: got %s, want %s", habit.UserID, userID)
	}
	if habit.Name != "habito 1" {
		t.Fatalf("GetHabitByID() returned wrong habit name: got %s, want %s", habit.Name, "habito 1")
	}
	if habit.Description.String != "descrição" {
		t.Fatalf("GetHabitByID() returned wrong description: got %s, want %s", habit.Description.String, "descrição")
	}
	if habit.ImageUrl.String != "teste.png" {
		t.Fatalf("GetHabitByID() returned wrong image URL: got %s, want %s", habit.ImageUrl.String, "teste.png")
	}
	if !habit.CreatedAt.Equal(createdAt) {
		t.Fatalf("GetHabitByID() returned wrong created_at: got %v, want %v", habit.CreatedAt, createdAt)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_GetHabitByID_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnError(sql.ErrNoRows)

	queries := db.New(mockDB)
	habitRepo := repository.NewHabitRepository(queries)
	routineRepo := repository.NewRoutineRepository(queries)
	userClient := &mockUserServiceClient{}
	service := NewHabitService(habitRepo, routineRepo, userClient)

	habit, err := service.GetHabitByID(context.Background(), habitID)

	// isso é pra ligar a struct e validar caso ela n seja nil
	if habit != (db.Habit{}) {
		t.Fatalf("GetHabitByID() expected empty habit struct, got %#v", habit)
	}

	if !errors.Is(err, AppErr.ErrHabitNotFound) {
		t.Fatalf("GetHabitByID() error = %v, want %v", err, AppErr.ErrHabitNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}

}

func TestHabitService_GetHabitByID_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}

	defer mockDB.Close()

	habitID := uuid.New()
	dbErr := errors.New("fatal database connection failure")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnError(dbErr)

	queries := db.New(mockDB)
	habitRepo := repository.NewHabitRepository(queries)
	routineRepo := repository.NewRoutineRepository(queries)
	userClient := &mockUserServiceClient{}
	service := NewHabitService(habitRepo, routineRepo, userClient)

	habit, err := service.GetHabitByID(context.Background(), habitID)

	if habit != (db.Habit{}) {
		t.Fatalf("GetHabitByID() expected empty habit struct, got %#v", habit)
	}

	if !errors.Is(err, dbErr) {
		t.Fatalf("GetHabitByID() error = %v, want %v", err, dbErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_CreateHabit(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}

	defer mockDB.Close()

	habitID := uuid.New()
	userID := uuid.New()
	name := "habito 1"
	desc := "desc 1"
	img_url := "img.png"
	createdAt := time.Now().UTC().Truncate(time.Microsecond)

	arg := repository.CreateHabitParams{
		UserID:      userID,
		Name:        name,
		Description: sql.NullString{String: desc, Valid: true},
		ImageUrl:    sql.NullString{String: img_url, Valid: true},
	}

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO habits (user_id, name, description, image_url, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id, user_id, name, description, image_url, created_at")).
		WithArgs(arg.UserID, arg.Name, arg.Description, arg.ImageUrl).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).
				AddRow(habitID, userID, name, desc, img_url, createdAt),
		)

	queries := db.New(mockDB)
	habitRepo := repository.NewHabitRepository(queries)
	routineRepo := repository.NewRoutineRepository(queries)
	userClient := &mockUserServiceClient{}
	service := NewHabitService(habitRepo, routineRepo, userClient)

	habit, err := service.CreateHabit(context.Background(), arg)

	if err != nil {
		t.Fatalf("CreateHabit() returned unexpected error: %v", err)
	}
	if habit.ID != habitID {
		t.Fatalf("CreateHabit() returned wrong habit ID: got %s, want %s", habit.ID, habitID)
	}
	if habit.UserID != userID {
		t.Fatalf("CreateHabit() returned wrong user ID: got %s, want %s", habit.UserID, userID)
	}
	if habit.Name != name {
		t.Fatalf("CreateHabit() returned wrong name: got %s, want %s", habit.Name, name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}

}

func TestHabitService_CreateHabit_DatabaseError(t *testing.T) {

}
