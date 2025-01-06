package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"styl-monolith/internal/users/adapters/awsAuth"
	"styl-monolith/internal/users/adapters/postgres"
	userRest "styl-monolith/internal/users/adapters/rest"
	userService "styl-monolith/internal/users/core/services"
	"styl-monolith/pkg/aws"
	postgresClient "styl-monolith/pkg/postgres"
)

// BuildUserCrudDomainHandler sets up and returns a user CRUD domain handler with all dependencies injected properly.
// It initializes database connection, repositories, and services required to handle user CRUD operations.
func BuildUserCrudDomainHandler(log *logrus.Logger) userRest.CrudRest {
	databaseClient := postgresClient.GetDatabaseInstance(log)
	userCrudRepository := postgres.NewUserRepository(log, databaseClient.GetDB())
	roleCrudRepository := postgres.NewRoleRepository(log, databaseClient.GetDB())
	crudService := userService.NewCrudService(log, userCrudRepository, roleCrudRepository)
	handler := userRest.NewCrudRestUser(crudService, log)
	return handler
}

// BuildLoginUserDomainHandler constructs and returns a LoginRest handler for managing user login-related functionalities.
func BuildLoginUserDomainHandler(log *logrus.Logger) userRest.LoginRest {
	awsRegion := viper.GetString("AWS_REGION")
	dynamoClient := aws.GetDynamoClientInstance(log, awsRegion).GetClient()
	cognitoClient := aws.GetCognitoClientInstance(log, awsRegion).GetClient()
	awsAuthRepository := awsAuth.NewAwsAuthRepository(log, dynamoClient, cognitoClient)
	loginService := userService.NewLoginService(log, awsAuthRepository)
	handler := userRest.NewLoginRestUser(&loginService, log)
	return handler
}
