package services

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"mime/multipart"
	"strconv"
	"styl-monolith/internal/media/core/domain"
	"styl-monolith/internal/media/core/ports"
	"styl-monolith/pkg/errorhandler"
	"styl-monolith/pkg/logger"
	"styl-monolith/pkg/utils"
)

type MediaService interface {
	// Media Operations
	UploadImage(imageData domain.PostImage, imageFile *multipart.FileHeader, userEmail string, postId uint) (domain.PostImage, error)
	DeleteImage(imageKey string) error
	GetImage(imageKey string) (domain.PostImage, error)
	ListImages(postId uint, page, size int) ([]domain.PostImage, int, error)
	SetMainImage(postId, imageId uint) error

	// Report Operations
	CreateReport(report domain.PostReport) (domain.PostReport, error)
	DeleteReport(reportId uint) error
	GetReport(reportId uint) (domain.PostReport, error)
	ListReports(page, size int) ([]domain.PostReport, int, error)
	ListReportsByPost(postId uint, page, size int) ([]domain.PostReport, int, error)
	ListReportsByUser(userId uint, page, size int) ([]domain.PostReport, int, error)
	GetPendingReports(page, size int) ([]domain.PostReport, error)
	ConfirmReport(reportId uint) error
	RejectReport(reportId uint) error

	// Post Operations
	CreatePost(post domain.Post) (domain.Post, error)
	DeletePost(postId uint) error
	HardDeletePost(postId uint) error
	ListPosts(page, size int) ([]domain.Post, int, error)
	GetPost(postId uint) (domain.Post, error)
	UpdatePost(postId uint, updatedPost domain.Post) (domain.Post, error)
	ListPostsByUser(userId uint, page, size int) ([]domain.Post, int, error)

	// Comment Operations
	CreateComment(comment domain.PostComment) (domain.PostComment, error)
	DeleteComment(commentId uint) error
	ListComments(postId uint, page, size int) ([]domain.PostComment, int, error)
	GetComment(commentId uint) (domain.PostComment, error)
	UpdateComment(commentId uint, updatedComment domain.PostComment) (domain.PostComment, error)

	// Like Operations
	CreateLike(like domain.PostLike) (domain.PostLike, error)
	DeleteLike(likeId uint) error
	ListLikes(page, size int) ([]domain.PostLike, int, error)
	GetLike(likeId uint) (domain.PostLike, error)
	ListPostLikes(postId uint, page, size int) ([]domain.PostLike, int, error)
}

type MediaServiceStruct struct {
	logger      *logrus.Logger
	mediaPort   ports.MediaPort
	reportPort  ports.ReportPort
	commentPort ports.CommentPort
	postPort    ports.PostPort
	likePort    ports.LikePort
	fileManager ports.FileManager
	s3Client    *s3.Client
}

// NewMediaService creates a new instance of MediaService.
func NewMediaService(
	log *logrus.Logger,
	mediaPort ports.MediaPort,
	reportPort ports.ReportPort,
	commentPort ports.CommentPort,
	postPort ports.PostPort,
	likePort ports.LikePort,
	fileManager ports.FileManager,
	s3Client *s3.Client,
) MediaService {
	return &MediaServiceStruct{
		logger:      log,
		mediaPort:   mediaPort,
		reportPort:  reportPort,
		commentPort: commentPort,
		postPort:    postPort,
		likePort:    likePort,
		fileManager: fileManager,
		s3Client:    s3Client,
	}
}

