package domain

type UserAuth struct {
	IdTokenHash  string
	AccessToken  string
	RefreshToken string
	IdToken      string
	UserId       uint
	Email        string
}

type JwtClaims struct {
	Username string `json:"cognito:username"`
	Email    string `json:"email"`
}
