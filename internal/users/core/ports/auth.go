package ports

import "styl-monolith/internal/users/core/domain"

type LoginPort interface {
	SignUpInCognito(user domain.User) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	SignInCognito(email, password string) (string, error)
	SaveTokensToDynamoDB(token string, user domain.User) error
	GetTokensFromDynamoDB(token string) (string, string, error)
	RefreshTokensWithCognito(refreshToken string) (string, string, error)
	SignOutCognito(token string) error
	DeleteTokensFromDynamoDB(token string) error
}
