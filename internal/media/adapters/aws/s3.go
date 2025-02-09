package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"styl-monolith/internal/media/core/domain"
	"styl-monolith/pkg/errorhandler"
)

type S3Service struct {
	s3Client *s3.Client
	logger   *logrus.Logger
	region   string
}

// NewS3Service initializes the S3 service with the provided bucket name.
func NewS3Service(s3Client *s3.Client, log *logrus.Logger, region string) *S3Service {
	return &S3Service{
		logger:   log,
		s3Client: s3Client,
		region:   region,
	}
}

// SaveImageToS3 saves an image to S3 and returns the saved image with its key or an error.
func (s *S3Service) SaveImageToS3(imageMetadata domain.PostImage, imageFile *multipart.FileHeader, bucketName, objectKey string, postId uint) (domain.PostImage, error) {
	// Open the image file
	file, err := imageFile.Open()
	if err != nil {
		s.logger.Errorf("Failed to open image file: %v", err)
		return domain.PostImage{}, errorhandler.NewDomainError(
			errorhandler.ErrFileOpenFailed,
			errorhandler.GetErrorMessage(errorhandler.ErrFileOpenFailed),
			err,
		)
	}
	defer file.Close()

	// Upload the file to S3
	_, err = s.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		s.logger.Errorf("Failed to upload file to S3: %v", err)
		return domain.PostImage{}, errorhandler.NewDomainError(
			errorhandler.ErrS3UploadFailed,
			errorhandler.GetErrorMessage(errorhandler.ErrS3UploadFailed),
			err,
		)
	}

	imageMetadata.S3Url = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, s.region, objectKey)
	imageMetadata.S3ObjectKey = objectKey
	imageMetadata.PostId = postId

	return imageMetadata, nil
}

// DeleteImageFromS3 deletes an image from S3 using the given bucket name and object key.
func (s *S3Service) DeleteImageFromS3(bucketName, objectKey string) error {
	_, err := s.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		s.logger.Errorf("Failed to delete image from S3: %v", err)
		return errorhandler.NewDomainError(
			errorhandler.ErrS3DeleteFailed,
			errorhandler.GetErrorMessage(errorhandler.ErrS3DeleteFailed),
			err,
		)
	}

	return nil
}
