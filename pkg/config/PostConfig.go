package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"styl-monolith/internal/media/adapters/fileManager"
	stylGrpc "styl-monolith/internal/media/adapters/grcp"
	mediaPostgresAdapter "styl-monolith/internal/media/adapters/postgres"
	mediaRestAdapter "styl-monolith/internal/media/adapters/rest"
	"styl-monolith/internal/media/core/services"
	"styl-monolith/pkg/aws"
	postgresClient "styl-monolith/pkg/postgres"
)

func BuildMediaCrudDomainHandler(log *logrus.Logger) mediaRestAdapter.CrudPostHandler {
	region := viper.GetString("AWS_REGION")
	databaseClient := postgresClient.GetDatabaseInstance(log)
	db := databaseClient.GetDB()
	fileMngmt := fileManager.NewFileManager(log)
	postgresMediaRepo := mediaPostgresAdapter.NewMediaCrudRepository(log, db)
	s3Client := aws.GetS3ClientInstance(log, region).GetClient()
	mediaService := services.NewMediaService(log, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, fileMngmt, s3Client)
	mediaHandler := mediaRestAdapter.NewCrudPostHandler(log, mediaService)
	return mediaHandler
}

func BuildMediaCrudDomainGrpc(log *logrus.Logger) stylGrpc.CrudPostGrcp {
	region := viper.GetString("AWS_REGION")
	databaseClient := postgresClient.GetDatabaseInstance(log)
	db := databaseClient.GetDB()
	fileMngmt := fileManager.NewFileManager(log)
	postgresMediaRepo := mediaPostgresAdapter.NewMediaCrudRepository(log, db)
	s3Client := aws.GetS3ClientInstance(log, region).GetClient()
	mediaService := services.NewMediaService(log, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, fileMngmt, s3Client)
	grpcMediaService := stylGrpc.NewCrudPostGrcp(mediaService, log)
	return *grpcMediaService
}
