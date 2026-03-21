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
	// vai pegar baseado nos sql q vão ser gerados

	return &model.User{
		ID:    id,
		Name:  "teste",
		Email: "teste@email.com",
	}, nil
}

func (r *UserRepository) SearchUser(ctx context.Context, name string, email string) (*model.User, error) {
	// fazer busca no banco, retorno vai pro return

	return nil, nil
}
