package service

import (
	"context"
	"database/sql"
	"errors"
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

	mock.ExpectQuery(`SELECT id, name, email FROM users WHERE id = \$1`).
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

	mock.ExpectQuery(`SELECT id, name, email FROM users WHERE id = \$1`).
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

	mock.ExpectQuery(`SELECT id, name, email FROM users WHERE`).
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

	mock.ExpectQuery(`SELECT id, name, email FROM users WHERE`).
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

	mock.ExpectExec(`UPDATE users SET`).
		WithArgs(params.Name, params.Email, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(`SELECT id, name, email FROM users WHERE id = \$1`).
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
