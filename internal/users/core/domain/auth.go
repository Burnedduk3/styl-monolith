package domain

type UserAuth struct {
	IdTokenHash       string
	AccessToken       string
	RefreshToken      string
	IdToken           string
	Email             string
	ExpiryAccessToken int32
}
