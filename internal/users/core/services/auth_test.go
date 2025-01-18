package services_test

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/internal/users/core/services"
	"styl-monolith/pkg/utils/jwt"
	"testing"
)

// MockLoginPort is a mock implementation of the LoginPort interface for testing purposes.
// It embeds the mock.Mock type to leverage mock behavior from the testify/mock package.
type MockLoginPort struct {
	mock.Mock
}

// GetTokensFromDynamoById retrieves user authentication tokens from DynamoDB using the provided token ID.
func (m *MockLoginPort) GetTokensFromDynamoById(tokenId string) (domain.UserAuth, error) {
	args := m.Called(tokenId)
	return args.Get(0).(domain.UserAuth), args.Error(1)
}

// GetTokensFromDynamoByEmail retrieves tokens for a user from DynamoDB based on their email address.
func (m *MockLoginPort) GetTokensFromDynamoByEmail(email string) (domain.UserAuth, error) {
	args := m.Called(email)
	return args.Get(0).(domain.UserAuth), args.Error(1)
}

// SignInCognito authenticates a user in Cognito using the provided email and password, returning UserAuth data or an error.
func (m *MockLoginPort) SignInCognito(email, password string) (domain.UserAuth, error) {
	args := m.Called(email, password)
	return args.Get(0).(domain.UserAuth), args.Error(1)
}

// SaveTokensToDynamoDB stores the provided user authentication tokens in DynamoDB and returns an error if the operation fails.
func (m *MockLoginPort) SaveTokensToDynamoDB(userAuth domain.UserAuth) error {
	args := m.Called(userAuth)
	return args.Error(0)
}

// FetchUserFromDatabase retrieves a user from the database using the provided email and idTokenHash. Returns a User and an error.
func (m *MockLoginPort) FetchUserFromDatabase(email, idTokenHash string) (domain.User, error) {
	args := m.Called(email, idTokenHash)
	return args.Get(0).(domain.User), args.Error(1)
}

// GetTokensFromDynamoByIdAndEmail retrieves user authentication tokens from DynamoDB using token ID and email.
// Returns a UserAuth object and an error if the operation fails.
func (m *MockLoginPort) GetTokensFromDynamoByIdAndEmail(tokenId, email string) (domain.UserAuth, error) {
	args := m.Called(tokenId, email)
	return args.Get(0).(domain.UserAuth), args.Error(1)
}

// SignOutCognito attempts to sign out the user from the Amazon Cognito service by invalidating their authentication tokens.
func (m *MockLoginPort) SignOutCognito(userAuth domain.UserAuth) error {
	args := m.Called(userAuth)
	return args.Error(0)
}

// DeleteTokensFromDynamoById removes authentication tokens from DynamoDB for the provided user authorization data.
// Takes a domain.UserAuth object and returns an error if the operation fails.
func (m *MockLoginPort) DeleteTokensFromDynamoById(userAuth domain.UserAuth) error {
	args := m.Called(userAuth)
	return args.Error(0)
}

// RefreshTokensWithCognito refreshes user authentication tokens using Cognito based on provided user credentials and username.
func (m *MockLoginPort) RefreshTokensWithCognito(userAuth domain.UserAuth, username string) (domain.UserAuth, error) {
	args := m.Called(userAuth, username)
	return args.Get(0).(domain.UserAuth), args.Error(1)
}

