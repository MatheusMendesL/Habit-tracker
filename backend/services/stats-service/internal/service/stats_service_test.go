package service

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbUser "shared/pb/user"
	"stats-service/db"
	AppErr "stats-service/internal/errors"
	"stats-service/internal/repository"
)

type mockStatsUserServiceClient struct {
	userByIDErrs []error
	userByIDIdx  int
}

func (m *mockStatsUserServiceClient) nextUserByIDErr() error {
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

func (m *mockStatsUserServiceClient) GetUserByID(ctx context.Context, in *pbUser.GetUserByIDRequest, opts ...grpc.CallOption) (*pbUser.GetUserByIDResponse, error) {
	if err := m.nextUserByIDErr(); err != nil {
		return nil, err
	}
	return &pbUser.GetUserByIDResponse{User: &pbUser.User{Id: in.UserId, Name: "Teste", Email: "teste@email.com"}}, nil
}

func (m *mockStatsUserServiceClient) GetUsersByIDs(ctx context.Context, in *pbUser.GetUsersByIDsRequest, opts ...grpc.CallOption) (*pbUser.GetUsersByIDsResponse, error) {
	return &pbUser.GetUsersByIDsResponse{}, nil
}

func (m *mockStatsUserServiceClient) EditPassword(ctx context.Context, in *pbUser.EditPasswordRequest, opts ...grpc.CallOption) (*pbUser.EditPasswordResponse, error) {
	return &pbUser.EditPasswordResponse{}, nil
}

func (m *mockStatsUserServiceClient) EditUser(ctx context.Context, in *pbUser.EditUserRequest, opts ...grpc.CallOption) (*pbUser.EditUserResponse, error) {
	return &pbUser.EditUserResponse{}, nil
}

func (m *mockStatsUserServiceClient) DeleteUser(ctx context.Context, in *pbUser.DeleteUserRequest, opts ...grpc.CallOption) (*pbUser.DeleteUserResponse, error) {
	return &pbUser.DeleteUserResponse{}, nil
}

func (m *mockStatsUserServiceClient) SearchUser(ctx context.Context, in *pbUser.SearchUserRequest, opts ...grpc.CallOption) (*pbUser.SearchUserResponse, error) {
	return &pbUser.SearchUserResponse{}, nil
}

func TestStatsService_CreateUserStats(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	updatedAt := createdAt.Add(time.Minute)

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO user_stats(user_id) VALUES ($1) RETURNING user_id, completed_habits, completed_routines, current_habit_streak, longest_habit_streak, current_routine_streak, longest_routine_streak, created_at, updated_at")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id", "completed_habits", "completed_routines", "current_habit_streak", "longest_habit_streak", "current_routine_streak", "longest_routine_streak", "created_at", "updated_at"}).
				AddRow(userID, 0, 0, 0, 0, 0, 0, createdAt, updatedAt),
		)

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	stats, err := service.CreateUserStats(context.Background(), userID)
	if err != nil {
		t.Fatalf("CreateUserStats() returned unexpected error: %v", err)
	}
	if stats.UserID != userID {
		t.Fatalf("CreateUserStats() wrong user ID: got %s, want %s", stats.UserID, userID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestStatsService_CreateUserStats_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	stats, err := service.CreateUserStats(context.Background(), uuid.Nil)
	if stats != (db.UserStats{}) {
		t.Fatalf("CreateUserStats() expected empty stats, got %#v", stats)
	}
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("CreateUserStats() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestStatsService_CreateUserStats_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	expectedErr := errors.New("database failure")
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO user_stats(user_id) VALUES ($1) RETURNING user_id, completed_habits, completed_routines, current_habit_streak, longest_habit_streak, current_routine_streak, longest_routine_streak, created_at, updated_at")).
		WithArgs(userID).
		WillReturnError(expectedErr)

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	stats, err := service.CreateUserStats(context.Background(), userID)
	if stats != (db.UserStats{}) {
		t.Fatalf("CreateUserStats() expected empty stats, got %#v", stats)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("CreateUserStats() error = %v, want %v", err, expectedErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestStatsService_GetUserStats(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	updatedAt := createdAt.Add(time.Minute)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, completed_habits, completed_routines, current_habit_streak, longest_habit_streak, current_routine_streak, longest_routine_streak, created_at, updated_at FROM user_stats WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id", "completed_habits", "completed_routines", "current_habit_streak", "longest_habit_streak", "current_routine_streak", "longest_routine_streak", "created_at", "updated_at"}).
				AddRow(userID, 12, 7, 4, 8, 3, 9, createdAt, updatedAt),
		)

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	stats, err := service.GetUserStats(context.Background(), userID)
	if err != nil {
		t.Fatalf("GetUserStats() returned unexpected error: %v", err)
	}
	if stats.UserID != userID {
		t.Fatalf("GetUserStats() wrong user ID: got %s, want %s", stats.UserID, userID)
	}
	if stats.CompletedHabits != 12 || stats.CompletedRoutines != 7 {
		t.Fatalf("GetUserStats() wrong counters: got %d/%d, want 12/7", stats.CompletedHabits, stats.CompletedRoutines)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestStatsService_GetUserStats_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	stats, err := service.GetUserStats(context.Background(), uuid.Nil)
	if stats != (db.UserStats{}) {
		t.Fatalf("GetUserStats() expected empty stats, got %#v", stats)
	}
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("GetUserStats() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestStatsService_GetUserStats_UserNotFound(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{userByIDErrs: []error{status.Error(codes.NotFound, AppErr.ErrUserNotFound.Error())}}, nil)

	stats, err := service.GetUserStats(context.Background(), userID)
	if stats != (db.UserStats{}) {
		t.Fatalf("GetUserStats() expected empty stats, got %#v", stats)
	}
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("GetUserStats() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}
}

func TestStatsService_GetUserStats_StatsNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, completed_habits, completed_routines, current_habit_streak, longest_habit_streak, current_routine_streak, longest_routine_streak, created_at, updated_at FROM user_stats WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	stats, err := service.GetUserStats(context.Background(), userID)
	if stats != (db.UserStats{}) {
		t.Fatalf("GetUserStats() expected empty stats, got %#v", stats)
	}
	if !errors.Is(err, AppErr.ErrUserStatsNotFound) {
		t.Fatalf("GetUserStats() error = %v, want %v", err, AppErr.ErrUserStatsNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestStatsService_GetUserStats_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	expectedErr := errors.New("database failure")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, completed_habits, completed_routines, current_habit_streak, longest_habit_streak, current_routine_streak, longest_routine_streak, created_at, updated_at FROM user_stats WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnError(expectedErr)

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	stats, err := service.GetUserStats(context.Background(), userID)
	if stats != (db.UserStats{}) {
		t.Fatalf("GetUserStats() expected empty stats, got %#v", stats)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("GetUserStats() error = %v, want %v", err, expectedErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestStatsService_DeleteUserStats(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	updatedAt := createdAt.Add(time.Minute)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, completed_habits, completed_routines, current_habit_streak, longest_habit_streak, current_routine_streak, longest_routine_streak, created_at, updated_at FROM user_stats WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id", "completed_habits", "completed_routines", "current_habit_streak", "longest_habit_streak", "current_routine_streak", "longest_routine_streak", "created_at", "updated_at"}).
				AddRow(userID, 1, 1, 2, 4, 3, 5, createdAt, updatedAt),
		)
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM user_stats WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	if err := service.DeleteUserStats(context.Background(), userID); err != nil {
		t.Fatalf("DeleteUserStats() returned unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestStatsService_DeleteUserStats_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	err = service.DeleteUserStats(context.Background(), uuid.Nil)
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("DeleteUserStats() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestStatsService_DeleteUserStats_UserNotFound(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{userByIDErrs: []error{status.Error(codes.NotFound, AppErr.ErrUserNotFound.Error())}}, nil)

	err = service.DeleteUserStats(context.Background(), userID)
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("DeleteUserStats() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}
}

func TestStatsService_DeleteUserStats_StatsNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, completed_habits, completed_routines, current_habit_streak, longest_habit_streak, current_routine_streak, longest_routine_streak, created_at, updated_at FROM user_stats WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	err = service.DeleteUserStats(context.Background(), userID)
	if !errors.Is(err, AppErr.ErrUserStatsNotFound) {
		t.Fatalf("DeleteUserStats() error = %v, want %v", err, AppErr.ErrUserStatsNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestStatsService_DeleteUserStats_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	updatedAt := createdAt.Add(time.Minute)
	expectedErr := errors.New("database failure")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, completed_habits, completed_routines, current_habit_streak, longest_habit_streak, current_routine_streak, longest_routine_streak, created_at, updated_at FROM user_stats WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id", "completed_habits", "completed_routines", "current_habit_streak", "longest_habit_streak", "current_routine_streak", "longest_routine_streak", "created_at", "updated_at"}).
				AddRow(userID, 1, 1, 2, 4, 3, 5, createdAt, updatedAt),
		)
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM user_stats WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnError(expectedErr)

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	err = service.DeleteUserStats(context.Background(), userID)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("DeleteUserStats() error = %v, want %v", err, expectedErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestStatsService_RegisterHabitCompletion(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	createdAt := completedAt.Add(-time.Hour)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, completed_habits, completed_routines, current_habit_streak, longest_habit_streak, current_routine_streak, longest_routine_streak, created_at, updated_at FROM user_stats WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id", "completed_habits", "completed_routines", "current_habit_streak", "longest_habit_streak", "current_routine_streak", "longest_routine_streak", "created_at", "updated_at"}).
				AddRow(userID, 0, 0, 0, 0, 0, 0, createdAt, createdAt),
		)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE user_stats SET completed_habits = completed_habits + 1, updated_at = NOW() WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	if err := service.RegisterHabitCompletion(context.Background(), userID, "habit-1", completedAt); err != nil {
		t.Fatalf("RegisterHabitCompletion() returned unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestStatsService_RegisterHabitCompletion_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	err = service.RegisterHabitCompletion(context.Background(), uuid.Nil, "habit-1", time.Now())
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("RegisterHabitCompletion() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestStatsService_RegisterHabitCompletion_UserNotFound(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{userByIDErrs: []error{status.Error(codes.NotFound, AppErr.ErrUserNotFound.Error())}}, nil)

	err = service.RegisterHabitCompletion(context.Background(), userID, "habit-1", time.Now())
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("RegisterHabitCompletion() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}
}

func TestStatsService_RegisterHabitCompletion_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	expectedErr := errors.New("database failure")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, completed_habits, completed_routines, current_habit_streak, longest_habit_streak, current_routine_streak, longest_routine_streak, created_at, updated_at FROM user_stats WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id", "completed_habits", "completed_routines", "current_habit_streak", "longest_habit_streak", "current_routine_streak", "longest_routine_streak", "created_at", "updated_at"}).
				AddRow(userID, 0, 0, 0, 0, 0, 0, completedAt, completedAt),
		)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE user_stats SET completed_habits = completed_habits + 1, updated_at = NOW() WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnError(expectedErr)

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	err = service.RegisterHabitCompletion(context.Background(), userID, "habit-1", completedAt)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("RegisterHabitCompletion() error = %v, want %v", err, expectedErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestStatsService_UndoHabitCompletion(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	createdAt := completedAt.Add(-time.Hour)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, completed_habits, completed_routines, current_habit_streak, longest_habit_streak, current_routine_streak, longest_routine_streak, created_at, updated_at FROM user_stats WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id", "completed_habits", "completed_routines", "current_habit_streak", "longest_habit_streak", "current_routine_streak", "longest_routine_streak", "created_at", "updated_at"}).
				AddRow(userID, 2, 0, 0, 0, 0, 0, createdAt, createdAt),
		)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE user_stats SET completed_habits = GREATEST(completed_habits - 1, 0), updated_at = NOW() WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	if err := service.UndoHabitCompletion(context.Background(), userID, "habit-1", completedAt); err != nil {
		t.Fatalf("UndoHabitCompletion() returned unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestStatsService_UndoHabitCompletion_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	err = service.UndoHabitCompletion(context.Background(), uuid.Nil, "habit-1", time.Now())
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("UndoHabitCompletion() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestStatsService_UndoHabitCompletion_UserNotFound(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{userByIDErrs: []error{status.Error(codes.NotFound, AppErr.ErrUserNotFound.Error())}}, nil)

	err = service.UndoHabitCompletion(context.Background(), userID, "habit-1", time.Now())
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("UndoHabitCompletion() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}
}

func TestStatsService_UndoHabitCompletion_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	expectedErr := errors.New("database failure")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, completed_habits, completed_routines, current_habit_streak, longest_habit_streak, current_routine_streak, longest_routine_streak, created_at, updated_at FROM user_stats WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id", "completed_habits", "completed_routines", "current_habit_streak", "longest_habit_streak", "current_routine_streak", "longest_routine_streak", "created_at", "updated_at"}).
				AddRow(userID, 2, 0, 0, 0, 0, 0, completedAt, completedAt),
		)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE user_stats SET completed_habits = GREATEST(completed_habits - 1, 0), updated_at = NOW() WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnError(expectedErr)

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	err = service.UndoHabitCompletion(context.Background(), userID, "habit-1", completedAt)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("UndoHabitCompletion() error = %v, want %v", err, expectedErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestStatsService_RegisterRoutineCompletion(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	createdAt := completedAt.Add(-time.Hour)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, completed_habits, completed_routines, current_habit_streak, longest_habit_streak, current_routine_streak, longest_routine_streak, created_at, updated_at FROM user_stats WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id", "completed_habits", "completed_routines", "current_habit_streak", "longest_habit_streak", "current_routine_streak", "longest_routine_streak", "created_at", "updated_at"}).
				AddRow(userID, 0, 0, 0, 0, 0, 0, createdAt, createdAt),
		)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE user_stats SET completed_routines = completed_routines + 1, updated_at = NOW() WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	if err := service.RegisterRoutineCompletion(context.Background(), userID, "routine-1", completedAt); err != nil {
		t.Fatalf("RegisterRoutineCompletion() returned unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestStatsService_RegisterRoutineCompletion_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	err = service.RegisterRoutineCompletion(context.Background(), uuid.Nil, "routine-1", time.Now())
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("RegisterRoutineCompletion() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestStatsService_RegisterRoutineCompletion_UserNotFound(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{userByIDErrs: []error{status.Error(codes.NotFound, AppErr.ErrUserNotFound.Error())}}, nil)

	err = service.RegisterRoutineCompletion(context.Background(), userID, "routine-1", time.Now())
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("RegisterRoutineCompletion() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}
}

func TestStatsService_RegisterRoutineCompletion_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	expectedErr := errors.New("database failure")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, completed_habits, completed_routines, current_habit_streak, longest_habit_streak, current_routine_streak, longest_routine_streak, created_at, updated_at FROM user_stats WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id", "completed_habits", "completed_routines", "current_habit_streak", "longest_habit_streak", "current_routine_streak", "longest_routine_streak", "created_at", "updated_at"}).
				AddRow(userID, 0, 0, 0, 0, 0, 0, completedAt, completedAt),
		)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE user_stats SET completed_routines = completed_routines + 1, updated_at = NOW() WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnError(expectedErr)

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	err = service.RegisterRoutineCompletion(context.Background(), userID, "routine-1", completedAt)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("RegisterRoutineCompletion() error = %v, want %v", err, expectedErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestStatsService_UndoRoutineCompletion(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	createdAt := completedAt.Add(-time.Hour)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, completed_habits, completed_routines, current_habit_streak, longest_habit_streak, current_routine_streak, longest_routine_streak, created_at, updated_at FROM user_stats WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id", "completed_habits", "completed_routines", "current_habit_streak", "longest_habit_streak", "current_routine_streak", "longest_routine_streak", "created_at", "updated_at"}).
				AddRow(userID, 0, 3, 0, 0, 0, 0, createdAt, createdAt),
		)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE user_stats SET completed_routines = GREATEST(completed_routines - 1, 0), updated_at = NOW() WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	if err := service.UndoRoutineCompletion(context.Background(), userID, "routine-1", completedAt); err != nil {
		t.Fatalf("UndoRoutineCompletion() returned unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestStatsService_UndoRoutineCompletion_InvalidArgument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	err = service.UndoRoutineCompletion(context.Background(), uuid.Nil, "routine-1", time.Now())
	if !errors.Is(err, AppErr.ErrInvalidArgument) {
		t.Fatalf("UndoRoutineCompletion() error = %v, want %v", err, AppErr.ErrInvalidArgument)
	}
}

func TestStatsService_UndoRoutineCompletion_UserNotFound(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{userByIDErrs: []error{status.Error(codes.NotFound, AppErr.ErrUserNotFound.Error())}}, nil)

	err = service.UndoRoutineCompletion(context.Background(), userID, "routine-1", time.Now())
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("UndoRoutineCompletion() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}
}

func TestStatsService_UndoRoutineCompletion_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	expectedErr := errors.New("database failure")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, completed_habits, completed_routines, current_habit_streak, longest_habit_streak, current_routine_streak, longest_routine_streak, created_at, updated_at FROM user_stats WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id", "completed_habits", "completed_routines", "current_habit_streak", "longest_habit_streak", "current_routine_streak", "longest_routine_streak", "created_at", "updated_at"}).
				AddRow(userID, 0, 3, 0, 0, 0, 0, completedAt, completedAt),
		)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE user_stats SET completed_routines = GREATEST(completed_routines - 1, 0), updated_at = NOW() WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnError(expectedErr)

	queries := db.New(mockDB)
	repo := repository.NewStatsRepository(queries)
	service := NewStatsService(repo, &mockStatsUserServiceClient{}, nil)

	err = service.UndoRoutineCompletion(context.Background(), userID, "routine-1", completedAt)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("UndoRoutineCompletion() error = %v, want %v", err, expectedErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}
