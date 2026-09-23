package service

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"user-service/db"
	AppErr "user-service/internal/errors"
	"user-service/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestUserService_GetUserByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users WHERE id = $1")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(userID, "Matheus", "matheus@email.com"),
		)

	queries := db.New(mockDB)
	userRepo := repository.NewUserRepository(queries)
	service := NewUserService(userRepo)

	user, err := service.GetUserByID(context.Background(), userID)
	if err != nil {
		t.Fatalf("GetUserByID() returned unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("GetUserByID() returned nil user")
	}
	if user.ID != userID {
		t.Fatalf("GetUserByID() returned wrong user ID: got %s, want %s", user.ID, userID)
	}
	if user.Name != "Matheus" {
		t.Fatalf("GetUserByID() returned wrong user name: got %s, want %s", user.Name, "Matheus")
	}
	if user.Email != "matheus@email.com" {
		t.Fatalf("GetUserByID() returned wrong user email: got %s, want %s", user.Email, "matheus@email.com")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users WHERE id = $1")).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	queries := db.New(mockDB)
	userRepo := repository.NewUserRepository(queries)
	service := NewUserService(userRepo)

	user, err := service.GetUserByID(context.Background(), userID)
	if user != nil {
		t.Fatalf("GetUserByID() expected nil user, got %#v", user)
	}
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("GetUserByID() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestUserService_GetUserByID_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	dbErr := errors.New("fatal database connection failure")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users WHERE id = $1")).
		WithArgs(userID).
		WillReturnError(dbErr)

	queries := db.New(mockDB)
	userRepo := repository.NewUserRepository(queries)
	service := NewUserService(userRepo)

	users, err := service.GetUserByID(context.Background(), userID)
	if users != nil {
		t.Fatalf("GetUserByID() expected nil user, got %#v", users)
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("GetUserByID() error = %v, want %v", err, dbErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestUserService_SearchUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userIDOne := uuid.New()
	userIDTwo := uuid.New()
	name := "Matheus"
	email := ""

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users WHERE")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(userIDOne, "Matheus", "matheus@email.com").
				AddRow(userIDTwo, "Matheus Mendes", "matheus.mendes@email.com"),
		)

	queries := db.New(mockDB)
	userRepo := repository.NewUserRepository(queries)
	service := NewUserService(userRepo)

	users, err := service.SearchUser(context.Background(), name, email)
	if err != nil {
		t.Fatalf("SearchUser() returned unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("SearchUser() returned %d users, want 2", len(users))
	}
	if users[0].Name != "Matheus" {
		t.Fatalf("SearchUser() first user name = %s, want %s", users[0].Name, "Matheus")
	}
	if users[1].Email != "matheus.mendes@email.com" {
		t.Fatalf("SearchUser() second user email = %s, want %s", users[1].Email, "matheus.mendes@email.com")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestUserService_SearchUser_EmptyResult(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	name := "usuario_inexistente"
	email := ""

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users WHERE")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email"}),
		)

	queries := db.New(mockDB)
	userRepo := repository.NewUserRepository(queries)
	service := NewUserService(userRepo)

	users, err := service.SearchUser(context.Background(), name, email)
	if err != nil {
		t.Fatalf("SearchUser() returned unexpected error: %v", err)
	}

	if len(users) != 0 {
		t.Fatalf("SearchUser() returned %d users, want 0", len(users))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestUserService_SearchUser_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	name := "Matheus"
	email := ""
	dbErr := errors.New("fatal database connection failure")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users WHERE")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(dbErr)

	queries := db.New(mockDB)
	userRepo := repository.NewUserRepository(queries)
	service := NewUserService(userRepo)

	users, err := service.SearchUser(context.Background(), name, email)

	if users != nil {
		t.Fatalf("SearchUser() expected nil user, got %#v", users)
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("SearchUser() error = %v, want %v", err, dbErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestUserService_EditUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	params := db.UpdateUserParams{
		Name:  "Matheus Updated",
		Email: "matheus.updated@email.com",
		ID:    userID,
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET")).
		WithArgs(params.Name, params.Email, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users WHERE id = $1")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(userID, params.Name, params.Email),
		)

	queries := db.New(mockDB)
	userRepo := repository.NewUserRepository(queries)
	service := NewUserService(userRepo)

	user, err := service.EditUser(context.Background(), params)
	if err != nil {
		t.Fatalf("EditUser() returned unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("EditUser() returned nil user")
	}
	if user.ID != userID {
		t.Fatalf("EditUser() returned wrong user ID: got %s, want %s", user.ID, userID)
	}
	if user.Name != params.Name {
		t.Fatalf("EditUser() returned wrong user name: got %s, want %s", user.Name, params.Name)
	}
	if user.Email != params.Email {
		t.Fatalf("EditUser() returned wrong user email: got %s, want %s", user.Email, params.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestUserService_EditUser_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	params := db.UpdateUserParams{
		Name:  "Matheus up",
		Email: "matheus.up@email.com",
		ID:    userID,
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET")).
		WithArgs(params.Name, params.Email, userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users WHERE id = $1")).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	queries := db.New(mockDB)
	userRepo := repository.NewUserRepository(queries)
	service := NewUserService(userRepo)

	user, err := service.EditUser(context.Background(), params)

	if user != nil {
		t.Fatalf("EditUser() expected nil user, got %#v", user)
	}
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("EditUser() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestUserService_EditUser_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	defer mockDB.Close()

	userID := uuid.New()
	params := db.UpdateUserParams{
		Name:  "Matheus up",
		Email: "matheus.up@email.com",
		ID:    userID,
	}
	dbErr := errors.New("fatal database connection failure")

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET")).
		WithArgs(params.Name, params.Email, userID).
		WillReturnError(dbErr)

	queries := db.New(mockDB)
	userRepo := repository.NewUserRepository(queries)
	service := NewUserService(userRepo)

	user, err := service.EditUser(context.Background(), params)
	if user != nil {
		t.Fatalf("EditUser() expected nil user, got %#v", user)
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("EditUser() error = %v, want %v", err, dbErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE id = $1")).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	queries := db.New(mockDB)
	userRepo := repository.NewUserRepository(queries)
	service := NewUserService(userRepo)

	err = service.DeleteUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("DeleteUser() returned unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestUserService_DeleteUser_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE id = $1")).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	queries := db.New(mockDB)
	userRepo := repository.NewUserRepository(queries)
	service := NewUserService(userRepo)

	err = service.DeleteUser(context.Background(), userID)

	if err != nil {
		t.Fatalf("DeleteUser() returned unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestUserService_DeleteUser_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	dbErr := errors.New("fatal database connection failure")
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE id = $1")).
		WithArgs(userID).
		WillReturnError(dbErr)

	queries := db.New(mockDB)
	userRepo := repository.NewUserRepository(queries)
	service := NewUserService(userRepo)

	err = service.DeleteUser(context.Background(), userID)

	if !errors.Is(err, dbErr) {
		t.Fatalf("DeleteUser() error = %v, want %v", err, dbErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestUserService_EditPassword(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	passParams := &db.UpdatePasswordParams{
		ID:       userID,
		Password: "new_secure_password_hash",
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET password")).
		WithArgs(passParams.Password, passParams.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	queries := db.New(mockDB)
	userRepo := repository.NewUserRepository(queries)
	service := NewUserService(userRepo)

	err = service.EditPassword(context.Background(), passParams)
	if err != nil {
		t.Fatalf("EditPassword() returned unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestUserService_EditPassword_RepositoryError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	userID := uuid.New()
	passParams := &db.UpdatePasswordParams{
		ID:       userID,
		Password: "new_secure_password_hash",
	}

	expectedErr := errors.New("database internal failure")

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET password")).
		WithArgs(passParams.Password, passParams.ID).
		WillReturnError(expectedErr)

	queries := db.New(mockDB)
	userRepo := repository.NewUserRepository(queries)
	service := NewUserService(userRepo)

	err = service.EditPassword(context.Background(), passParams)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("EditPassword() error = %v, want %v", err, expectedErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestUserService_GetUsersByIDs(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	idOne := uuid.New()
	idTwo := uuid.New()
	ids := []uuid.UUID{idOne, idTwo}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users WHERE id")).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(idOne, "Matheus", "matheus@email.com").
				AddRow(idTwo, "Mendes", "mendes@email.com"),
		)

	queries := db.New(mockDB)
	userRepo := repository.NewUserRepository(queries)
	service := NewUserService(userRepo)

	users, err := service.GetUsersByIDs(context.Background(), ids)
	if err != nil {
		t.Fatalf("GetUsersByIDs() returned unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("GetUsersByIDs() returned %d users, want 2", len(users))
	}
	if users[0].Name != "Matheus" {
		t.Fatalf("GetUsersByIDs() first user name = %s, want %s", users[0].Name, "Matheus")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestUserService_GetUsersByIDs_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	ids := []uuid.UUID{uuid.New()}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users WHERE id")).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(AppErr.ErrUserNotFound)

	queries := db.New(mockDB)
	userRepo := repository.NewUserRepository(queries)
	service := NewUserService(userRepo)

	users, err := service.GetUsersByIDs(context.Background(), ids)
	if users != nil {
		t.Fatalf("GetUsersByIDs() expected nil users, got %#v", users)
	}
	if !errors.Is(err, AppErr.ErrUserNotFound) {
		t.Fatalf("GetUsersByIDs() error = %v, want %v", err, AppErr.ErrUserNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}

func TestUserService_GetUsersByIDs_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	ids := []uuid.UUID{uuid.New()}
	dbErr := errors.New("fatal database connection failure")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users WHERE id")).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(dbErr)

	queries := db.New(mockDB)
	userRepo := repository.NewUserRepository(queries)
	service := NewUserService(userRepo)

	users, err := service.GetUsersByIDs(context.Background(), ids)
	if users != nil {
		t.Fatalf("GetUsersByIDs() expected nil users, got %#v", users)
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("GetUsersByIDs() error = %v, want %v", err, dbErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations were not met: %v", err)
	}
}
