package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	awsAdapater "styl-monolith/internal/media/adapters/aws"
	stylGrpc "styl-monolith/internal/media/adapters/grcp"
	mediaPostgresAdapter "styl-monolith/internal/media/adapters/postgres"
	mediaRestAdapter "styl-monolith/internal/media/adapters/rest"
	"styl-monolith/internal/media/core/services"
	"styl-monolith/pkg/aws"
	postgresClient "styl-monolith/pkg/postgres"
)

func BuildMediaCrudDomainHandler(log *logrus.Logger) mediaRestAdapter.CrudPostHandler {
	region := viper.GetString("AWS_REGION")
	s3MediaBucketName := viper.GetString("AWS_S3_MEDIA_BUCKET_NAME")
	databaseClient := postgresClient.GetDatabaseInstance(log)
	db := databaseClient.GetDB()
	postgresMediaRepo := mediaPostgresAdapter.NewMediaCrudRepository(log, db)
	s3Client := aws.GetS3ClientInstance(log, region).GetClient()
	mediaPort := awsAdapater.NewS3Service(s3Client, log, region)
	mediaService := services.NewMediaService(log, mediaPort, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, s3MediaBucketName)
	mediaHandler := mediaRestAdapter.NewCrudPostHandler(log, mediaService)
	return mediaHandler
}

func BuildMediaCrudDomainGrpc(log *logrus.Logger) stylGrpc.CrudPostGrcp {
	region := viper.GetString("AWS_REGION")
	s3MediaBucketName := viper.GetString("AWS_S3_MEDIA_BUCKET_NAME")
	databaseClient := postgresClient.GetDatabaseInstance(log)
	db := databaseClient.GetDB()
	postgresMediaRepo := mediaPostgresAdapter.NewMediaCrudRepository(log, db)
	s3Client := aws.GetS3ClientInstance(log, region).GetClient()
	mediaPort := awsAdapater.NewS3Service(s3Client, log, region)
	mediaService := services.NewMediaService(log, mediaPort, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, s3MediaBucketName)
	grpcMediaService := stylGrpc.NewCrudPostGrcp(mediaService, log)
	return *grpcMediaService
}
