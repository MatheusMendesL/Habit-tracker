package service

import (
	"habit-service/internal/repository"
	pbUser "shared/pb/user"
)

type HabitService struct {
	pbUser.UserServiceClient
	repo *repository.HabitRepository
}

func NewHabitService(r *repository.HabitRepository, userClient pbUser.UserServiceClient) *HabitService {
	return &HabitService{
		repo:              r,
		UserServiceClient: userClient,
	}
}
