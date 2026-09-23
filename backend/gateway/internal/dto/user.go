package dto

type GetUserByIDRequest struct {
	ID string `json:"id"`
}

type GetUserByIDResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
