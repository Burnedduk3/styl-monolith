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
	"styl-monolith/pkg/utils"
)

// init configures the application by loading environment variables and setting up the logger.
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
	viper.SetDefault("API_BASE_URL", "http://127.0.0.1:1323")

	// Verify required configurations
	requiredVars := []string{"DB_HOST", "DB_USER", "DB_NAME", "DB_PASSWORD", "DB_SSLMODE"}
	for _, key := range requiredVars {
		if !viper.IsSet(key) {
			panic(fmt.Sprintf("Missing required configuration %s", key))
		}
	}
}

func RegisterUsersDomainRoutes(e *echo.Echo, logger *logrus.Logger) {
	crudUserDomainHandler := config.BuildUserCrudDomainHandler(logger)
	loginUserDomainHandler := config.BuildLoginUserDomainHandler(logger)
	utils.RegisterRoutesAutomatically(e, crudUserDomainHandler, "api/v1", logger)
	utils.RegisterRoutesAutomatically(e, loginUserDomainHandler, "api/v1", logger)
}

// main initializes and starts both an HTTP server using Echo and a gRPC server on their respective ports.
// It handles server errorhandler gracefully and logs critical messages.
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
		RegisterUsersDomainRoutes(e, logger)
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
	logger.Info("Starting gRPC server on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("Failed to serve gRPC server: %v", err)
	}
}
