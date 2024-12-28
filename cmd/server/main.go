package main

import (
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
	viper.SetDefault("LOG_LEVEL", "debug") // Default log level

	// Configure logger log level from environment variable
	logLevel := viper.GetString("LOG_LEVEL")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic(fmt.Sprintf("Invalid log level: %s", logLevel))
	}
	logrus.SetLevel(level)

	// Verify required configurations
	requiredVars := []string{"DB_HOST", "DB_USER", "DB_NAME", "DB_PASSWORD", "DB_SSLMODE"}
	for _, key := range requiredVars {
		if !viper.IsSet(key) {
			panic(fmt.Sprintf("Missing required configuration %s", key))
		}
	}
}

func RegisterUsersDomainRoutes(e *echo.Echo, logger *logrus.Logger) {
	crudUserDomainHandler := config.BuildUserCrudDomainHandler()
	loginUserDomainHandler := config.BuildLoginUserDomainHandler()
	utils.RegisterRoutesAutomatically(e, crudUserDomainHandler, "api/v1", logger)
	utils.RegisterRoutesAutomatically(e, loginUserDomainHandler, "api/v1", logger)
}

// main initializes and starts both an HTTP server using Echo and a gRPC server on their respective ports.
// It handles server errorhandler gracefully and logs critical messages.
func main() {
	logger := log.GetLogger(logrus.GetLevel())
	// Spin up Echo HTTP server on port 1323
	e := echo.New()
	go func() {
		RegisterUsersDomainRoutes(e, logger)
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
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
