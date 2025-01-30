package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
	"styl-monolith/pkg/logger"
	"sync"
)

type S3Client struct {
	client *s3.Client
}

var (
	s3Instance *S3Client
	s3Once     sync.Once
)

// GetS3ClientInstance creates or returns the single instance of S3Client
func GetS3ClientInstance(logs *logrus.Logger, region string) *S3Client {
	logs.Debug(logger.OpeningS3Connection)

	// Ensure the S3 client is created only once (singleton)
	s3Once.Do(func() {
		// Load AWS configuration
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region), // Sets the region dynamically
		)
		if err != nil {
			logs.Fatalf(fmt.Sprintf("failed to load AWS configuration: %s", err))
		}

		// Create S3 client
		client := s3.NewFromConfig(cfg)
		logs.Debug("S3 client initialized successfully")

		// Initialize the singleton instance
		s3Instance = &S3Client{client: client}
	})
	return s3Instance
}

// GetClient returns the S3 client instance
func (s *S3Client) GetClient() *s3.Client {
	return s.client
}
