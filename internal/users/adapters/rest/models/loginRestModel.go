package models

type LoginRestRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRestResponse struct {
	TokenHash string
	User      UserResponse
}
