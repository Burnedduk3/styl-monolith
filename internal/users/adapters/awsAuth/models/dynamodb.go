package models

type Tokens struct {
	TokenId      string
	UserId       string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}
