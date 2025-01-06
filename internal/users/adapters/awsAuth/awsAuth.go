package awsAuth

import (
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/sirupsen/logrus"
	"styl-monolith/internal/users/core/domain"
)

type AwsAuth struct {
	log           *logrus.Logger
	dynamoClient  *dynamodb.Client
	cognitoClient *cognitoidentityprovider.Client
}

func NewAwsAuthRepository(log *logrus.Logger, dynamoClient *dynamodb.Client, cognitoClient *cognitoidentityprovider.Client) *AwsAuth {
	return &AwsAuth{log: log, dynamoClient: dynamoClient, cognitoClient: cognitoClient}
}

func (a *AwsAuth) SignUpInCognito(user domain.User) (string, error) {
	// TODO: Implement logic to sign up user in Cognito using a.cognitoClient
	// Example steps:
	// - Use CognitoIdentityProvider's SignUp API to register the user
	// - Return user ID and/or credentials from the Cognito response
	// - Handle errors
	return "", nil
}

func (a *AwsAuth) GetUserByEmail(email string) (domain.User, error) {
	// TODO: Implement logic to fetch user information by email
	// Example steps:
	// - Query Cognito or DynamoDB to fetch user details (depending on the implementation)
	// - Populate and return the User domain object
	// - Handle errors
	return domain.User{}, nil
}

func (a *AwsAuth) SignInCognito(email, password string) (string, error) {
	// TODO: Implement logic to sign in the user via Cognito
	// Example steps:
	// - Use CognitoIdentityProvider's AdminInitiateAuth or InitiateAuth API
	// - Pass email and password to authenticate the user
	// - Return generated access token or error
	return "", nil
}

func (a *AwsAuth) SaveTokensToDynamoDB(token string, user domain.User) error {
	// TODO: Implement logic to save tokens in DynamoDB
	// Example steps:
	// - Use DynamoDB PutItem API to save the token and user information
	// - Handle errors
	return nil
}

func (a *AwsAuth) GetTokensFromDynamoDB(token string) (string, string, error) {
	// TODO: Implement logic to retrieve tokens from DynamoDB
	// Example steps:
	// - Use DynamoDB GetItem API to fetch the tokens associated with the provided key
	// - Return access and refresh tokens (or error if not found)
	return "", "", nil
}

func (a *AwsAuth) RefreshTokensWithCognito(refreshToken string) (string, string, error) {
	// TODO: Implement logic to refresh tokens using Cognito
	// Example steps:
	// - Use CognitoIdentityProvider's InitiateAuth or AdminInitiateAuth API
	// - Pass the refresh token to generate new access and refresh tokens
	// - Return new tokens or error
	return "", "", nil
}

func (a *AwsAuth) SignOutCognito(token string) error {
	// TODO: Implement logic to sign out user in Cognito
	// Example steps:
	// - Use CognitoIdentityProvider's GlobalSignOut or RevokeToken API
	// - Pass the access token to sign the user out
	// - Handle errors
	return nil
}

func (a *AwsAuth) DeleteTokensFromDynamoDB(token string) error {
	// TODO: Implement logic to delete tokens from DynamoDB
	// Example steps:
	// - Use DynamoDB DeleteItem API to remove the token and its associated data
	// - Handle errors
	return nil
}
