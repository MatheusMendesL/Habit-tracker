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

func (s *UserService) SearchUser(ctx context.Context, name string, email string) (*db.User, error) {
	user, err := s.repo.SearchUser(ctx, name, email)

	if err != nil {
		return &db.User{}, err
	}

	return &db.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
