package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/sirupsen/logrus"
	"styl-monolith/pkg/logger"
	"sync"
)

type CognitoClient struct {
	client *cognitoidentityprovider.Client
}

var (
	cognitoInstance *CognitoClient
	cognitoOnce     sync.Once
)

// GetCognitoClientInstance creates or returns the single instance of CognitoClient
func GetCognitoClientInstance(logs *logrus.Logger, region string) *CognitoClient {
	logs.Debug(logger.OpeningDatabaseConnection)

	// Ensure the Cognito client is created only once (singleton)
	cognitoOnce.Do(func() {
		// Load AWS configuration
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region), // Sets the region dynamically
		)
		if err != nil {
			logs.Fatalf(fmt.Sprintf("failed to load AWS configuration: %s", err))
		}

		// Create Cognito Identity Provider client
		client := cognitoidentityprovider.NewFromConfig(cfg)
		logs.Debug("Cognito client initialized successfully")

		// Initialize the singleton instance
		cognitoInstance = &CognitoClient{client: client}
	})
	return cognitoInstance
}

// GetClient returns the Cognito Identity Provider client instance
func (c *CognitoClient) GetClient() *cognitoidentityprovider.Client {
	return c.client
}
