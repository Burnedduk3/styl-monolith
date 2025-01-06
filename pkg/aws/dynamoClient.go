package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/sirupsen/logrus"
	"styl-monolith/pkg/logger"
	"sync"
)

type DynamoClient struct {
	client *dynamodb.Client
}

var (
	dynamoInstance *DynamoClient
	dynamoOnce     sync.Once
)

// GetDynamoClientInstance creates or returns the single instance of DynamoClient
func GetDynamoClientInstance(logs *logrus.Logger, region string) *DynamoClient {
	logs.Debug(logger.OpeningDynamoConnection)
	dynamoOnce.Do(func() {
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
		)
		if err != nil {
			logs.Fatalf(fmt.Sprintf("failed to load AWS configuration: %v", err))
		}

		/* Create DynamoDB client */
		client := dynamodb.NewFromConfig(cfg)
		logs.Debug("DynamoDB client initialized successfully")

		// Initialize the singleton instance
		dynamoInstance = &DynamoClient{client: client}
	})
	return dynamoInstance
}

// GetClient returns the DynamoDB client instance
func (d *DynamoClient) GetClient() *dynamodb.Client {
	return d.client
}
