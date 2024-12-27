package config

import (
	"github.com/sirupsen/logrus"
	userRest "styl-monolith/internal/users/adapters/rest"
	userService "styl-monolith/internal/users/core/services"
	"styl-monolith/pkg/logger"
)

func BuildUserCrudDomainHandler() userRest.CrudRest {
	log := logger.GetLogger(logrus.DebugLevel)
	crudService := userService.NewCrudService(log)
	handler := userRest.NewCrudRestUser(&crudService, log)
	return handler
}

func BuildLoginUserDomainHandler() userRest.LoginRest {
	log := logger.GetLogger(logrus.DebugLevel)
	loginService := userService.NewLoginService(log)
	handler := userRest.NewLoginRestUser(&loginService, log)
	return handler
}
