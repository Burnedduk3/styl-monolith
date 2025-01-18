package services

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/sirupsen/logrus"
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/internal/users/core/ports"
	"styl-monolith/pkg/utils/jwt"
)

type LoginServiceStruct struct {
	logger  *logrus.Logger
	awsAuth ports.LoginPort
}

type LoginService interface {
	Login(email, password string) (domain.UserAuth, domain.User, error)
	SignUp(domain.User) (domain.User, error)
	Logout(tokenId, email string) error
	Refresh(tokenId, email string) (domain.UserAuth, error)
	ChangePassword() error
	DeleteUser(email, tokenHash string) error
}

func NewLoginService(log *logrus.Logger, loginPort ports.LoginPort) LoginService {
	return &LoginServiceStruct{
		logger:  log,
		awsAuth: loginPort,
	}
}

func (l *LoginServiceStruct) Login(email, password string) (domain.UserAuth, domain.User, error) {
	l.logger.Debug("Starting login method")
	l.logger.Debug("signing in user with email: ", email)
	userAuth, err := l.awsAuth.SignInCognito(email, password)
	if err != nil {
		return domain.UserAuth{}, domain.User{}, err
	}
	h := sha256.New()
	h.Write([]byte(userAuth.IdToken))
	idTokenHash := base64.StdEncoding.EncodeToString(h.Sum(nil))
	userAuth.IdTokenHash = idTokenHash
	l.logger.Debug("SavingTokenToDatabase")
	err = l.awsAuth.SaveTokensToDynamoDB(userAuth)
	if err != nil {
		return userAuth, domain.User{}, err
	}
	l.logger.Debug("Fetching User from database with email: ", email)
	user, err := l.awsAuth.FetchUserFromDatabase(email, userAuth.IdTokenHash)
	if err != nil {
		return userAuth, domain.User{}, err
	}
	l.logger.Debug("Tokens Successfully saved to database saved to database")
	return userAuth, user, nil
}

func (l *LoginServiceStruct) SignUp(user domain.User) (domain.User, error) {
	createdUser, err := l.awsAuth.SignUpInCognito(user)
	if err != nil {
		return domain.User{}, err
	}
	return createdUser, nil
}

func (l *LoginServiceStruct) Logout(tokenId, email string) error {
	l.logger.Debug("Starting logout method")
	l.logger.Debug("Fetching tokens from database with tokenId: ", tokenId)
	userAuth, err := l.awsAuth.GetTokensFromDynamoByIdAndEmail(tokenId, email)
	if err != nil {
		l.logger.Debug("Error fetching tokens from database with tokenId: ", tokenId)
		return err
	}
	l.logger.Debug("signin out user with tokenId: ", tokenId)
	err = l.awsAuth.SignOutCognito(userAuth)
	if err != nil {
		l.logger.Debug("Error signing Out from cognito database: ", tokenId)
		return err
	}
	err = l.awsAuth.DeleteTokensFromDynamoById(userAuth)
	if err != nil {
		l.logger.Debug("Error deleting tokens from database with tokenId: ", tokenId)
	}
	return nil
}

func (l *LoginServiceStruct) Refresh(tokenId, email string) (domain.UserAuth, error) {
	l.logger.Debug("Starting refresh method")
	l.logger.Debug("Fetching tokens from database with tokenId: ", tokenId)
	userAuth, err := l.awsAuth.GetTokensFromDynamoByIdAndEmail(tokenId, email)
	if err != nil {
		l.logger.Debug("Error fetching tokens from database with tokenId: ", tokenId)
		return userAuth, err
	}
	jwtClaims, err := jwt.DecodeIdJWT(userAuth.IdToken)
	if err != nil {
		l.logger.Error("Unable to decode JWT")
		return userAuth, err
	}
	if jwtClaims.Email != email {
		l.logger.Error("Emails do not match, request email: ", email, " Token Email ", userAuth.Email)
		return userAuth, err
	}
	l.logger.Debug("Refreshing tokens with tokenId: ", tokenId)
	NewUserAuth, err := l.awsAuth.RefreshTokensWithCognito(userAuth, jwtClaims.Username)
	NewUserAuth.RefreshToken = userAuth.RefreshToken
	if err != nil {
		l.logger.Debug("Error refreshing tokens with tokenId: ", tokenId)
		return userAuth, err
	}
	l.logger.Debug("Updating tokens in database with tokenId: ", tokenId)
	err = l.awsAuth.DeleteTokensFromDynamoById(userAuth)
	if err != nil {
		l.logger.Debug("Error deleting tokens from database with tokenId: ", tokenId)
		return userAuth, err
	}
	h := sha256.New()
	h.Write([]byte(NewUserAuth.IdToken))
	idTokenHash := base64.StdEncoding.EncodeToString(h.Sum(nil))
	NewUserAuth.IdTokenHash = idTokenHash
	err = l.awsAuth.SaveTokensToDynamoDB(NewUserAuth)
	if err != nil {
		l.logger.Debug("Error saving tokens to database with tokenId: ", tokenId)
		return userAuth, err
	}
	return NewUserAuth, nil
}

func (l *LoginServiceStruct) ChangePassword() error {
	// TODO: Implement ChangePassword logic
	return nil
}

func (l *LoginServiceStruct) DeleteUser(email, tokenHash string) error {
	l.logger.Debug("Starting delete method")
	l.logger.Debug("Fetching User from database with email: ", email)
	user, err := l.awsAuth.FetchUserFromDatabase(email, tokenHash)
	if err != nil {
		l.logger.Debug("Error fetching User from database with email: ", email)
		return err
	}
	l.logger.Debug("Deleting User from database with id: ", user.ID)
	err = l.awsAuth.DeleteUserFromDatabase(user.ID)
	if err != nil {
		l.logger.Debug("Error fetching User from database with email: ", email)
		return err
	}
	err = l.awsAuth.DeleteUserInCognito(user.Email)
	if err != nil {
		l.logger.Debug("Error fetching User from database with email: ", email)
		return err
	}
	return nil
}
