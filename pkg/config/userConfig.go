package config

import (
	"github.com/sirupsen/logrus"
	"styl-monolith/internal/users/adapters/postgres"
	userRest "styl-monolith/internal/users/adapters/rest"
	userService "styl-monolith/internal/users/core/services"
	"styl-monolith/pkg/aws"
)

func BuildUserCrudDomainHandler(log *logrus.Logger) userRest.CrudRest {
	databaseClient := aws.GetDatabaseInstance(log)
	userCrudRepository := postgres.NewUserRepository(log, databaseClient.GetDB())
	roleCrudRepository := postgres.NewRoleRepository(log, databaseClient.GetDB())
	crudService := userService.NewCrudService(log, userCrudRepository, roleCrudRepository)
	handler := userRest.NewCrudRestUser(crudService, log)
	return handler
}

func BuildLoginUserDomainHandler(log *logrus.Logger) userRest.LoginRest {
	loginService := userService.NewLoginService(log)
	handler := userRest.NewLoginRestUser(&loginService, log)
	return handler
}