func (s *MediaServiceStruct) UploadImage(imageData domain.PostImage, imageFile *multipart.FileHeader, userEmail string, postId uint) (domain.PostImage, error) {
	s.logger.Info("Uploading image to S3")

	// Generate email hash
	emailHash := utils.GenerateHashFromString(userEmail)

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

	// Generate the S3 object key
	objectKey := emailHash + "/" + strconv.FormatUint(uint64(postId), 10) + "/" + utils.GenerateHashFromString(imageData.Filename)

	// Retrieve the S3 bucket name and region from environment variables
	s3BucketName := viper.GetString("AWS_S3_MEDIA_BUCKET_NAME")
	s3Region := viper.GetString("AWS_REGION")
	if s3BucketName == "" {
		s.logger.Error("Missing environment variable: AWS_S3_MEDIA_BUCKET_NAME")
		return domain.PostImage{}, errorhandler.NewDomainError(
			errorhandler.ErrMissingEnvVariable,
			"Missing environment variable: AWS_S3_MEDIA_BUCKET_NAME",
			nil,
		)
	}

	// Upload the file to S3
	_, err = s.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s3BucketName),
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

	// Generate the public S3 URL
	imageData.S3Url = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s3BucketName, s3Region, objectKey)
	imageData.S3ObjectKey = objectKey
	imageData.PostId = postId
	// Save image metadata to the database
	postImage, err := s.mediaPort.SaveMetadataOfImageOfS3(imageData)

	if err != nil {
		s.logger.Errorf("Failed to save metadata for S3 image: %v", err)
		return domain.PostImage{}, errorhandler.NewDomainError(
			errorhandler.ErrMetadataSaveFailed,
			errorhandler.GetErrorMessage(errorhandler.ErrMetadataSaveFailed),
			err,
		)
	}

	return postImage, nil
}

func (s *MediaServiceStruct) DeleteImage(imageKey string) error {
	s.logger.Infof("Deleting image from S3 with key: %s", imageKey)
	return s.mediaPort.DeleteImageFromS3(imageKey)
}

func (s *MediaServiceStruct) GetImage(imageKey string) (domain.PostImage, error) {
	s.logger.Infof("Getting image with key: %s", imageKey)
	return s.mediaPort.GetImageFromS3(imageKey)
}

func (s *MediaServiceStruct) ListImages(postId uint, page, size int) ([]domain.PostImage, int, error) {
	s.logger.Infof("Listing images for post ID: %d", postId)
	return s.mediaPort.ListImagesFromS3(postId, page, size)
}

func (s *MediaServiceStruct) SetMainImage(postId, imageId uint) error {
	s.logger.Infof("Setting image ID: %d as main image for post ID: %d", imageId, postId)
	return s.mediaPort.SetPostMainImage(postId, imageId)
}

// Report Operations
func (s *MediaServiceStruct) CreateReport(report domain.PostReport) (domain.PostReport, error) {
	s.logger.Info("Creating a post report")
	return s.reportPort.CreateReport(report)
}

func (s *MediaServiceStruct) DeleteReport(reportId uint) error {
	s.logger.Infof("Deleting report ID: %d", reportId)
	return s.reportPort.DeleteReport(reportId)
}

func (s *MediaServiceStruct) GetReport(reportId uint) (domain.PostReport, error) {
	s.logger.Infof("Getting report ID: %d", reportId)
	return s.reportPort.GetReport(reportId)
}

func (s *MediaServiceStruct) ListReports(page, size int) ([]domain.PostReport, int, error) {
	s.logger.Info("Listing all reports")
	return s.reportPort.ListReports(page, size)
}

func (s *MediaServiceStruct) ListReportsByPost(postId uint, page, size int) ([]domain.PostReport, int, error) {
	s.logger.Infof("Listing reports for post ID: %d", postId)
	return s.reportPort.ListReportsByPost(postId, page, size)
}

func (s *MediaServiceStruct) ListReportsByUser(userId uint, page, size int) ([]domain.PostReport, int, error) {
	s.logger.Infof("Listing reports by user ID: %d", userId)
	return s.reportPort.ListReportsByUser(userId, page, size)
}

func (s *MediaServiceStruct) GetPendingReports(page, size int) ([]domain.PostReport, error) {
	s.logger.Info("Listing pending reports")
	return s.reportPort.GetPendingReports(page, size)
}

func (s *MediaServiceStruct) ConfirmReport(reportId uint) error {
	s.logger.Infof("Confirming report ID: %d", reportId)
	return s.reportPort.ConfirmReport(reportId)
}

func (s *MediaServiceStruct) RejectReport(reportId uint) error {
	s.logger.Infof("Rejecting report ID: %d", reportId)
	return s.reportPort.RejectReport(reportId)
}

