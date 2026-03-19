package repository

import (
	"context"
	model "user-service/internal/model"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindByID(ctx context.Context, id int32) (*model.User, error) {
	return &model.User{
		ID:    id,
		Name:  "teste",
		Email: "teste@email.com",
	}, nil
}
