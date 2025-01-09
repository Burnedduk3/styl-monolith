package models

type Tokens struct {
	TokenId        string
	TokenCognitoId string
	UserId         string
	AccessToken    string
	RefreshToken   string
	ExpiresIn      int
}
