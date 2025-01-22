package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"styl-monolith/internal/users/adapters/awsAuth"
	"styl-monolith/internal/users/adapters/postgres"
	userRest "styl-monolith/internal/users/adapters/rest"
	"styl-monolith/internal/users/core/ports"
	userService "styl-monolith/internal/users/core/services"
	"styl-monolith/pkg/aws"
	postgresClient "styl-monolith/pkg/postgres"
)

// Constants for configuration keys
const (
	AwsRegionKey = "AWS_REGION"
)

// initializeDatabaseComponents sets up the database client and repositories (user & role).
func initializeDatabaseComponents(log *logrus.Logger) (*gorm.DB, ports.UserPort, ports.RolePort, ports.UserReportRepository) {
	databaseClient := postgresClient.GetDatabaseInstance(log)
	db := databaseClient.GetDB()
	userRepository := postgres.NewUserRepository(log, db)
	roleRepository := postgres.NewRoleRepository(log, db)
	userReportRepository := postgres.NewUserReportRepository(log, db)
	return db, userRepository, roleRepository, userReportRepository
}

// BuildUserCrudDomainHandler creates and configures the User CRUD domain handler with its dependencies.
func BuildUserCrudDomainHandler(log *logrus.Logger) userRest.CrudRest {
	_, userRepo, roleRepo, userReportRepository := initializeDatabaseComponents(log)
	crudService := userService.NewCrudService(log, userRepo, roleRepo, userReportRepository)
	return userRest.NewCrudRestUser(crudService, log)
}

// BuildLoginUserDomainHandler constructs and returns the login domain handler with AWS auth support.
func BuildLoginUserDomainHandler(log *logrus.Logger) userRest.LoginRest {
	_, userRepo, roleRepo, userReportRepository := initializeDatabaseComponents(log)
	crudService := userService.NewCrudService(log, userRepo, roleRepo, userReportRepository)
	awsAuthRepo := CreateAwsAuthClient(log)
	loginService := userService.NewLoginService(log, awsAuthRepo)
	return userRest.NewLoginRestUser(loginService, crudService, log)
}

// CreateAwsAuthClient initializes and returns an AWS-based authentication client using Cognito and DynamoDB services.
func CreateAwsAuthClient(log *logrus.Logger) ports.LoginPort {
	awsRegion := viper.GetString(AwsRegionKey)
	dynamoClient := aws.GetDynamoClientInstance(log, awsRegion).GetClient()
	cognitoClient := aws.GetCognitoClientInstance(log, awsRegion).GetClient()
	return awsAuth.NewAwsAuthRepository(log, dynamoClient, cognitoClient)
}
