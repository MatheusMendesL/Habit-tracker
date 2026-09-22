package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"habit-service/db"
	AppErr "habit-service/internal/errors"
	"habit-service/internal/repository"
	pbUser "shared/pb/user"
)

type mockRoutineUserServiceClient struct {
	userByIDErrs []error
	userByIDIdx  int
}

func (m *mockRoutineUserServiceClient) nextUserByIDErr() error {
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

func (m *mockRoutineUserServiceClient) GetUserByID(ctx context.Context, in *pbUser.GetUserByIDRequest, opts ...grpc.CallOption) (*pbUser.GetUserByIDResponse, error) {
	if err := m.nextUserByIDErr(); err != nil {
		return nil, err
	}
	return &pbUser.GetUserByIDResponse{User: &pbUser.User{Id: in.UserId, Name: "Teste", Email: "teste@email.com"}}, nil
}

func (m *mockRoutineUserServiceClient) GetUsersByIDs(ctx context.Context, in *pbUser.GetUsersByIDsRequest, opts ...grpc.CallOption) (*pbUser.GetUsersByIDsResponse, error) {
	return &pbUser.GetUsersByIDsResponse{}, nil
}

func (m *mockRoutineUserServiceClient) EditPassword(ctx context.Context, in *pbUser.EditPasswordRequest, opts ...grpc.CallOption) (*pbUser.EditPasswordResponse, error) {
	return &pbUser.EditPasswordResponse{}, nil
}

func (m *mockRoutineUserServiceClient) EditUser(ctx context.Context, in *pbUser.EditUserRequest, opts ...grpc.CallOption) (*pbUser.EditUserResponse, error) {
	return &pbUser.EditUserResponse{}, nil
}

func (m *mockRoutineUserServiceClient) DeleteUser(ctx context.Context, in *pbUser.DeleteUserRequest, opts ...grpc.CallOption) (*pbUser.DeleteUserResponse, error) {
	return &pbUser.DeleteUserResponse{}, nil
}

func (m *mockRoutineUserServiceClient) SearchUser(ctx context.Context, in *pbUser.SearchUserRequest, opts ...grpc.CallOption) (*pbUser.SearchUserResponse, error) {
	return &pbUser.SearchUserResponse{}, nil
}

func newRoutineTestService(t *testing.T, mockDB *sql.DB, userClient pbUser.UserServiceClient) *RoutineService {
	t.Helper()
	queries := db.New(mockDB)
	habitRepo := repository.NewHabitRepository(queries)
	routineRepo := repository.NewRoutineRepository(queries)
	return NewRoutineService(routineRepo, userClient, NewHabitService(habitRepo, routineRepo, userClient))
}

func TestRoutineService_CreateRoutine(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	arg := repository.CreateRoutineParams{UserID: userID, Name: "rotina 1"}

	mock.ExpectQuery("(?s).*INSERT INTO routines .* RETURNING id, user_id, name, created_at.*").
		WithArgs(userID, arg.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, userID, "rotina 1", createdAt))

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	routine, err := service.CreateRoutine(context.Background(), arg)
	if err != nil {
		t.Fatalf("CreateRoutine() returned unexpected error: %v", err)
	}
	if routine.ID != routineID || routine.UserID != userID || routine.Name != arg.Name {
		t.Fatalf("CreateRoutine() returned unexpected result: %#v", routine)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_CreateRoutine_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	result, err := service.CreateRoutine(context.Background(), repository.CreateRoutineParams{UserID: uuid.Nil, Name: ""})
	if result != (db.Routine{}) {
		t.Fatalf("CreateRoutine() expected empty routine, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("CreateRoutine() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestRoutineService_CreateRoutine_UserNotFound(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userClient := &mockRoutineUserServiceClient{userByIDErrs: []error{status.Error(codes.NotFound, AppErr.ErrUserNotFound.Error())}}
	service := newRoutineTestService(t, mockDB, userClient)
	result, err := service.CreateRoutine(context.Background(), repository.CreateRoutineParams{UserID: uuid.New(), Name: "rotina 1"})
	if result != (db.Routine{}) {
		t.Fatalf("CreateRoutine() expected empty routine, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("CreateRoutine() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}
}

func TestRoutineService_CreateRoutine_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	arg := repository.CreateRoutineParams{UserID: userID, Name: "rotina 1"}
	dbErr := errors.New("database failure")
	mock.ExpectQuery("(?s).*INSERT INTO routines .* RETURNING id, user_id, name, created_at.*").
		WithArgs(userID, arg.Name).
		WillReturnError(dbErr)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	result, err := service.CreateRoutine(context.Background(), arg)
	if result != (db.Routine{}) {
		t.Fatalf("CreateRoutine() expected empty routine, got %#v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("CreateRoutine() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_GetRoutineByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, userID, "rotina 1", createdAt))

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	result, err := service.GetRoutineByID(context.Background(), routineID)
	if err != nil {
		t.Fatalf("GetRoutineByID() returned unexpected error: %v", err)
	}
	if result.ID != routineID || result.UserID != userID || result.Name != "rotina 1" {
		t.Fatalf("GetRoutineByID() returned unexpected result: %#v", result)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_GetRoutineByID_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	result, err := service.GetRoutineByID(context.Background(), uuid.Nil)
	if result != (db.Routine{}) {
		t.Fatalf("GetRoutineByID() expected empty routine, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("GetRoutineByID() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestRoutineService_GetRoutineByID_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnError(sql.ErrNoRows)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	result, err := service.GetRoutineByID(context.Background(), routineID)
	if result != (db.Routine{}) {
		t.Fatalf("GetRoutineByID() expected empty routine, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrRoutineNotFound) {
		t.Fatalf("GetRoutineByID() error = %v, want %v", err, AppErr.ErrRoutineNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_GetRoutineByID_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	dbErr := errors.New("database failure")
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnError(dbErr)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	result, err := service.GetRoutineByID(context.Background(), routineID)
	if result != (db.Routine{}) {
		t.Fatalf("GetRoutineByID() expected empty routine, got %#v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("GetRoutineByID() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_EditRoutine(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	params := db.UpdateRoutineParams{ID: routineID, Name: sql.NullString{String: "nova rotina", Valid: true}}

	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, userID, "rotina 1", createdAt))
	mock.ExpectQuery("(?s).*UPDATE routines.* WHERE id = \\$2 RETURNING id, user_id, name, created_at.*").
		WithArgs(params.Name, params.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, userID, "nova rotina", createdAt))

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	result, err := service.EditRoutine(context.Background(), params)
	if err != nil {
		t.Fatalf("EditRoutine() returned unexpected error: %v", err)
	}
	if result.Name != "nova rotina" {
		t.Fatalf("EditRoutine() returned unexpected result: %#v", result)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_EditRoutine_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	result, err := service.EditRoutine(context.Background(), db.UpdateRoutineParams{ID: uuid.Nil})
	if result != (db.Routine{}) {
		t.Fatalf("EditRoutine() expected empty routine, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("EditRoutine() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestRoutineService_EditRoutine_RoutineNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnError(sql.ErrNoRows)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	result, err := service.EditRoutine(context.Background(), db.UpdateRoutineParams{ID: routineID})
	if result != (db.Routine{}) {
		t.Fatalf("EditRoutine() expected empty routine, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrRoutineNotFound) {
		t.Fatalf("EditRoutine() error = %v, want %v", err, AppErr.ErrRoutineNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_EditRoutine_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	params := db.UpdateRoutineParams{ID: routineID, Name: sql.NullString{String: "nova rotina", Valid: true}}
	dbErr := errors.New("database failure")
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, uuid.New(), "rotina 1", createdAt))
	mock.ExpectQuery("(?s).*UPDATE routines.* WHERE id = \\$2 RETURNING id, user_id, name, created_at.*").
		WithArgs(params.Name, params.ID).
		WillReturnError(dbErr)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	result, err := service.EditRoutine(context.Background(), params)
	if result != (db.Routine{}) {
		t.Fatalf("EditRoutine() expected empty routine, got %#v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("EditRoutine() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_DeleteRoutine(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, uuid.New(), "rotina 1", createdAt))
	mock.ExpectExec("(?s).*DELETE.* FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.DeleteRoutine(context.Background(), routineID); err != nil {
		t.Fatalf("DeleteRoutine() returned unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_DeleteRoutine_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.DeleteRoutine(context.Background(), uuid.Nil); !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("DeleteRoutine() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestRoutineService_DeleteRoutine_RoutineNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnError(sql.ErrNoRows)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.DeleteRoutine(context.Background(), routineID); !errors.Is(err, AppErr.ErrRoutineNotFound) {
		t.Fatalf("DeleteRoutine() error = %v, want %v", err, AppErr.ErrRoutineNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_DeleteRoutine_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	dbErr := errors.New("database failure")
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, uuid.New(), "rotina 1", createdAt))
	mock.ExpectExec("(?s).*DELETE.* FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnError(dbErr)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.DeleteRoutine(context.Background(), routineID); !errors.Is(err, dbErr) {
		t.Fatalf("DeleteRoutine() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_ListRoutinesByUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE user_id = \\$1.*").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(uuid.New(), userID, "rotina 1", createdAt))

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	result, err := service.ListRoutinesByUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("ListRoutinesByUser() returned unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].UserID != userID {
		t.Fatalf("ListRoutinesByUser() returned unexpected result: %#v", result)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_ListRoutinesByUser_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	result, err := service.ListRoutinesByUser(context.Background(), uuid.Nil)
	if len(result) != 0 {
		t.Fatalf("ListRoutinesByUser() expected empty slice, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("ListRoutinesByUser() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestRoutineService_ListRoutinesByUser_UserNotFound(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userClient := &mockRoutineUserServiceClient{userByIDErrs: []error{status.Error(codes.NotFound, AppErr.ErrUserNotFound.Error())}}
	service := newRoutineTestService(t, mockDB, userClient)
	result, err := service.ListRoutinesByUser(context.Background(), uuid.New())
	if len(result) != 0 {
		t.Fatalf("ListRoutinesByUser() expected empty slice, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("ListRoutinesByUser() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}
}

func TestRoutineService_ListRoutinesByUser_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	dbErr := errors.New("database failure")
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE user_id = \\$1.*").
		WithArgs(userID).
		WillReturnError(dbErr)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	result, err := service.ListRoutinesByUser(context.Background(), userID)
	if len(result) != 0 {
		t.Fatalf("ListRoutinesByUser() expected empty slice, got %#v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("ListRoutinesByUser() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_AddHabitToRoutine(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	habitID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, userID, "rotina 1", createdAt))
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = \\$1.*").
		WithArgs(habitID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).AddRow(habitID, userID, "habito 1", "desc", "img.png", createdAt))
	mock.ExpectExec("INSERT INTO routine_habits .* VALUES \\(\\$1, \\$2\\)").
		WithArgs(routineID, habitID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.AddHabitToRoutine(context.Background(), db.AddHabitToRoutineParams{RoutineID: routineID, HabitID: habitID}); err != nil {
		t.Fatalf("AddHabitToRoutine() returned unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_AddHabitToRoutine_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.AddHabitToRoutine(context.Background(), db.AddHabitToRoutineParams{RoutineID: uuid.Nil, HabitID: uuid.New()}); !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("AddHabitToRoutine() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestRoutineService_AddHabitToRoutine_RoutineNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	habitID := uuid.New()
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnError(sql.ErrNoRows)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.AddHabitToRoutine(context.Background(), db.AddHabitToRoutineParams{RoutineID: routineID, HabitID: habitID}); !errors.Is(err, AppErr.ErrRoutineNotFound) {
		t.Fatalf("AddHabitToRoutine() error = %v, want %v", err, AppErr.ErrRoutineNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_AddHabitToRoutine_HabitNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	habitID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, uuid.New(), "rotina 1", createdAt))
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = \\$1.*").
		WithArgs(habitID).
		WillReturnError(sql.ErrNoRows)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.AddHabitToRoutine(context.Background(), db.AddHabitToRoutineParams{RoutineID: routineID, HabitID: habitID}); !errors.Is(err, AppErr.ErrHabitNotFound) {
		t.Fatalf("AddHabitToRoutine() error = %v, want %v", err, AppErr.ErrHabitNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_AddHabitToRoutine_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	habitID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	dbErr := errors.New("database failure")
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, userID, "rotina 1", createdAt))
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = \\$1.*").
		WithArgs(habitID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).AddRow(habitID, userID, "habito 1", "desc", "img.png", createdAt))
	mock.ExpectExec("INSERT INTO routine_habits .* VALUES \\(\\$1, \\$2\\)").
		WithArgs(routineID, habitID).
		WillReturnError(dbErr)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.AddHabitToRoutine(context.Background(), db.AddHabitToRoutineParams{RoutineID: routineID, HabitID: habitID}); !errors.Is(err, dbErr) {
		t.Fatalf("AddHabitToRoutine() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_RemoveHabitFromRoutine(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	habitID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, userID, "rotina 1", createdAt))
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = \\$1.*").
		WithArgs(habitID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).AddRow(habitID, userID, "habito 1", "desc", "img.png", createdAt))
	mock.ExpectExec("(?s).*DELETE.* FROM routine_habits WHERE routine_id = \\$1 AND habit_id = \\$2.*").
		WithArgs(routineID, habitID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.RemoveHabitFromRoutine(context.Background(), db.RemoveHabitFromRoutineParams{RoutineID: routineID, HabitID: habitID}); err != nil {
		t.Fatalf("RemoveHabitFromRoutine() returned unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_RemoveHabitFromRoutine_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.RemoveHabitFromRoutine(context.Background(), db.RemoveHabitFromRoutineParams{RoutineID: uuid.Nil, HabitID: uuid.New()}); !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("RemoveHabitFromRoutine() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestRoutineService_RemoveHabitFromRoutine_RoutineNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	habitID := uuid.New()
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnError(sql.ErrNoRows)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.RemoveHabitFromRoutine(context.Background(), db.RemoveHabitFromRoutineParams{RoutineID: routineID, HabitID: habitID}); !errors.Is(err, AppErr.ErrRoutineNotFound) {
		t.Fatalf("RemoveHabitFromRoutine() error = %v, want %v", err, AppErr.ErrRoutineNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_RemoveHabitFromRoutine_HabitNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	habitID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, uuid.New(), "rotina 1", createdAt))
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = \\$1.*").
		WithArgs(habitID).
		WillReturnError(sql.ErrNoRows)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.RemoveHabitFromRoutine(context.Background(), db.RemoveHabitFromRoutineParams{RoutineID: routineID, HabitID: habitID}); !errors.Is(err, AppErr.ErrHabitNotFound) {
		t.Fatalf("RemoveHabitFromRoutine() error = %v, want %v", err, AppErr.ErrHabitNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_RemoveHabitFromRoutine_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	habitID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	dbErr := errors.New("database failure")
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, userID, "rotina 1", createdAt))
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, description, image_url, created_at FROM habits WHERE id = \\$1.*").
		WithArgs(habitID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "image_url", "created_at"}).AddRow(habitID, userID, "habito 1", "desc", "img.png", createdAt))
	mock.ExpectExec("(?s).*DELETE.* FROM routine_habits WHERE routine_id = \\$1 AND habit_id = \\$2.*").
		WithArgs(routineID, habitID).
		WillReturnError(dbErr)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.RemoveHabitFromRoutine(context.Background(), db.RemoveHabitFromRoutineParams{RoutineID: routineID, HabitID: habitID}); !errors.Is(err, dbErr) {
		t.Fatalf("RemoveHabitFromRoutine() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_MarkRoutineCompleted(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, uuid.New(), "rotina 1", completedAt))
	mock.ExpectExec("INSERT INTO routine_logs .* VALUES \\(\\$1, \\$2\\).* ON CONFLICT.*").
		WithArgs(routineID, completedAt).
		WillReturnResult(sqlmock.NewResult(0, 1))

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.MarkRoutineCompleted(context.Background(), db.MarkRoutineCompletedParams{RoutineID: routineID, CompletedAt: completedAt}); err != nil {
		t.Fatalf("MarkRoutineCompleted() returned unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_MarkRoutineCompleted_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.MarkRoutineCompleted(context.Background(), db.MarkRoutineCompletedParams{RoutineID: uuid.Nil, CompletedAt: time.Now()}); !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("MarkRoutineCompleted() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestRoutineService_MarkRoutineCompleted_RoutineNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnError(sql.ErrNoRows)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.MarkRoutineCompleted(context.Background(), db.MarkRoutineCompletedParams{RoutineID: routineID, CompletedAt: completedAt}); !errors.Is(err, AppErr.ErrRoutineNotFound) {
		t.Fatalf("MarkRoutineCompleted() error = %v, want %v", err, AppErr.ErrRoutineNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_MarkRoutineCompleted_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	dbErr := errors.New("database failure")
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, uuid.New(), "rotina 1", completedAt))
	mock.ExpectExec("INSERT INTO routine_logs .* VALUES \\(\\$1, \\$2\\).* ON CONFLICT.*").
		WithArgs(routineID, completedAt).
		WillReturnError(dbErr)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.MarkRoutineCompleted(context.Background(), db.MarkRoutineCompletedParams{RoutineID: routineID, CompletedAt: completedAt}); !errors.Is(err, dbErr) {
		t.Fatalf("MarkRoutineCompleted() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_UnmarkRoutineCompleted(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, uuid.New(), "rotina 1", completedAt))
	mock.ExpectExec("(?s).*DELETE.* FROM routine_logs WHERE routine_id = \\$1.*").
		WithArgs(routineID, completedAt).
		WillReturnResult(sqlmock.NewResult(0, 1))

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.UnmarkRoutineCompleted(context.Background(), db.UnmarkRoutineCompletedParams{RoutineID: routineID, CompletedAt: completedAt}); err != nil {
		t.Fatalf("UnmarkRoutineCompleted() returned unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_UnmarkRoutineCompleted_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.UnmarkRoutineCompleted(context.Background(), db.UnmarkRoutineCompletedParams{RoutineID: uuid.Nil, CompletedAt: time.Now()}); !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("UnmarkRoutineCompleted() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestRoutineService_UnmarkRoutineCompleted_RoutineNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnError(sql.ErrNoRows)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.UnmarkRoutineCompleted(context.Background(), db.UnmarkRoutineCompletedParams{RoutineID: routineID, CompletedAt: completedAt}); !errors.Is(err, AppErr.ErrRoutineNotFound) {
		t.Fatalf("UnmarkRoutineCompleted() error = %v, want %v", err, AppErr.ErrRoutineNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_UnmarkRoutineCompleted_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	dbErr := errors.New("database failure")
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, uuid.New(), "rotina 1", completedAt))
	mock.ExpectExec("(?s).*DELETE.* FROM routine_logs WHERE routine_id = \\$1.*").
		WithArgs(routineID, completedAt).
		WillReturnError(dbErr)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	if err := service.UnmarkRoutineCompleted(context.Background(), db.UnmarkRoutineCompletedParams{RoutineID: routineID, CompletedAt: completedAt}); !errors.Is(err, dbErr) {
		t.Fatalf("UnmarkRoutineCompleted() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_GetRoutineLogs(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	startDate := time.Now().UTC().Truncate(time.Microsecond)
	endDate := startDate.Add(24 * time.Hour)
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, uuid.New(), "rotina 1", startDate))
	mock.ExpectQuery("(?s).*SELECT routine_id,.* FROM routine_logs WHERE routine_id = \\$1.* ORDER BY completed_at DESC.*").
		WithArgs(routineID, startDate, endDate).
		WillReturnRows(sqlmock.NewRows([]string{"routine_id", "completed_at"}).AddRow(routineID, startDate.Add(time.Hour)))

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	result, err := service.GetRoutineLogs(context.Background(), db.GetRoutineLogsParams{RoutineID: routineID, StartDate: startDate, EndDate: endDate})
	if err != nil {
		t.Fatalf("GetRoutineLogs() returned unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].RoutineID != routineID {
		t.Fatalf("GetRoutineLogs() unexpected result: %#v", result)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_GetRoutineLogs_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	result, err := service.GetRoutineLogs(context.Background(), db.GetRoutineLogsParams{RoutineID: uuid.Nil, StartDate: time.Now(), EndDate: time.Now().Add(time.Hour)})
	if len(result) != 0 {
		t.Fatalf("GetRoutineLogs() expected empty slice, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("GetRoutineLogs() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestRoutineService_GetRoutineLogs_RoutineNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	startDate := time.Now().UTC().Truncate(time.Microsecond)
	endDate := startDate.Add(24 * time.Hour)
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnError(sql.ErrNoRows)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	result, err := service.GetRoutineLogs(context.Background(), db.GetRoutineLogsParams{RoutineID: routineID, StartDate: startDate, EndDate: endDate})
	if len(result) != 0 {
		t.Fatalf("GetRoutineLogs() expected empty slice, got %#v", result)
	}
	if !errors.Is(err, AppErr.ErrRoutineNotFound) {
		t.Fatalf("GetRoutineLogs() error = %v, want %v", err, AppErr.ErrRoutineNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestRoutineService_GetRoutineLogs_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	routineID := uuid.New()
	startDate := time.Now().UTC().Truncate(time.Microsecond)
	endDate := startDate.Add(24 * time.Hour)
	dbErr := errors.New("database failure")
	mock.ExpectQuery("(?s).*SELECT id, user_id, name, created_at FROM routines WHERE id = \\$1.*").
		WithArgs(routineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).AddRow(routineID, uuid.New(), "rotina 1", startDate))
	mock.ExpectQuery("(?s).*SELECT routine_id,.* FROM routine_logs WHERE routine_id = \\$1.* ORDER BY completed_at DESC.*").
		WithArgs(routineID, startDate, endDate).
		WillReturnError(dbErr)

	service := newRoutineTestService(t, mockDB, &mockRoutineUserServiceClient{})
	result, err := service.GetRoutineLogs(context.Background(), db.GetRoutineLogsParams{RoutineID: routineID, StartDate: startDate, EndDate: endDate})
	if len(result) != 0 {
		t.Fatalf("GetRoutineLogs() expected empty slice, got %#v", result)
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("GetRoutineLogs() error = %v, want %v", err, dbErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}
