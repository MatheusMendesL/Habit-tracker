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

func (r *UserRepository) DeleteUser(ctx context.Context, id int32) error {
	err := r.q.DeleteUser(ctx, id)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, req db.UpdateUserParams) (*db.User, error) {
	err := r.q.UpdateUser(ctx, req)

	if err != nil {
		return nil, err
	}

	user, err := r.FindByID(ctx, req.ID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

/*func (r *UserRepository) StartFollowing(ctx context.Context, followerID, followeeID int32) error {
	params := db.StartFollowingParams{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
	return r.q.StartFollowing(ctx, params)

}

func (r *UserRepository) StopFollowing(ctx context.Context, followerID, followeeID int32) error {
	params := db.UnfollowParams{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
	return r.q.Unfollow(ctx, params)

}

func (r *UserRepository) ListFollowers(ctx context.Context, userID int32) ([]*db.User, error) {
	rows, err := r.q.ListFollowers(ctx, userID)
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

func (r *UserRepository) ListFollowing(ctx context.Context, userID int32) ([]*db.User, error) {
	rows, err := r.q.ListFollowing(ctx, userID)
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
}*/

func (r *UserRepository) UpdatePassword(ctx context.Context, req *db.UpdatePasswordParams) error {
	params := db.UpdatePasswordParams{
		ID:       req.ID,
		Password: req.Password,
	}

	return r.q.UpdatePassword(ctx, params)
}

func (r *UserRepository) GetUsersByIDs(ctx context.Context, ids []int32) ([]db.GetUsersByIDsRow, error) {
	if len(ids) == 0 {
		return []db.GetUsersByIDsRow{}, nil
	}

	users, err := r.q.GetUsersByIDs(ctx, ids)

	if err != nil {
		return []db.GetUsersByIDsRow{}, err
	}

	m := make(map[int32]db.GetUsersByIDsRow, len(users))
	for _, u := range users {
		m[u.ID] = u
	}

	ordered := make([]db.GetUsersByIDsRow, 0, len(ids))
	for _, id := range ids {
		if u, ok := m[id]; ok {
			ordered = append(ordered, u)
		}
	}

	return ordered, nil
}
