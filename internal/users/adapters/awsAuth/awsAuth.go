package awsAuth

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"styl-monolith/internal/users/adapters/rest/models"
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/pkg/errorhandler"
	"styl-monolith/pkg/logger"
	"styl-monolith/pkg/utils"
)

type AwsAuth struct {
	log           *logrus.Logger
	dynamoClient  *dynamodb.Client
	cognitoClient *cognitoidentityprovider.Client
	httpClient    *http.Client
	baseUrl       string
}

func NewAwsAuthRepository(log *logrus.Logger, dynamoClient *dynamodb.Client, cognitoClient *cognitoidentityprovider.Client, client *http.Client, baseUrl string) *AwsAuth {
	return &AwsAuth{log: log, dynamoClient: dynamoClient, cognitoClient: cognitoClient, httpClient: client, baseUrl: baseUrl}
}

func (a *AwsAuth) SignUpInCognito(user domain.User) (domain.User, error) {
	// Prepare user attributes for the new user
	userPoolID := viper.GetString("AWS_COGNITO_USER_POOL_ID")
	userAttributes := []types.AttributeType{
		{
			Name:  aws.String("email"),
			Value: aws.String(user.Email),
		},
		{
			Name:  aws.String("phone_number"),
			Value: aws.String(fmt.Sprintf("%s%s", user.CountryCode, user.Phone)),
		},
		{
			Name:  aws.String("birthdate"),
			Value: aws.String(user.Birthday),
		},
		{
			Name:  aws.String("given_name"),
			Value: aws.String(user.Name),
		},
		{
			Name:  aws.String("family_name"),
			Value: aws.String(user.LastName),
		},
		{
			Name:  aws.String("email_verified"), // Ensure email is marked as verified
			Value: aws.String("true"),
		},
		{
			Name:  aws.String("phone_number_verified"), // Ensure phone number is marked as verified
			Value: aws.String("true"),
		},
	}

	tempPassword := utils.GenerateRandomString(8)

	// Create the AdminCreateUserInput
	input := &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId:        &userPoolID,    // Cognito User Pool ID
		Username:          &user.Username, // The user's username
		UserAttributes:    userAttributes,
		TemporaryPassword: aws.String(tempPassword),        // Optional: Set a temporary password
		MessageAction:     types.MessageActionTypeSuppress, // Suppress welcome email (optional)
	}
	a.log.Info(fmt.Sprintf(logger.CreatingUserWithEmail, user.Email))
	// Call Cognito AdminCreateUser API
	_, err := a.cognitoClient.AdminCreateUser(context.TODO(), input)
	if err != nil {

		return domain.User{}, errorhandler.NewDomainError(
			errorhandler.ErrCreatingUserInCognitoUserPool,
			errorhandler.GetErrorMessage(errorhandler.ErrCreatingUserInCognitoUserPool),
			err)
	}
	a.log.Info(logger.SettingUserPassword)
	// Confirm the new password as definitive
	_, err = a.cognitoClient.AdminSetUserPassword(context.TODO(), &cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: &userPoolID,
		Username:   &user.Username,
		Password:   aws.String(user.Password), // Set the definitive password
		Permanent:  true,                      // Mark the password as permanent
	})
	if err != nil {
		return domain.User{}, errorhandler.NewDomainError(
			errorhandler.ErrSettingPermanentPassword,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrCreatingUserInCognitoUserPool), user.Username),
			err)
	}
	a.log.Info(logger.SuccessfullyCreatedCognitoUser)
	a.log.Info(fmt.Sprintf(logger.CreateUserOnDB, user.Email))
	payloadBytes, err := json.Marshal(user.UserToUserPayload())
	if err != nil {
		return domain.User{}, errorhandler.NewDomainError(
			errorhandler.ErrSettingPermanentPassword,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrCreatingUserInCognitoUserPool), user.Username),
			err)
	}
	// Create a new POST request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", a.baseUrl, "api/v1/user"), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return domain.User{}, errorhandler.NewDomainError(
			errorhandler.ErrSettingPermanentPassword,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrCreatingUserInCognitoUserPool), user.Username),
			err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return domain.User{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.User{}, err
	}

	// Check for non-2xx status codes and handle errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return domain.User{}, err
	}

	umarshaledUser := models.UserResponse{}
	_ = json.Unmarshal(body, &umarshaledUser)
	return umarshaledUser.ToDomainUser(), nil
}

