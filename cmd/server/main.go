package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"net/http"
	log "styl-monolith/pkg/logger"
)

// init configures the application by loading environment variables and setting up the logger.
func init() {
	// Initialize Viper to load configurations
	viper.AutomaticEnv()     // Automatically read from environment variables
	viper.SetEnvPrefix("DB") // Environment variable prefix "DB_"
	viper.SetDefault("HOST", "127.0.0.1")
	viper.SetDefault("USER", "root")
	viper.SetDefault("NAME", "test_db")
	viper.SetDefault("PASSWORD", "password")
	viper.SetDefault("DB_SSLMODE", "disable")

	// Verify required configurations
	requiredVars := []string{"HOST", "USER", "NAME", "PASSWORD"}
	for _, key := range requiredVars {
		if !viper.IsSet(key) {
			panic(fmt.Sprintf("Missing required configuration %s", key))
		}
	}
}

// main initializes and starts both an HTTP server using Echo and a gRPC server on their respective ports.
// It handles server errors gracefully and logs critical messages.
func main() {
	logger := log.GetLogger(logrus.InfoLevel)
	// Spin up Echo HTTP server on port 1323
	e := echo.New()
	go func() {
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
