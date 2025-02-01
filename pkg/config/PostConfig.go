package config

import (
	"github.com/sirupsen/logrus"
	stylGrpc "styl-monolith/internal/media/adapters/grcp"
	mediaPostgresAdapter "styl-monolith/internal/media/adapters/postgres"
	mediaRestAdapter "styl-monolith/internal/media/adapters/rest"
	"styl-monolith/internal/media/core/services"
	postgresClient "styl-monolith/pkg/postgres"
)

func BuildMediaCrudDomainHandler(log *logrus.Logger) mediaRestAdapter.CrudPostHandler {
	databaseClient := postgresClient.GetDatabaseInstance(log)
	db := databaseClient.GetDB()
	postgresMediaRepo := mediaPostgresAdapter.NewMediaCrudRepository(log, db)
	mediaService := services.NewMediaService(log, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo)
	mediaHandler := mediaRestAdapter.NewCrudPostHandler(log, mediaService)
	return mediaHandler
}

func BuildMediaCrudDomainGrpc(log *logrus.Logger) stylGrpc.CrudPostGrcp {
	databaseClient := postgresClient.GetDatabaseInstance(log)
	db := databaseClient.GetDB()
	postgresMediaRepo := mediaPostgresAdapter.NewMediaCrudRepository(log, db)
	mediaService := services.NewMediaService(log, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo, postgresMediaRepo)
	grpcMediaService := stylGrpc.NewCrudPostGrcp(mediaService, log)
	return *grpcMediaService
}
