package postgres

import (
	"github.com/sirupsen/logrus"
	"styl-monolith/pkg/aws"
)

type RoleRepository struct {
	log      *logrus.Logger
	pgClient aws.Database
}

func NewRoleRepository(log *logrus.Logger) *RoleRepository {
	return &RoleRepository{log: log}
}
