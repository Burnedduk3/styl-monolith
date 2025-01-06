package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"styl-monolith/internal/users/adapters/postgres/models"
	log "styl-monolith/pkg/logger"
	"styl-monolith/pkg/postgres"
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

	// Verify required configurations
	requiredVars := []string{"DB_HOST", "DB_USER", "DB_NAME", "DB_PASSWORD", "DB_SSLMODE"}
	for _, key := range requiredVars {
		if !viper.IsSet(key) {
			panic(fmt.Sprintf("Missing required configuration %s", key))
		}
	}
}

func main() {
	// Configure logger log level from environment variable
	logLevel := viper.GetString("LOG_LEVEL")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic(fmt.Sprintf("Invalid log level: %s", logLevel))
	}
	logger := log.GetLogger(level)

	databaseClient := postgres.GetDatabaseInstance(logger)
	db := databaseClient.GetDB()
	// Perform database migration
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		logger.Fatal(err)
	}
	// Perform database migration
	err = db.AutoMigrate(&models.Role{})
	if err != nil {
		logger.Fatal(err)
	}

}
