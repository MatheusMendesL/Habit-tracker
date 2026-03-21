package repository

import (
	"context"
	"user-service/db"
)

type UserRepository struct {
	q *db.Queries
}

func NewUserRepository(q *db.Queries) *UserRepository {
	return &UserRepository{q: q}
}

func (r *UserRepository) FindByID(ctx context.Context, id int32) (*db.User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &db.User{
		ID:    row.ID,
		Name:  row.Name,
		Email: row.Email,
	}, nil
}

func (r *UserRepository) SearchUser(ctx context.Context, name string, email string) ([]*db.User, error) {
	params := db.SearchUserParams{
		Name:  name,
		Email: email,
	}
	rows, err := r.q.SearchUser(ctx, params)

	if err != nil {
		return nil, err
	}

	users := make([]*db.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, &db.User{
			ID:    row.ID,
			Name:  row.Name,
			Email: row.Email,
		})
	}
	return users, nil
}
