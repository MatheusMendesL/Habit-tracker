package dto

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GetUserByIDRequest struct {
	ID string `json:"id"`
}

type GetUserByIDResponse struct {
	User *User `json:"user"`
}

type SearchUsersRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

type SearchUsersResponse struct {
	Users []*User `json:"users"`
}

type DeleteUserRequest struct {
	ID string `json:"id"`
}

type DeleteUserResponse struct {
	Success bool `json:"success"`
}

type GetUsersByIDsRequest struct {
	IDs []string `json:"ids"`
}

type GetUsersByIDsResponse struct {
	Users []*User `json:"users"`
}

type EditUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

type EditUserResponse struct {
	User *User `json:"user"`
}

type EditPasswordRequest struct {
	NewPassword string `json:"new_password"`
}

type EditPasswordResponse struct {
	Success bool `json:"success"`
}
