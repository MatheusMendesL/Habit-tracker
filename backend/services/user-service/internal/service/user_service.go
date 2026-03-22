package service

import (
	"context"
	"user-service/db"
	"user-service/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) GetUserByID(ctx context.Context, id int32) (*db.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) SearchUser(ctx context.Context, name string, email string) ([]*db.User, error) {
	return s.repo.SearchUser(ctx, name, email)
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, user db.UpdateUserParams) (*db.User, error) {
	return s.repo.UpdateUser(ctx, user)
}
