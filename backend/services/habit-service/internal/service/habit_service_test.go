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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockUserServiceClient struct {
	userByIDErrs []error
	userByIDIdx  int
}

func (m *mockUserServiceClient) nextUserByIDErr() error {
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

func (m *mockUserServiceClient) GetUserByID(ctx context.Context, in *pbUser.GetUserByIDRequest, opts ...grpc.CallOption) (*pbUser.GetUserByIDResponse, error) {
	if err := m.nextUserByIDErr(); err != nil {
		return nil, err
	}
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

func newHabitTestService(t *testing.T, mockDB *sql.DB, userClient pbUser.UserServiceClient) *HabitService {
	t.Helper()
	queries := db.New(mockDB)
	return NewHabitService(repository.NewHabitRepository(queries), repository.NewRoutineRepository(queries), userClient)
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

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})

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

func TestHabitService_GetHabitByID_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.GetHabitByID(context.Background(), uuid.Nil)
	if result != (db.Habit{}) {
		t.Fatalf("GetHabitByID() expected empty habit, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("GetHabitByID() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestHabitService_GetHabitByID_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnError(sql.ErrNoRows)

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.GetHabitByID(context.Background(), habitID)
	if result != (db.Habit{}) {
		t.Fatalf("GetHabitByID() expected empty habit, got %#v", result)
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

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.GetHabitByID(context.Background(), habitID)
	if result != (db.Habit{}) {
		t.Fatalf("GetHabitByID() expected empty habit, got %#v", result)
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
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	arg := repository.CreateHabitParams{
		UserID:      userID,
		Name:        "habito 1",
		Description: sql.NullString{String: "desc 1", Valid: true},
		ImageUrl:    sql.NullString{String: "img.png", Valid: true},
	}

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO habits (user_id, name, description, image_url, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id, user_id, name, description, image_url, created_at")).
		WithArgs(arg.UserID, arg.Name, arg.Description, arg.ImageUrl).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).
				AddRow(habitID, userID, "habito 1", "desc 1", "img.png", createdAt),
		)

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.CreateHabit(context.Background(), arg)
	if err != nil {
		t.Fatalf("CreateHabit() returned unexpected error: %v", err)
	}
	if result.ID != habitID || result.UserID != userID || result.Name != "habito 1" {
		t.Fatalf("CreateHabit() returned unexpected result: %#v", result)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_CreateHabit_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.CreateHabit(context.Background(), repository.CreateHabitParams{UserID: uuid.Nil, Name: ""})
	if result != (db.Habit{}) {
		t.Fatalf("CreateHabit() expected empty habit, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("CreateHabit() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestHabitService_CreateHabit_UserNotFound(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	userClient := &mockUserServiceClient{userByIDErrs: []error{status.Error(codes.NotFound, AppErr.ErrUserNotFound.Error())}}
	service := newHabitTestService(t, mockDB, userClient)
	result, err := service.CreateHabit(context.Background(), repository.CreateHabitParams{UserID: uuid.New(), Name: "habito 1"})
	if result != (db.Habit{}) {
		t.Fatalf("CreateHabit() expected empty habit, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("CreateHabit() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}
}

func TestHabitService_CreateHabit_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	arg := repository.CreateHabitParams{UserID: userID, Name: "habito 1"}
	dbErr := errors.New("fatal database connection failure")
	mock.ExpectQuery("INSERT INTO habits .* RETURNING id, user_id, name, description, image_url, created_at").
		WithArgs(arg.UserID, arg.Name, arg.Description, arg.ImageUrl).
		WillReturnError(dbErr)

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.CreateHabit(context.Background(), arg)
	if result != (db.Habit{}) {
		t.Fatalf("CreateHabit() expected empty habit, got %#v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("CreateHabit() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_EditHabit(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	params := db.UpdateHabitParams{
		ID:          habitID,
		Name:        sql.NullString{String: "novo nome", Valid: true},
		Description: sql.NullString{String: "nova desc", Valid: true},
		ImageUrl:    sql.NullString{String: "novo.png", Valid: true},
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).AddRow(habitID, userID, "habito 1", "desc", "img.png", createdAt))
	mock.ExpectQuery("UPDATE habits.* WHERE id = \\$4 RETURNING id, user_id, name, description, image_url, created_at").
		WithArgs(params.Name, params.Description, params.ImageUrl, params.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).AddRow(habitID, userID, "novo nome", "nova desc", "novo.png", createdAt))

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.EditHabit(context.Background(), params)
	if err != nil {
		t.Fatalf("EditHabit() returned unexpected error: %v", err)
	}
	if result.Name != "novo nome" || result.Description.String != "nova desc" || result.ImageUrl.String != "novo.png" {
		t.Fatalf("EditHabit() returned unexpected result: %#v", result)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_EditHabit_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.EditHabit(context.Background(), db.UpdateHabitParams{ID: uuid.Nil})
	if result != (db.Habit{}) {
		t.Fatalf("EditHabit() expected empty habit, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("EditHabit() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestHabitService_EditHabit_HabitNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnError(sql.ErrNoRows)

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.EditHabit(context.Background(), db.UpdateHabitParams{ID: habitID})
	if result != (db.Habit{}) {
		t.Fatalf("EditHabit() expected empty habit, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrHabitNotFound) {
		t.Fatalf("EditHabit() error = %v, want %v", err, AppErr.ErrHabitNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_EditHabit_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	params := db.UpdateHabitParams{ID: habitID, Name: sql.NullString{String: "novo nome", Valid: true}}
	dbErr := errors.New("fatal database connection failure")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).AddRow(habitID, uuid.New(), "habito 1", "desc", "img.png", createdAt))
	mock.ExpectQuery("UPDATE habits.* WHERE id = \\$4 RETURNING id, user_id, name, description, image_url, created_at").
		WithArgs(params.Name, sql.NullString{}, sql.NullString{}, params.ID).
		WillReturnError(dbErr)

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.EditHabit(context.Background(), params)
	if result != (db.Habit{}) {
		t.Fatalf("EditHabit() expected empty habit, got %#v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("EditHabit() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_DeleteHabit(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).AddRow(habitID, uuid.New(), "habito 1", "desc", "img.png", createdAt))
	mock.ExpectExec("DELETE.* FROM habits WHERE id = \\$1").
		WithArgs(habitID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	if err := service.DeleteHabit(context.Background(), habitID); err != nil {
		t.Fatalf("DeleteHabit() returned unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_DeleteHabit_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	if err := service.DeleteHabit(context.Background(), uuid.Nil); !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("DeleteHabit() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestHabitService_DeleteHabit_HabitNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnError(sql.ErrNoRows)

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	if err := service.DeleteHabit(context.Background(), habitID); !errors.Is(err, AppErr.ErrHabitNotFound) {
		t.Fatalf("DeleteHabit() error = %v, want %v", err, AppErr.ErrHabitNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_DeleteHabit_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	dbErr := errors.New("fatal database connection failure")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).AddRow(habitID, uuid.New(), "habito 1", "desc", "img.png", createdAt))
	mock.ExpectExec("DELETE.* FROM habits WHERE id = \\$1").
		WithArgs(habitID).
		WillReturnError(dbErr)

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	if err := service.DeleteHabit(context.Background(), habitID); !errors.Is(err, dbErr) {
		t.Fatalf("DeleteHabit() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_ListHabitsByUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).AddRow(uuid.New(), userID, "habito 1", "desc", "img.png", createdAt))

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.ListHabitsByUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("ListHabitsByUser() returned unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].UserID != userID {
		t.Fatalf("ListHabitsByUser() returned unexpected result: %#v", result)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_ListHabitsByUser_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.ListHabitsByUser(context.Background(), uuid.Nil)
	if len(result) != 0 {
		t.Fatalf("ListHabitsByUser() expected empty slice, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("ListHabitsByUser() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestHabitService_ListHabitsByUser_UserNotFound(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	userClient := &mockUserServiceClient{userByIDErrs: []error{status.Error(codes.NotFound, AppErr.ErrUserNotFound.Error())}}
	service := newHabitTestService(t, mockDB, userClient)
	result, err := service.ListHabitsByUser(context.Background(), uuid.New())
	if len(result) != 0 {
		t.Fatalf("ListHabitsByUser() expected empty slice, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("ListHabitsByUser() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}
}

func TestHabitService_ListHabitsByUser_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	dbErr := errors.New("fatal database connection failure")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnError(dbErr)

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.ListHabitsByUser(context.Background(), userID)
	if len(result) != 0 {
		t.Fatalf("ListHabitsByUser() expected empty slice, got %#v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("ListHabitsByUser() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_ListHabitsByRoutine(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, created_at FROM routines WHERE id = $1")).
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, uuid.New(), "rotina 1", createdAt))
	mock.ExpectQuery("SELECT h.id,.* FROM habits h.*JOIN routine_habits rh ON rh.habit_id = h.id WHERE rh.routine_id = \\$1").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).AddRow(uuid.New(), uuid.New(), "habito 1", "desc", "img.png", createdAt))

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.ListHabitsByRoutine(context.Background(), routineID)
	if err != nil {
		t.Fatalf("ListHabitsByRoutine() returned unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("ListHabitsByRoutine() returned unexpected result: %#v", result)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_ListHabitsByRoutine_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.ListHabitsByRoutine(context.Background(), uuid.Nil)
	if len(result) != 0 {
		t.Fatalf("ListHabitsByRoutine() expected empty slice, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("ListHabitsByRoutine() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestHabitService_ListHabitsByRoutine_RoutineNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, created_at FROM routines WHERE id = $1")).
		WithArgs(routineID).
		WillReturnError(sql.ErrNoRows)

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.ListHabitsByRoutine(context.Background(), routineID)
	if len(result) != 0 {
		t.Fatalf("ListHabitsByRoutine() expected empty slice, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrRoutineNotFound) {
		t.Fatalf("ListHabitsByRoutine() error = %v, want %v", err, AppErr.ErrRoutineNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_ListHabitsByRoutine_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	dbErr := errors.New("fatal database connection failure")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, created_at FROM routines WHERE id = $1")).
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, uuid.New(), "rotina 1", createdAt))
	mock.ExpectQuery("SELECT h.id,.* FROM habits h.*JOIN routine_habits rh ON rh.habit_id = h.id WHERE rh.routine_id = \\$1").
		WithArgs(routineID).
		WillReturnError(dbErr)

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.ListHabitsByRoutine(context.Background(), routineID)
	if len(result) != 0 {
		t.Fatalf("ListHabitsByRoutine() expected empty slice, got %#v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("ListHabitsByRoutine() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_MarkHabitCompleted(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).AddRow(habitID, uuid.New(), "habito 1", "desc", "img.png", completedAt))
	mock.ExpectExec("INSERT INTO habit_logs .* VALUES \\(\\$1, \\$2\\).* ON CONFLICT.*").
		WithArgs(habitID, completedAt).
		WillReturnResult(sqlmock.NewResult(0, 1))

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	if err := service.MarkHabitCompleted(context.Background(), db.MarkHabitCompletedParams{HabitID: habitID, CompletedAt: completedAt}); err != nil {
		t.Fatalf("MarkHabitCompleted() returned unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_MarkHabitCompleted_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	if err := service.MarkHabitCompleted(context.Background(), db.MarkHabitCompletedParams{HabitID: uuid.Nil, CompletedAt: time.Now()}); !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("MarkHabitCompleted() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestHabitService_MarkHabitCompleted_HabitNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnError(sql.ErrNoRows)

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	if err := service.MarkHabitCompleted(context.Background(), db.MarkHabitCompletedParams{HabitID: habitID, CompletedAt: completedAt}); !errors.Is(err, AppErr.ErrHabitNotFound) {
		t.Fatalf("MarkHabitCompleted() error = %v, want %v", err, AppErr.ErrHabitNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_MarkHabitCompleted_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	dbErr := errors.New("fatal database connection failure")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).AddRow(habitID, uuid.New(), "habito 1", "desc", "img.png", completedAt))
	mock.ExpectExec("INSERT INTO habit_logs .* VALUES \\(\\$1, \\$2\\).* ON CONFLICT.*").
		WithArgs(habitID, completedAt).
		WillReturnError(dbErr)

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	if err := service.MarkHabitCompleted(context.Background(), db.MarkHabitCompletedParams{HabitID: habitID, CompletedAt: completedAt}); !errors.Is(err, dbErr) {
		t.Fatalf("MarkHabitCompleted() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_UnmarkHabitCompleted(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).AddRow(habitID, uuid.New(), "habito 1", "desc", "img.png", completedAt))
	mock.ExpectExec("DELETE.* FROM habit_logs WHERE habit_id = \\$1.*").
		WithArgs(habitID, completedAt).
		WillReturnResult(sqlmock.NewResult(0, 1))

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	if err := service.UnmarkHabitCompleted(context.Background(), db.UnmarkHabitCompletedParams{HabitID: habitID, CompletedAt: completedAt}); err != nil {
		t.Fatalf("UnmarkHabitCompleted() returned unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_UnmarkHabitCompleted_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	if err := service.UnmarkHabitCompleted(context.Background(), db.UnmarkHabitCompletedParams{HabitID: uuid.Nil, CompletedAt: time.Now()}); !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("UnmarkHabitCompleted() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestHabitService_UnmarkHabitCompleted_HabitNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnError(sql.ErrNoRows)

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	if err := service.UnmarkHabitCompleted(context.Background(), db.UnmarkHabitCompletedParams{HabitID: habitID, CompletedAt: completedAt}); !errors.Is(err, AppErr.ErrHabitNotFound) {
		t.Fatalf("UnmarkHabitCompleted() error = %v, want %v", err, AppErr.ErrHabitNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_UnmarkHabitCompleted_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	dbErr := errors.New("fatal database connection failure")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).AddRow(habitID, uuid.New(), "habito 1", "desc", "img.png", completedAt))
	mock.ExpectExec("DELETE.* FROM habit_logs WHERE habit_id = \\$1.*").
		WithArgs(habitID, completedAt).
		WillReturnError(dbErr)

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	if err := service.UnmarkHabitCompleted(context.Background(), db.UnmarkHabitCompletedParams{HabitID: habitID, CompletedAt: completedAt}); !errors.Is(err, dbErr) {
		t.Fatalf("UnmarkHabitCompleted() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_GetHabitLogs(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	startDate := time.Now().UTC().Truncate(time.Microsecond)
	endDate := startDate.Add(24 * time.Hour)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).AddRow(habitID, uuid.New(), "habito 1", "desc", "img.png", startDate))
	mock.ExpectQuery("SELECT habit_id,.* FROM habit_logs WHERE habit_id = \\$1.* ORDER BY completed_at DESC").
		WithArgs(habitID, startDate, endDate).
		WillReturnRows(sqlmock.NewRows([]string{"habit_id", "completed_at"}).AddRow(habitID, startDate.Add(time.Hour)))

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.GetHabitLogs(context.Background(), db.GetHabitLogsParams{HabitID: habitID, StartDate: startDate, EndDate: endDate})
	if err != nil {
		t.Fatalf("GetHabitLogs() returned unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].HabitID != habitID {
		t.Fatalf("GetHabitLogs() returned unexpected result: %#v", result)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_GetHabitLogs_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.GetHabitLogs(context.Background(), db.GetHabitLogsParams{HabitID: uuid.Nil, StartDate: time.Now(), EndDate: time.Now().Add(time.Hour)})
	if len(result) != 0 {
		t.Fatalf("GetHabitLogs() expected empty slice, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("GetHabitLogs() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestHabitService_GetHabitLogs_HabitNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	startDate := time.Now().UTC().Truncate(time.Microsecond)
	endDate := startDate.Add(24 * time.Hour)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnError(sql.ErrNoRows)

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.GetHabitLogs(context.Background(), db.GetHabitLogsParams{HabitID: habitID, StartDate: startDate, EndDate: endDate})
	if len(result) != 0 {
		t.Fatalf("GetHabitLogs() expected empty slice, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrHabitNotFound) {
		t.Fatalf("GetHabitLogs() error = %v, want %v", err, AppErr.ErrHabitNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestHabitService_GetHabitLogs_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create the sqlmock: %v", err)
	}
	defer mockDB.Close()

	habitID := uuid.New()
	startDate := time.Now().UTC().Truncate(time.Microsecond)
	endDate := startDate.Add(24 * time.Hour)
	dbErr := errors.New("fatal database connection failure")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = $1")).
		WithArgs(habitID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).AddRow(habitID, uuid.New(), "habito 1", "desc", "img.png", startDate))
	mock.ExpectQuery("SELECT habit_id,.* FROM habit_logs WHERE habit_id = \\$1.* ORDER BY completed_at DESC").
		WithArgs(habitID, startDate, endDate).
		WillReturnError(dbErr)

	service := newHabitTestService(t, mockDB, &mockUserServiceClient{})
	result, err := service.GetHabitLogs(context.Background(), db.GetHabitLogsParams{HabitID: habitID, StartDate: startDate, EndDate: endDate})
	if len(result) != 0 {
		t.Fatalf("GetHabitLogs() expected empty slice, got %#v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("GetHabitLogs() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}
