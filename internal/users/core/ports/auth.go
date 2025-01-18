package ports

import "styl-monolith/internal/users/core/domain"

type LoginPort interface {
	SignUpInCognito(user domain.User) (domain.User, error)
	SignInCognito(email, password string) (domain.UserAuth, error)
	SaveTokensToDynamoDB(userAuth domain.UserAuth) error
	FetchUserFromDatabase(email, tokenHash string) (domain.User, error)
	DeleteUserFromDatabase(id uint) error
	DeleteUserInCognito(email string) error
	GetTokensFromDynamoById(tokenId string) (domain.UserAuth, error)
	GetTokensFromDynamoByIdAndEmail(tokenId, email string) (domain.UserAuth, error)
	GetTokensFromDynamoByEmail(email string) (domain.UserAuth, error)
	RefreshTokensWithCognito(userAuth domain.UserAuth, username string) (domain.UserAuth, error)
	SignOutCognito(userAuth domain.UserAuth) error
	DeleteTokensFromDynamoById(userAuth domain.UserAuth) error
}
