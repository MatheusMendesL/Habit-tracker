package service

import (
	"context"
	"user-service/db"
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
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) SearchUser(ctx context.Context, name string, email string) ([]*db.User, error) {
	return s.repo.SearchUser(ctx, name, email)
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *UserService) EditUser(ctx context.Context, user db.UpdateUserParams) (*db.User, error) {
	return s.repo.EditUser(ctx, user)
}

func (s *UserService) EditPassword(ctx context.Context, pass *db.UpdatePasswordParams) error {
	return s.repo.EditPassword(ctx, pass)
}

func (s *UserService) GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]db.GetUsersByIDsRow, error) {
	return s.repo.GetUsersByIDs(ctx, ids)
}
