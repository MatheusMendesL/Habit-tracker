package repository

import (
	"context"
	"user-service/db"

	model "user-service/internal/model"
)

type UserRepository struct {
	q *db.Queries
}

func NewUserRepository(q *db.Queries) *UserRepository {
	return &UserRepository{q: q}
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
