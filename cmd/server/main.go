package main

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"styl-monolith/pkg/config"
	log "styl-monolith/pkg/logger"
	"styl-monolith/pkg/middleware"
	"styl-monolith/pkg/utils"

	pb "styl-monolith/generated/proto/media"
)

// apiVersion specifies the current API version being used in the application, typically for routing or version control.
var apiVersion = "v1"

// initialPath defines the base API path using the current API version concatenated to "/api/".
var initialPath = "/api/" + apiVersion

// init initializes the application configuration using Viper and sets default values for required environment variables.
// It ensures that mandatory configuration variables are set; otherwise, it panics with a message.
func init() {
	// Initialize Viper to load configurations
	viper.AutomaticEnv() // Automatically read from environment variables
	viper.SetDefault("DB_HOST", "127.0.0.1")
	viper.SetDefault("DB_USER", "root")
	viper.SetDefault("DB_NAME", "test_db")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("LOG_LEVEL", "debug") // Default log level
	viper.SetDefault("AWS_REGION", "us-east-1")
	viper.SetDefault("AWS_COGNITO_USER_POOL_ID", "us-east-1_4fXnfhGa9")
	viper.SetDefault("AWS_COGNITO_USER_POOL_CLIENT_ID", "MISSING")
	viper.SetDefault("AWS_COGNITO_USER_POOL_CLIENT_SECRET", "MISSING")
	viper.SetDefault("API_BASE_URL", "http://127.0.0.1:1323")
	viper.SetDefault("AWS_DYNAMODB_TOKENS_TABLE", "staging-auth-tokens")
	viper.SetDefault("AWS_S3_MEDIA_BUCKET_NAME", "staging-post-images-scrzo9nn")
	viper.SetDefault("AWS_REGION", "us-east-1")

	// Verify required configurations
	requiredVars := []string{"DB_HOST", "DB_USER", "DB_NAME", "DB_PASSWORD", "DB_SSLMODE", "AWS_S3_MEDIA_BUCKET_NAME"}
	for _, key := range requiredVars {
		if !viper.IsSet(key) {
			panic(fmt.Sprintf("Missing required configuration %s", key))
		}
	}
}

// registerUsersDomainRoutes registers CRUD and login endpoints for the users domain, with appropriate middleware for routing.
func registerUsersDomainRoutes(e *echo.Echo, logger *logrus.Logger) {
	crudUserDomainHandler := config.BuildUserCrudDomainHandler(logger)
	loginUserDomainHandler := config.BuildLoginUserDomainHandler(logger)
	dynamodbRepo := config.CreateAwsAuthClient(logger)
	utils.RegisterRoutesAutomatically(e, crudUserDomainHandler, initialPath, logger, true, middleware.ValidateAccessTokenWithRepository(dynamodbRepo))
	utils.RegisterRoutesAutomatically(e, loginUserDomainHandler, initialPath, logger, false, nil)
}

// registerMediaDomainRoutes configures and registers routes for the media domain in the provided Echo instance.
func registerMediaDomainRoutes(e *echo.Echo, logger *logrus.Logger) {
	crudMediaDomainHandler := config.BuildMediaCrudDomainHandler(logger)
	dynamodbRepo := config.CreateAwsAuthClient(logger)
	utils.RegisterRoutesAutomatically(e, crudMediaDomainHandler, initialPath, logger, true, middleware.ValidateAccessTokenWithRepository(dynamodbRepo))
}

// registerMediaDomainRoutes configures and registers routes for the media domain in the provided Echo instance.
func registerMediaDomainGrpc(serv *grpc.Server, logger *logrus.Logger) *grpc.Server {
	grpcHandler := config.BuildMediaCrudDomainGrpc(logger)
	pb.RegisterCrudPostServiceServer(serv, &grpcHandler)
	return serv
}

// main initializes and starts both the Echo HTTP server and the gRPC server for handling various services and routes.
func main() {
	// Configure logger log level from environment variable
	logLevel := viper.GetString("LOG_LEVEL")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic(fmt.Sprintf("Invalid log level: %s", logLevel))
	}
	logger := log.GetLogger(level)
	// Spin up Echo HTTP server on port 1323
	e := echo.New()
	go func() {
		registerUsersDomainRoutes(e, logger)
		registerMediaDomainRoutes(e, logger)
		logger.Info("Starting Echo server on port 1323")
		if err := e.Start(":1323"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Failed to start Echo server: %v", err)
		}
	}()

	// Spin up gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatal("Failed to listen on port 50051: %v", err)
	}
	grpcServer := grpc.NewServer()
	grpcServer = registerMediaDomainGrpc(grpcServer, logger)
	logger.Info("Starting gRPC server on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("Failed to serve gRPC server: %v", err)
	}
}
