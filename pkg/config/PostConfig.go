package config

import (
	"github.com/sirupsen/logrus"
	"styl-monolith/internal/media/adapters/fileManager"
	stylGrpc "styl-monolith/internal/media/adapters/grcp"
	mediaPostgresAdapter "styl-monolith/internal/media/adapters/postgres"
	mediaRestAdapter "styl-monolith/internal/media/adapters/rest"
	"styl-monolith/internal/media/core/services"
	postgresClient "styl-monolith/pkg/postgres"
)

func BuildMediaCrudDomainHandler(log *logrus.Logger) mediaRestAdapter.CrudPostHandler {
	databaseClient := postgresClient.GetDatabaseInstance(log)
	db := databaseClient.GetDB()
	fileMngmt := fileManager.NewFileManager(log)
	postgresMediaRepo := mediaPostgresAdapter.NewMediaCrudRepository(log, db)
	mediaService := services.NewMediaService(log, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, fileMngmt)
	mediaHandler := mediaRestAdapter.NewCrudPostHandler(log, mediaService)
	return mediaHandler
}

func BuildMediaCrudDomainGrpc(log *logrus.Logger) stylGrpc.CrudPostGrcp {
	databaseClient := postgresClient.GetDatabaseInstance(log)
	db := databaseClient.GetDB()
	fileMngmt := fileManager.NewFileManager(log)
	postgresMediaRepo := mediaPostgresAdapter.NewMediaCrudRepository(log, db)
	mediaService := services.NewMediaService(log, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, fileMngmt)
	grpcMediaService := stylGrpc.NewCrudPostGrcp(mediaService, log)
	return *grpcMediaService
}
