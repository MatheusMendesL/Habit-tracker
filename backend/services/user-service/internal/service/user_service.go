package service

import (
	"context"
	"errors"
	"user-service/db"
	AppErr "user-service/internal/errors"
	"user-service/internal/repository"

	"github.com/google/uuid"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*db.User, error) {
	user, err := s.repo.FindByID(ctx, id)

	if err != nil {
		if errors.Is(err, AppErr.ErrUserNotFound) {
			return nil, AppErr.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) SearchUser(ctx context.Context, name string, email string) ([]*db.User, error) {
	users, err := s.repo.SearchUser(ctx, name, email)

	if err != nil {
		if errors.Is(err, AppErr.ErrUserNotFound) {
			return nil, AppErr.ErrUserNotFound
		}
		return nil, err
	}

	return users, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// ok
	return s.repo.DeleteUser(ctx, id)
}

func (s *UserService) EditUser(ctx context.Context, params db.UpdateUserParams) (*db.User, error) {
	user, err := s.repo.EditUser(ctx, params)

	if err != nil {
		if errors.Is(err, AppErr.ErrUserNotFound) {
			return nil, AppErr.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) EditPassword(ctx context.Context, pass *db.UpdatePasswordParams) error {
	// ok
	return s.repo.EditPassword(ctx, pass)
}

func (s *UserService) GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]db.GetUsersByIDsRow, error) {
	users, err := s.repo.GetUsersByIDs(ctx, ids)

	if err != nil {
		if errors.Is(err, AppErr.ErrUserNotFound) {
			return nil, AppErr.ErrUserNotFound
		}
		return nil, err
	}

	return users, nil
}