// SignUpInCognito registers a new user in Cognito using the provided user details and returns the created user or an error.
func (m *MockLoginPort) SignUpInCognito(user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

// DeleteUserFromDatabase removes a user from the database using the provided user ID and returns an error if the operation fails.
func (m *MockLoginPort) DeleteUserFromDatabase(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// DeleteUserInCognito removes a user from the Cognito user pool using their email. Returns an error if the operation fails.
func (m *MockLoginPort) DeleteUserInCognito(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

// TestLoginService_Login tests the Login function of the LoginService to verify authentication and user retrieval behavior.
func TestLoginService_Login(t *testing.T) {
	mockLoginPort := new(MockLoginPort)
	logger := logrus.New()
	service := services.NewLoginService(logger, mockLoginPort)

	email := "test@example.com"
	password := "password"
	userAuth := domain.UserAuth{
		IdToken: "dummyToken",
	}
	user := domain.User{
		ID:    12345,
		Email: email,
	}

	mockLoginPort.On("SignInCognito", email, password).Return(userAuth, nil)
	mockLoginPort.On("SaveTokensToDynamoDB", mock.Anything).Return(nil)
	mockLoginPort.On("FetchUserFromDatabase", email, mock.Anything).Return(user, nil)

	resultAuth, resultUser, err := service.Login(email, password)

	assert.NoError(t, err)
	assert.Equal(t, userAuth.IdToken, resultAuth.IdToken)
	assert.Equal(t, user.Email, resultUser.Email)
	mockLoginPort.AssertExpectations(t)
}

// TestLoginService_SignUp tests the SignUp function of the LoginService to verify user signup functionality with mock dependencies.
// It checks for proper interaction with the SignUpInCognito method of MockLoginPort and validates the returned user details.
func TestLoginService_SignUp(t *testing.T) {
	mockLoginPort := new(MockLoginPort)
	logger := logrus.New()
	service := services.NewLoginService(logger, mockLoginPort)

	user := domain.User{
		ID:    1241241,
		Email: "test@example.com",
	}

	mockLoginPort.On("SignUpInCognito", user).Return(user, nil)

	createdUser, err := service.SignUp(user)

	assert.NoError(t, err)
	assert.Equal(t, user.ID, createdUser.ID)
	mockLoginPort.AssertExpectations(t)
}

// TestLoginService_Logout tests the Logout function of LoginService by simulating token retrieval, sign-out, and deletion.
func TestLoginService_Logout(t *testing.T) {
	mockLoginPort := new(MockLoginPort)
	logger := logrus.New()
	service := services.NewLoginService(logger, mockLoginPort)

	tokenId := "dummyTokenId"
	email := "test@example.com"
	userAuth := domain.UserAuth{
		IdToken: "dummyToken",
	}

	mockLoginPort.On("GetTokensFromDynamoByIdAndEmail", tokenId, email).Return(userAuth, nil)
	mockLoginPort.On("SignOutCognito", userAuth).Return(nil)
	mockLoginPort.On("DeleteTokensFromDynamoById", userAuth).Return(nil)

	err := service.Logout(tokenId, email)

	assert.NoError(t, err)
	mockLoginPort.AssertExpectations(t)
}

// TestLoginService_Refresh validates the functionality of refreshing user tokens via the LoginService.
// It ensures proper interaction with the mocked LoginPort, including fetching, refreshing, deleting, and saving tokens.
// Assertions verify that the service correctly returns the updated tokens and no errors occur during the operation.
func TestLoginService_Refresh(t *testing.T) {
	mockLoginPort := new(MockLoginPort)
	logger := logrus.New()
	service := services.NewLoginService(logger, mockLoginPort)

	tokenId := "dummyTokenId"
	email := "test@example.com"
	userAuth := domain.UserAuth{
		IdToken:      "dummyToken",
		Email:        email,
		RefreshToken: "dummyRefreshToken",
	}
	newUserAuth := domain.UserAuth{
		IdToken:      "newDummyToken",
		Email:        email,
		RefreshToken: "dummyRefreshToken",
	}
	jwtClaims := jwt.IdJwtClaims{
		Username: "testUser",
		Email:    email,
	}

	mockLoginPort.On("GetTokensFromDynamoByIdAndEmail", tokenId, email).Return(userAuth, nil)
	mockLoginPort.On("RefreshTokensWithCognito", userAuth, jwtClaims.Username).Return(newUserAuth, nil)
	mockLoginPort.On("DeleteTokensFromDynamoById", userAuth).Return(nil)
	mockLoginPort.On("SaveTokensToDynamoDB", newUserAuth).Return(nil)

	refreshedAuth, err := service.Refresh(tokenId, email)

	assert.NoError(t, err)
	assert.Equal(t, newUserAuth.IdToken, refreshedAuth.IdToken)
	mockLoginPort.AssertExpectations(t)
}
