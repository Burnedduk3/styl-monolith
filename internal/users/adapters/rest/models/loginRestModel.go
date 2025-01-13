package models

type LoginRestRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRestResponse struct {
	TokenHash string
	User      UserResponse
}

type TokenRequest struct {
	Email   string `json:"email"`
	TokenId string `json:"token_id"`
}
