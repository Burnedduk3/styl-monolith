package services

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/sirupsen/logrus"
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/internal/users/core/ports"
)

type LoginServiceStruct struct {
	logger  *logrus.Logger
	awsAuth ports.LoginPort
}

type LoginService interface {
	Login(email, password string) (domain.UserAuth, domain.User, error)
	SignUp(domain.User) (domain.User, error)
	Logout() error
	Refresh() error
	ChangePassword() error
}

func NewLoginService(log *logrus.Logger, loginPort ports.LoginPort) LoginService {
	return &LoginServiceStruct{
		logger:  log,
		awsAuth: loginPort,
	}
}

func (l *LoginServiceStruct) Login(email, password string) (domain.UserAuth, domain.User, error) {
	l.logger.Debug("Starting login method")
	l.logger.Debug("Fetching User from database with email: ", email)
	user, err := l.awsAuth.FetchUserFromDatabase(email)
	if err != nil {
		return domain.UserAuth{}, domain.User{}, err
	}
	l.logger.Debug("signing in user with email: ", email)
	userAuth, err := l.awsAuth.SignInCognito(email, password)
	if err != nil {
		return domain.UserAuth{}, domain.User{}, err
	}
	userAuth.UserId = user.ID
	h := sha256.New()
	h.Write([]byte(userAuth.IdToken))
	idTokenHash := base64.StdEncoding.EncodeToString(h.Sum(nil))
	userAuth.IdTokenHash = idTokenHash
	l.logger.Debug("User fetched from database, saving token to database")
	err = l.awsAuth.SaveTokensToDynamoDB(userAuth)
	if err != nil {
		return domain.UserAuth{}, domain.User{}, err
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

func (l *LoginServiceStruct) Logout() error {
	// TODO: Implement Logout logic
	return nil
}

func (l *LoginServiceStruct) Refresh() error {
	// TODO: Implement Refresh logic
	return nil
}

func (l *LoginServiceStruct) ChangePassword() error {
	// TODO: Implement ChangePassword logic
	return nil
}
