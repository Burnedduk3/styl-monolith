package services

import (
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"strconv"
	"styl-monolith/internal/media/core/domain"
	"styl-monolith/internal/media/core/ports"
	"styl-monolith/pkg/errorhandler"
	"styl-monolith/pkg/logger"
	"styl-monolith/pkg/utils"
)

// MediaService defines an interface for handling media operations such as uploads, reports, posts, comments, and likes.
type MediaService interface {
	// Media Operations
	UploadImage(imagesData []domain.PostImage, imageFiles []*multipart.FileHeader, userEmail string, postId uint) ([]domain.PostImage, error)
	DeleteImage(imageKey string) error

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
	SetMainImage(postId, imageId uint) error
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

// MediaServiceStruct defines the structure for managing media services and integrations with various ports and resources.
type MediaServiceStruct struct {
	logger            *logrus.Logger
	mediaPort         ports.MediaPort
	reportPort        ports.ReportPort
	commentPort       ports.CommentPort
	postPort          ports.PostPort
	likePort          ports.LikePort
	s3MediaBucketName string
}

// NewMediaService initializes and returns a MediaService instance with the provided dependencies and configurations.
func NewMediaService(
	log *logrus.Logger,
	mediaPort ports.MediaPort,
	reportPort ports.ReportPort,
	commentPort ports.CommentPort,
	postPort ports.PostPort,
	likePort ports.LikePort,
	s3MediaBucketName string,
) MediaService {
	return &MediaServiceStruct{
		logger:            log,
		mediaPort:         mediaPort,
		reportPort:        reportPort,
		commentPort:       commentPort,
		postPort:          postPort,
		likePort:          likePort,
		s3MediaBucketName: s3MediaBucketName,
	}
}

// UploadImage uploads multiple images to S3 concurrently, saves their metadata, and associates them with a given post.
// It returns a list of successfully uploaded images or an error in case of failure.
func (s *MediaServiceStruct) UploadImage(imagesData []domain.PostImage, imageFiles []*multipart.FileHeader, userEmail string, postId uint) ([]domain.PostImage, error) {
	s.logger.Info("Uploading images to S3 concurrently")

	emailHash := utils.GenerateHashFromString(userEmail)
	var savedImages []domain.PostImage
	savedImageChan := make(chan domain.PostImage, len(imagesData))
	errorChan := make(chan error, len(imagesData))

	for index, imageMetadata := range imagesData {
		go func(index int, imageMetadata domain.PostImage) {
			objectKey := emailHash + "/" + strconv.FormatUint(uint64(postId), 10) + "/" + utils.GenerateHashFromString(imageMetadata.Filename) + "." + imageMetadata.Format

			imageData, err := s.mediaPort.SaveImageToS3(imageMetadata, imageFiles[index], s.s3MediaBucketName, objectKey, postId)
			if err != nil {
				s.logger.Errorf("Failed to save metadata for S3 image: %v", err)
				errorChan <- err
				return
			}

			postImage, err := s.postPort.SaveMetadataOfImageOfS3(imageData)
			if err != nil {
				s.logger.Errorf("Failed to save metadata for S3 image: %v", err)
				errorChan <- errorhandler.NewDomainError(
					errorhandler.ErrMetadataSaveFailed,
					errorhandler.GetErrorMessage(errorhandler.ErrMetadataSaveFailed),
					err,
				)
				return
			}

			savedImageChan <- postImage
		}(index, imageMetadata)
	}

	for i := 0; i < len(imagesData); i++ {
		select {
		case savedImage := <-savedImageChan:
			savedImages = append(savedImages, savedImage)
		case err := <-errorChan:
			return nil, err
		}
	}

	return savedImages, nil
}

// DeleteImage removes an image from the S3 storage using the provided image key and bucket name. Returns an error if it fails.
func (s *MediaServiceStruct) DeleteImage(imageKey string) error {
	s.logger.Infof("Deleting image from S3 with key: %s", imageKey)
	return s.mediaPort.DeleteImageFromS3(s.s3MediaBucketName, imageKey)
}

// SetMainImage sets the provided image ID as the main image for the specified post ID and logs the operation.
func (s *MediaServiceStruct) SetMainImage(postId, imageId uint) error {
	s.logger.Infof("Setting image ID: %d as main image for post ID: %d", imageId, postId)
	return s.postPort.SetPostMainImage(postId, imageId)
}

// CreateReport creates a new report for a specific post and returns the created report or an error.
func (s *MediaServiceStruct) CreateReport(report domain.PostReport) (domain.PostReport, error) {
	s.logger.Info("Creating a post report")
	return s.reportPort.CreateReport(report)
}

// DeleteReport removes a report identified by its unique ID from the system and returns an error if the operation fails.
func (s *MediaServiceStruct) DeleteReport(reportId uint) error {
	s.logger.Infof("Deleting report ID: %d", reportId)
	return s.reportPort.DeleteReport(reportId)
}

// GetReport retrieves a post report by its unique report ID and returns the report along with any potential error.
func (s *MediaServiceStruct) GetReport(reportId uint) (domain.PostReport, error) {
	s.logger.Infof("Getting report ID: %d", reportId)
	return s.reportPort.GetReport(reportId)
}

// ListReports retrieves a paginated list of all post reports with the specified page number and size.
func (s *MediaServiceStruct) ListReports(page, size int) ([]domain.PostReport, int, error) {
	s.logger.Info("Listing all reports")
	return s.reportPort.ListReports(page, size)
}

// ListReportsByPost fetches a paginated list of reports for a specific post by its ID, along with the total number of pages.
func (s *MediaServiceStruct) ListReportsByPost(postId uint, page, size int) ([]domain.PostReport, int, error) {
	s.logger.Infof("Listing reports for post ID: %d", postId)
	return s.reportPort.ListReportsByPost(postId, page, size)
}

// ListReportsByUser retrieves a paginated list of post reports generated by a specific user, based on user ID, page, and size parameters.
func (s *MediaServiceStruct) ListReportsByUser(userId uint, page, size int) ([]domain.PostReport, int, error) {
	s.logger.Infof("Listing reports by user ID: %d", userId)
	return s.reportPort.ListReportsByUser(userId, page, size)
}

// GetPendingReports retrieves a paginated list of post reports that are pending resolution.
func (s *MediaServiceStruct) GetPendingReports(page, size int) ([]domain.PostReport, error) {
	s.logger.Info("Listing pending reports")
	return s.reportPort.GetPendingReports(page, size)
}

// ConfirmReport confirms a specific report identified by its ID. It logs the action and delegates the task to the reportPort.
func (s *MediaServiceStruct) ConfirmReport(reportId uint) error {
	s.logger.Infof("Confirming report ID: %d", reportId)
	return s.reportPort.ConfirmReport(reportId)
}

// RejectReport rejects a specific post report identified by its report ID and logs the operation.
func (s *MediaServiceStruct) RejectReport(reportId uint) error {
	s.logger.Infof("Rejecting report ID: %d", reportId)
	return s.reportPort.RejectReport(reportId)
}

// CreatePost creates a new post by forwarding the provided post object to the postPort for processing.
func (s *MediaServiceStruct) CreatePost(post domain.Post) (domain.Post, error) {
	s.logger.Info("Creating a new post")
	return s.postPort.CreatePost(post)
}

// DeletePost deletes a post with the given postId by invoking the corresponding method on the postPort. It logs the operation.
func (s *MediaServiceStruct) DeletePost(postId uint) error {
	s.logger.Infof("Deleting post ID: %d", postId)
	return s.postPort.DeletePost(postId)
}

// HardDeletePost permanently deletes a post identified by its ID by delegating the operation to the post port.
func (s *MediaServiceStruct) HardDeletePost(postId uint) error {
	s.logger.Infof("Deleting post ID: %d", postId)
	return s.postPort.DeletePost(postId)
}

// ListPosts retrieves a paginated list of posts, returning the list of posts, total pages, and an error if any occurs.
func (s *MediaServiceStruct) ListPosts(page, size int) ([]domain.Post, int, error) {
	s.logger.Info("Listing paginated posts")
	posts, totalPages, err := s.postPort.ListPosts(page, size)
	if err != nil {
		return []domain.Post{}, 0, err
	}
	return posts, totalPages, nil
}

// GetPost retrieves a post by its unique identifier.
// It returns the post object and an error if any issues occur during retrieval.
func (s *MediaServiceStruct) GetPost(postId uint) (domain.Post, error) {
	s.logger.Infof("Getting post ID: %d", postId)
	return s.postPort.GetPost(postId)
}

// UpdatePost updates an existing post with the given post ID using the provided updated post data.
// Returns the updated post and any error encountered during the update process.
func (s *MediaServiceStruct) UpdatePost(postId uint, updatedPost domain.Post) (domain.Post, error) {
	s.logger.Infof("Updating post ID: %d", postId)
	return s.postPort.UpdatePost(postId, updatedPost)
}

// ListPostsByUser retrieves a paginated list of posts created by a specific user.
// It returns a slice of posts, the total number of pages, and an error if any operation fails.
func (s *MediaServiceStruct) ListPostsByUser(userId uint, page, size int) ([]domain.Post, int, error) {
	s.logger.Infof("Listing posts for user ID: %d", userId)
	return s.postPort.ListPostsByUser(userId, page, size)
}

// CreateComment creates a new comment for a specific post and returns the created comment or an error if it fails.
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

// DeleteComment removes a comment identified by its ID and returns an error if the deletion is unsuccessful.
func (s *MediaServiceStruct) DeleteComment(commentId uint) error {
	s.logger.Infof("Deleting comment ID: %d", commentId)
	return s.commentPort.DeleteComment(commentId)
}

// ListComments retrieves a paginated list of comments for a specified post ID, with the total page count as a return value.
func (s *MediaServiceStruct) ListComments(postId uint, page, size int) ([]domain.PostComment, int, error) {
	s.logger.Infof("Listing comments for post ID: %d", postId)
	return s.commentPort.ListComments(postId, page, size)
}

// GetComment retrieves a comment by its ID and returns the corresponding PostComment or an error if retrieval fails.
func (s *MediaServiceStruct) GetComment(commentId uint) (domain.PostComment, error) {
	s.logger.Infof("Getting comment ID: %d", commentId)
	return s.commentPort.GetCommentById(commentId)
}

// UpdateComment updates an existing comment by its ID with the provided updated comment data and returns the updated comment.
func (s *MediaServiceStruct) UpdateComment(commentId uint, updatedComment domain.PostComment) (domain.PostComment, error) {
	s.logger.Infof("Updating comment ID: %d", commentId)
	return s.commentPort.UpdateComment(commentId, updatedComment)
}

// CreateLike creates a like for a given post and user, and logs the process. It interacts with the likePort to persist the like.
func (s *MediaServiceStruct) CreateLike(like domain.PostLike) (domain.PostLike, error) {
	s.logger.Infof("Creating like for post ID: %d by user ID: %d", like.PostId, like.UserId)
	return s.likePort.CreateLike(like.UserId, like.PostId)
}

// DeleteLike removes a like identified by the provided like ID and logs the operation.
func (s *MediaServiceStruct) DeleteLike(likeId uint) error {
	s.logger.Infof("Deleting like ID: %d", likeId)
	return s.likePort.DeleteLike(likeId)
}

// ListLikes retrieves a paginated list of post likes, returning the likes, total pages, and any encountered error.
func (s *MediaServiceStruct) ListLikes(page, size int) ([]domain.PostLike, int, error) {
	s.logger.Info("Listing all likes")
	return s.likePort.ListLikes(page, size)
}

// GetLike retrieves a like with the specified likeId from the data source and returns it along with any encountered error.
func (s *MediaServiceStruct) GetLike(likeId uint) (domain.PostLike, error) {
	s.logger.Infof("Getting like ID: %d", likeId)
	return s.likePort.GetLike(likeId)
}

// ListPostLikes retrieves a paginated list of likes for a specific post based on its ID.
// Returns the likes, the total number of pages, and an error if any occurred during the operation.
func (s *MediaServiceStruct) ListPostLikes(postId uint, page, size int) ([]domain.PostLike, int, error) {
	s.logger.Infof("Listing likes for post ID: %d", postId)
	return s.likePort.GetPostLikes(postId, page, size)
}
