package dto

type FollowRequest struct {
	FollowerID string `json:"-"`
	FolloweeID string `json:"followee_id"`
}

type FollowResponse struct {
	Success bool `json:"success"`
}

type SocialUsersResponse struct {
	Users []*User `json:"users"`
}
