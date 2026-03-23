package service

import (
	"context"
	"user-service/db"
	AppErr "user-service/internal/errors"
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

func (s *UserService) StartFollowing(ctx context.Context, followerId int32, followeeId int32) error {
	if followerId == followeeId {
		return AppErr.ErrInvalidArgument
	}

	return s.repo.StartFollowing(ctx, followerId, followeeId)
}

func (s *UserService) StopFollowing(ctx context.Context, followerId int32, followeeId int32) error {
	return s.repo.StopFollowing(ctx, followerId, followeeId)
}

func (s *UserService) ListFollowers(ctx context.Context, userId int32) ([]*db.User, error) {
	return s.repo.ListFollowers(ctx, userId)
}

func (s *UserService) ListFollowing(ctx context.Context, userId int32) ([]*db.User, error) {
	return s.repo.ListFollowing(ctx, userId)
}