// Post Operations
func (s *MediaServiceStruct) CreatePost(post domain.Post) (domain.Post, error) {
	s.logger.Info("Creating a new post")
	return s.postPort.CreatePost(post)
}

func (s *MediaServiceStruct) DeletePost(postId uint) error {
	s.logger.Infof("Deleting post ID: %d", postId)
	return s.postPort.DeletePost(postId)
}

func (s *MediaServiceStruct) HardDeletePost(postId uint) error {
	s.logger.Infof("Deleting post ID: %d", postId)
	return s.postPort.DeletePost(postId)
}

func (s *MediaServiceStruct) ListPosts(page, size int) ([]domain.Post, int, error) {
	s.logger.Info("Listing paginated posts")
	posts, totalPages, err := s.postPort.ListPosts(page, size)
	if err != nil {
		return []domain.Post{}, 0, err
	}
	return posts, totalPages, nil
}

func (s *MediaServiceStruct) GetPost(postId uint) (domain.Post, error) {
	s.logger.Infof("Getting post ID: %d", postId)
	return s.postPort.GetPost(postId)
}

func (s *MediaServiceStruct) UpdatePost(postId uint, updatedPost domain.Post) (domain.Post, error) {
	s.logger.Infof("Updating post ID: %d", postId)
	return s.postPort.UpdatePost(postId, updatedPost)
}

func (s *MediaServiceStruct) ListPostsByUser(userId uint, page, size int) ([]domain.Post, int, error) {
	s.logger.Infof("Listing posts for user ID: %d", userId)
	return s.postPort.ListPostsByUser(userId, page, size)
}

// Comment Operations
func (s *MediaServiceStruct) CreateComment(comment domain.PostComment) (domain.PostComment, error) {
	s.logger.Info(logger.CreatingNewComment)
	post, err := s.GetPost(comment.PostId)
	if err != nil {
		s.logger.Errorf(logger.PostNotFound, comment.PostId)
		return domain.PostComment{}, err
	}
	comment.PostId = post.Id
	return s.commentPort.CreateComment(comment)
}

func (s *MediaServiceStruct) DeleteComment(commentId uint) error {
	s.logger.Infof("Deleting comment ID: %d", commentId)
	return s.commentPort.DeleteComment(commentId)
}

func (s *MediaServiceStruct) ListComments(postId uint, page, size int) ([]domain.PostComment, int, error) {
	s.logger.Infof("Listing comments for post ID: %d", postId)
	return s.commentPort.ListComments(postId, page, size)
}

func (s *MediaServiceStruct) GetComment(commentId uint) (domain.PostComment, error) {
	s.logger.Infof("Getting comment ID: %d", commentId)
	return s.commentPort.GetCommentById(commentId)
}

func (s *MediaServiceStruct) UpdateComment(commentId uint, updatedComment domain.PostComment) (domain.PostComment, error) {
	s.logger.Infof("Updating comment ID: %d", commentId)
	return s.commentPort.UpdateComment(commentId, updatedComment)
}

// Like Operations
func (s *MediaServiceStruct) CreateLike(like domain.PostLike) (domain.PostLike, error) {
	s.logger.Infof("Creating like for post ID: %d by user ID: %d", like.PostId, like.UserId)
	return s.likePort.CreateLike(like.UserId, like.PostId)
}

func (s *MediaServiceStruct) DeleteLike(likeId uint) error {
	s.logger.Infof("Deleting like ID: %d", likeId)
	return s.likePort.DeleteLike(likeId)
}

func (s *MediaServiceStruct) ListLikes(page, size int) ([]domain.PostLike, int, error) {
	s.logger.Info("Listing all likes")
	return s.likePort.ListLikes(page, size)
}

func (s *MediaServiceStruct) GetLike(likeId uint) (domain.PostLike, error) {
	s.logger.Infof("Getting like ID: %d", likeId)
	return s.likePort.GetLike(likeId)
}

func (s *MediaServiceStruct) ListPostLikes(postId uint, page, size int) ([]domain.PostLike, int, error) {
	s.logger.Infof("Listing likes for post ID: %d", postId)
	return s.likePort.GetPostLikes(postId, page, size)
}
