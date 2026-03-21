package service

import (
	"context"
	model "user-service/internal/model"
	"user-service/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) GetUserByID(ctx context.Context, id int32) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) SearchUser(ctx context.Context, name string, email string) (*model.User, error) {
	return s.repo.SearchUser(ctx, name, email)
}
