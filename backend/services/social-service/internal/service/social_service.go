package service

import "social/internal/repository"

type SocialService struct {
	repo *repository.SocialRepository
}

func NewSocialService(r *repository.SocialRepository) *SocialService {
	return &SocialService{repo: r}
}
