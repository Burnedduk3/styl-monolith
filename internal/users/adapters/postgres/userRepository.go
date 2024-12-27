package postgres

import "github.com/sirupsen/logrus"

type UserRepository struct {
	log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{log: log}
}