// calculateSecretHash calculates Cognito's required SECRET_HASH
func (a *AwsAuth) calculateSecretHash(clientSecret, clientId, username string) string {
	h := hmac.New(sha256.New, []byte(clientSecret))
	h.Write([]byte(username + clientId))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (a *AwsAuth) SignInCognito(email, password string) (domain.UserAuth, error) {
	// Initiate the authentication request
	clientId := viper.GetString("AWS_COGNITO_USER_POOL_CLIENT_ID")
	clientSecret := viper.GetString("AWS_COGNITO_USER_POOL_CLIENT_SECRET")
	secretHash := a.calculateSecretHash(clientSecret, clientId, email)
	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(clientId),
		AuthParameters: map[string]string{
			"USERNAME":    email,
			"PASSWORD":    password,
			"SECRET_HASH": secretHash,
		},
	}
	authResponse, err := a.cognitoClient.InitiateAuth(context.TODO(), authInput)
	if err != nil {
		return domain.UserAuth{}, fmt.Errorf("authentication failed: %v", err)
	}
	userAuth := domain.UserAuth{
		AccessToken:  *authResponse.AuthenticationResult.AccessToken,
		RefreshToken: *authResponse.AuthenticationResult.RefreshToken,
		IdToken:      *authResponse.AuthenticationResult.IdToken,
		UserId:       0,
		Email:        email,
	}
	return userAuth, nil
}

func (a *AwsAuth) FetchUserFromDatabase(email string) (domain.User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", a.baseUrl, fmt.Sprintf("api/v1/user?email=%s", email)), nil)
	if err != nil {
		return domain.User{}, errorhandler.NewDomainError(
			errorhandler.ErrUserNotFound,
			errorhandler.GetErrorMessage(errorhandler.ErrUserNotFound),
			err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return domain.User{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.User{}, err
	}

	// Check for non-2xx status codes and handle errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return domain.User{}, err
	}

	umarshaledUser := models.UserResponse{}
	_ = json.Unmarshal(body, &umarshaledUser)
	return umarshaledUser.ToDomainUser(), nil
}

func (a *AwsAuth) SaveTokensToDynamoDB(userAuth domain.UserAuth) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(viper.GetString("AWS_DYNAMODB_TOKENS_TABLE")),
		Item: map[string]dynamodbTypes.AttributeValue{
			"id_token_cognito": &dynamodbTypes.AttributeValueMemberS{Value: userAuth.IdToken},
			"access_token":     &dynamodbTypes.AttributeValueMemberS{Value: userAuth.AccessToken},
			"refresh_token":    &dynamodbTypes.AttributeValueMemberS{Value: userAuth.RefreshToken},
			"token_id":         &dynamodbTypes.AttributeValueMemberS{Value: userAuth.IdTokenHash},
			"user_id":          &dynamodbTypes.AttributeValueMemberN{Value: fmt.Sprintf("%d", userAuth.UserId)},
			"user_email":       &dynamodbTypes.AttributeValueMemberS{Value: userAuth.Email},
		},
	}

	_, err := a.dynamoClient.PutItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to save tokens to DynamoDB: %v", err)
	}
	return nil
}

func (a *AwsAuth) GetTokensFromDynamoDB(tokenId string) (domain.UserAuth, error) {
	// TODO: Implement logic to retrieve tokens from DynamoDB
	// Example steps:
	// - Use DynamoDB GetItem API to fetch the tokens associated with the provided key
	// - Return access and refresh tokens (or error if not found)
	return domain.UserAuth{}, nil
}

func (a *AwsAuth) GetTokensFromDynamoByEmail(email string) (domain.UserAuth, error) {
	// TODO: Implement logic to retrieve tokens from DynamoDB
	// Example steps:
	// - Use DynamoDB GetItem API to fetch the tokens associated with the provided key
	// - Return access and refresh tokens (or error if not found)
	return domain.UserAuth{}, nil
}

func (a *AwsAuth) RefreshTokensWithCognito(userAuth domain.UserAuth) (domain.UserAuth, error) {
	// TODO: Implement logic to refresh tokens using Cognito
	// Example steps:
	// - Use CognitoIdentityProvider's InitiateAuth or AdminInitiateAuth API
	// - Pass the refresh token to generate new access and refresh tokens
	// - Return new tokens or error
	return domain.UserAuth{}, nil
}

func (a *AwsAuth) SignOutCognito(tokenId string) error {
	// TODO: Implement logic to sign out user in Cognito
	// Example steps:
	// - Use CognitoIdentityProvider's GlobalSignOut or RevokeToken API
	// - Pass the access token to sign the user out
	// - Handle errors
	return nil
}

func (a *AwsAuth) DeleteTokensFromDynamoDB(tokenId string) error {
	// TODO: Implement logic to delete tokens from DynamoDB
	// Example steps:
	// - Use DynamoDB DeleteItem API to remove the token and its associated data
	// - Handle errors
	return nil
}
