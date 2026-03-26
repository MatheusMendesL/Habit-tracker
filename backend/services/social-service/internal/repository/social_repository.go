package repository

import "social/db"

type SocialRepository struct {
	q *db.Queries
}

func NewSocialRepository(q *db.Queries) *SocialRepository {
	return &SocialRepository{q: q}
}
