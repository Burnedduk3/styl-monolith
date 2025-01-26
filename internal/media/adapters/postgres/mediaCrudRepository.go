package rest

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"styl-monolith/internal/media/core/domain"
	"styl-monolith/pkg/errorhandler"
)

type MediaCrudRepository struct {
	log  *logrus.Logger
	conn *gorm.DB
}

// NewMediaCrudRepository creates a new instance of MediaCrudRepository.
func NewMediaCrudRepository(log *logrus.Logger, conn *gorm.DB) *MediaCrudRepository {
	return &MediaCrudRepository{log: log, conn: conn}
}

// UploadImageToS3 adds a new image to the database and returns its S3 URL.
func (r *MediaCrudRepository) UploadImageToS3(image domain.PostImage) (string, error) {
	result := r.conn.Create(&image)
	if result.Error != nil {
		err := errorhandler.NewDomainError(
			errorhandler.ErrImageDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrImageDatabaseUnableToCompleteOperation),
			result.Error,
		)
		return "", err
	}
	return image.S3Url, nil
}

// DeleteImageFromS3 deletes an image record from the database by its S3 URL.
func (r *MediaCrudRepository) DeleteImageFromS3(imageKey string) error {
	result := r.conn.Where("s3_url = ?", imageKey).Delete(&domain.PostImage{})
	if result.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrImageDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrImageDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	if result.RowsAffected == 0 {
		return errorhandler.NewDomainError(
			errorhandler.ErrImageNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrImageNotFound), imageKey),
			nil,
		)
	}
	return nil
}

// GetImageFromS3 retrieves an image from the database by its S3 URL.
func (r *MediaCrudRepository) GetImageFromS3(imageKey string) (domain.PostImage, error) {
	var image domain.PostImage
	result := r.conn.Where("s3_url = ?", imageKey).First(&image)
	if result.Error != nil {
		return domain.PostImage{}, errorhandler.NewDomainError(
			errorhandler.ErrImageNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrImageNotFound), imageKey),
			result.Error,
		)
	}
	return image, nil
}

// ListImagesFromS3 lists images associated with a given post, paginated.
func (r *MediaCrudRepository) ListImagesFromS3(postId uint, page, size int) ([]domain.PostImage, error) {
	var images []domain.PostImage
	offset := (page - 1) * size
	result := r.conn.Where("post_id = ?", postId).Limit(size).Offset(offset).Find(&images)
	if result.Error != nil {
		return nil, errorhandler.NewDomainError(
			errorhandler.ErrImageDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrImageDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return images, nil
}

// SetPostMainImage sets a specific image as the main image for a post.
func (r *MediaCrudRepository) SetPostMainImage(postId, imageId uint) error {
	// Reset main image for all images of the post
	reset := r.conn.Model(&domain.PostImage{}).Where("post_id = ?", postId).Update("is_main", false)
	if reset.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrImageDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrImageDatabaseUnableToCompleteOperation),
			reset.Error,
		)
	}

	// Set the specified image as the main image
	result := r.conn.Model(&domain.PostImage{}).Where("id = ? AND post_id = ?", imageId, postId).Update("is_main", true)
	if result.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrImageDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrImageDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	if result.RowsAffected == 0 {
		return errorhandler.NewDomainError(
			errorhandler.ErrImageNotFound,
			fmt.Sprintf("Image with ID %d for Post ID %d not found", imageId, postId),
			nil,
		)
	}
	return nil
}

// CreateComment adds a new comment to the database.
func (r *MediaCrudRepository) CreateComment(comment domain.PostComment) (domain.PostComment, error) {
	result := r.conn.Create(&comment)
	if result.Error != nil {
		err := errorhandler.NewDomainError(
			errorhandler.ErrCommentDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrCommentDatabaseUnableToCompleteOperation),
			result.Error,
		)
		return domain.PostComment{}, err
	}
	return comment, nil
}

// DeleteComment deletes a comment from the database by its ID.
func (r *MediaCrudRepository) DeleteComment(commentId uint) error {
	result := r.conn.Delete(&domain.PostComment{}, commentId)
	if result.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrCommentDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrCommentDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	if result.RowsAffected == 0 {
		return errorhandler.NewDomainError(
			errorhandler.ErrCommentNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrCommentNotFound), commentId),
			nil,
		)
	}
	return nil
}

// ListComments retrieves a paginated list of comments for a specific post.
func (r *MediaCrudRepository) ListComments(postId uint, page, size int) ([]domain.PostComment, error) {
	var comments []domain.PostComment
	offset := (page - 1) * size
	result := r.conn.Where("post_id = ?", postId).Limit(size).Offset(offset).Find(&comments)
	if result.Error != nil {
		return nil, errorhandler.NewDomainError(
			errorhandler.ErrCommentDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrCommentDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return comments, nil
}

// GetCommentById retrieves a single comment by its ID.
func (r *MediaCrudRepository) GetCommentById(commentId uint) (domain.PostComment, error) {
	var comment domain.PostComment
	result := r.conn.First(&comment, commentId)
	if result.Error != nil {
		return domain.PostComment{}, errorhandler.NewDomainError(
			errorhandler.ErrCommentNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrCommentNotFound), commentId),
			result.Error,
		)
	}
	return comment, nil
}

// UpdateComment updates an existing comment in the database.
func (r *MediaCrudRepository) UpdateComment(commentId uint, updatedComment domain.PostComment) (domain.PostComment, error) {
	result := r.conn.Model(&domain.PostComment{}).Where("id = ?", commentId).Updates(updatedComment)
	if result.Error != nil {
		return domain.PostComment{}, errorhandler.NewDomainError(
			errorhandler.ErrCommentDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrCommentDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	if result.RowsAffected == 0 {
		return domain.PostComment{}, errorhandler.NewDomainError(
			errorhandler.ErrCommentNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrCommentNotFound), commentId),
			nil,
		)
	}
	return updatedComment, nil
}

func (r *MediaCrudRepository) CreatePost(post domain.Post) (domain.Post, error) {
	result := r.conn.Create(&post)
	if result.Error != nil {
		return domain.Post{}, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return post, nil
}

// DeletePost deletes a post from the database by its ID.
func (r *MediaCrudRepository) DeletePost(postId uint) error {
	result := r.conn.Delete(&domain.Post{}, postId)
	if result.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	if result.RowsAffected == 0 {
		return errorhandler.NewDomainError(
			errorhandler.ErrPostNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrPostNotFound), postId),
			nil,
		)
	}
	return nil
}

// ListPosts retrieves a paginated list of posts.
func (r *MediaCrudRepository) ListPosts(page, size int) ([]domain.Post, error) {
	var posts []domain.Post
	offset := (page - 1) * size
	result := r.conn.Limit(size).Offset(offset).Find(&posts)
	if result.Error != nil {
		return nil, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return posts, nil
}

// GetPost retrieves a single post by its ID.
func (r *MediaCrudRepository) GetPost(postId uint) (domain.Post, error) {
	var post domain.Post
	result := r.conn.First(&post, postId)
	if result.Error != nil {
		return domain.Post{}, errorhandler.NewDomainError(
			errorhandler.ErrPostNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrPostNotFound), postId),
			result.Error,
		)
	}
	return post, nil
}

// UpdatePost updates an existing post in the database.
func (r *MediaCrudRepository) UpdatePost(postId uint, updatedPost domain.Post) (domain.Post, error) {
	result := r.conn.Model(&domain.Post{}).Where("id = ?", postId).Updates(updatedPost)
	if result.Error != nil {
		return domain.Post{}, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	if result.RowsAffected == 0 {
		return domain.Post{}, errorhandler.NewDomainError(
			errorhandler.ErrPostNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrPostNotFound), postId),
			nil,
		)
	}
	return updatedPost, nil
}

// IncreasePostViewCount increments the view count for a post.
func (r *MediaCrudRepository) IncreasePostViewCount(postId uint) error {
	result := r.conn.Model(&domain.Post{}).Where("id = ?", postId).Update("view_count", gorm.Expr("view_count + ?", 1))
	if result.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	if result.RowsAffected == 0 {
		return errorhandler.NewDomainError(
			errorhandler.ErrPostNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrPostNotFound), postId),
			nil,
		)
	}
	return nil
}

// ListPostsByUser retrieves all posts created by a specific user, paginated.
func (r *MediaCrudRepository) ListPostsByUser(userId uint, page, size int) ([]domain.Post, error) {
	var posts []domain.Post
	offset := (page - 1) * size
	result := r.conn.Where("user_id = ?", userId).Limit(size).Offset(offset).Find(&posts)
	if result.Error != nil {
		return nil, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return posts, nil
}

// GetPostsByTag retrieves all posts associated with a specific tag, paginated.
func (r *MediaCrudRepository) GetPostsByTag(tagId uint, page, size int) ([]domain.Post, error) {
	var posts []domain.Post
	offset := (page - 1) * size
	result := r.conn.
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Where("post_tags.tag_id = ?", tagId).
		Limit(size).Offset(offset).Find(&posts)
	if result.Error != nil {
		return nil, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return posts, nil
}

// ReportPost adds a report to a post.
func (r *MediaCrudRepository) ReportPost(postId uint, reason string, userId uint) error {
	report := domain.PostReport{PostId: postId, Reason: reason, UserId: userId, Resolved: false}
	result := r.conn.Create(&report)
	if result.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return nil
}

// MarkPostAsSpam sets a post's `is_spam` status to `true`.
func (r *MediaCrudRepository) MarkPostAsSpam(postId uint) error {
	result := r.conn.Model(&domain.Post{}).Where("id = ?", postId).Update("is_spam", true)
	if result.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	if result.RowsAffected == 0 {
		return errorhandler.NewDomainError(
			errorhandler.ErrPostNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrPostNotFound), postId),
			nil,
		)
	}
	return nil
}

// ApprovePost sets a post's `is_approved` status to `true`.
func (r *MediaCrudRepository) ApprovePost(postId uint) error {
	result := r.conn.Model(&domain.Post{}).Where("id = ?", postId).Update("is_approved", true)
	if result.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	if result.RowsAffected == 0 {
		return errorhandler.NewDomainError(
			errorhandler.ErrPostNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrPostNotFound), postId),
			nil,
		)
	}
	return nil
}

// ListFlaggedPosts retrieves all posts marked as spam or flagged, paginated.
func (r *MediaCrudRepository) ListFlaggedPosts(page, size int) ([]domain.Post, error) {
	var posts []domain.Post
	offset := (page - 1) * size
	result := r.conn.Where("is_spam = ? OR is_flagged = ?", true, true).Limit(size).Offset(offset).Find(&posts)
	if result.Error != nil {
		return nil, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return posts, nil
}

// TogglePostPrivacy updates the privacy of a post.
func (r *MediaCrudRepository) TogglePostPrivacy(postId uint, isPublic bool) error {
	result := r.conn.Model(&domain.Post{}).Where("id = ?", postId).Update("is_public", isPublic)
	if result.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	if result.RowsAffected == 0 {
		return errorhandler.NewDomainError(
			errorhandler.ErrPostNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrPostNotFound), postId),
			nil,
		)
	}
	return nil
}

// CreateLike adds a new like to a post by a specific user.
func (r *MediaCrudRepository) CreateLike(userId, postId uint) (domain.PostLike, error) {
	like := domain.PostLike{UserId: userId, PostId: postId}
	result := r.conn.Create(&like)
	if result.Error != nil {
		return domain.PostLike{}, errorhandler.NewDomainError(
			errorhandler.ErrLikeDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrLikeDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return like, nil
}

// DeleteLike removes a like from the database by its ID.
func (r *MediaCrudRepository) DeleteLike(likeId uint) error {
	result := r.conn.Delete(&domain.PostLike{}, likeId)
	if result.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrLikeDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrLikeDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	if result.RowsAffected == 0 {
		return errorhandler.NewDomainError(
			errorhandler.ErrLikeNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrLikeNotFound), likeId),
			nil,
		)
	}
	return nil
}

// ListLikes retrieves a paginated list of all likes in the system.
func (r *MediaCrudRepository) ListLikes(page, size int) ([]domain.PostLike, error) {
	var likes []domain.PostLike
	offset := (page - 1) * size
	result := r.conn.Limit(size).Offset(offset).Find(&likes)
	if result.Error != nil {
		return nil, errorhandler.NewDomainError(
			errorhandler.ErrLikeDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrLikeDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return likes, nil
}

// GetLike retrieves a single like by its ID.
func (r *MediaCrudRepository) GetLike(likeId uint) (domain.PostLike, error) {
	var like domain.PostLike
	result := r.conn.First(&like, likeId)
	if result.Error != nil {
		return domain.PostLike{}, errorhandler.NewDomainError(
			errorhandler.ErrLikeNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrLikeNotFound), likeId),
			result.Error,
		)
	}
	return like, nil
}

// GetPostLikes retrieves a list of likes for a specific post, paginated.
func (r *MediaCrudRepository) GetPostLikes(postId uint, page, size int) ([]domain.PostLike, error) {
	var likes []domain.PostLike
	offset := (page - 1) * size
	result := r.conn.Where("post_id = ?", postId).Limit(size).Offset(offset).Find(&likes)
	if result.Error != nil {
		return nil, errorhandler.NewDomainError(
			errorhandler.ErrLikeDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrLikeDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return likes, nil
}

// ResolveReport marks a specific report as resolved.
func (r *MediaCrudRepository) ResolveReport(reportId uint) error {
	result := r.conn.Model(&domain.PostReport{}).Where("id = ?", reportId).Update("resolved", true)
	if result.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	if result.RowsAffected == 0 {
		return errorhandler.NewDomainError(
			errorhandler.ErrReportNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrReportNotFound), reportId),
			nil,
		)
	}
	return nil
}

// ListReportsForPost retrieves all reports for a specific post, paginated.
func (r *MediaCrudRepository) ListReportsForPost(postId uint, page, size int) ([]domain.PostReport, error) {
	var reports []domain.PostReport
	offset := (page - 1) * size
	result := r.conn.Where("post_id = ?", postId).Limit(size).Offset(offset).Find(&reports)
	if result.Error != nil {
		return nil, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return reports, nil
}

// CreateReport creates a new report for a specific post.
func (r *MediaCrudRepository) CreateReport(report domain.PostReport) (domain.PostReport, error) {
	result := r.conn.Create(&report)
	if result.Error != nil {
		return domain.PostReport{}, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return report, nil
}

// DeleteReport deletes a report by its ID.
func (r *MediaCrudRepository) DeleteReport(reportId uint) error {
	result := r.conn.Delete(&domain.PostReport{}, reportId)
	if result.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	if result.RowsAffected == 0 {
		return errorhandler.NewDomainError(
			errorhandler.ErrReportNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrReportNotFound), reportId),
			nil,
		)
	}
	return nil
}

// ListReports retrieves a paginated list of all reports.
func (r *MediaCrudRepository) ListReports(page, size int) ([]domain.PostReport, error) {
	var reports []domain.PostReport
	offset := (page - 1) * size
	result := r.conn.Limit(size).Offset(offset).Find(&reports)
	if result.Error != nil {
		return nil, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return reports, nil
}

// GetReport retrieves a single report by its ID.
func (r *MediaCrudRepository) GetReport(reportId uint) (domain.PostReport, error) {
	var report domain.PostReport
	result := r.conn.First(&report, reportId)
	if result.Error != nil {
		return domain.PostReport{}, errorhandler.NewDomainError(
			errorhandler.ErrReportNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrReportNotFound), reportId),
			result.Error,
		)
	}
	return report, nil
}

// UpdateReport updates an existing report by its ID.
func (r *MediaCrudRepository) UpdateReport(reportId uint, updatedReport domain.PostReport) (domain.PostReport, error) {
	result := r.conn.Model(&domain.PostReport{}).Where("id = ?", reportId).Updates(updatedReport)
	if result.Error != nil {
		return domain.PostReport{}, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	if result.RowsAffected == 0 {
		return domain.PostReport{}, errorhandler.NewDomainError(
			errorhandler.ErrReportNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrReportNotFound), reportId),
			nil,
		)
	}
	return updatedReport, nil
}

// ListReportsByPost retrieves reports associated with a specific post, paginated.
func (r *MediaCrudRepository) ListReportsByPost(postId uint, page, size int) ([]domain.PostReport, error) {
	var reports []domain.PostReport
	offset := (page - 1) * size
	result := r.conn.Where("post_id = ?", postId).Limit(size).Offset(offset).Find(&reports)
	if result.Error != nil {
		return nil, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return reports, nil
}

// ListReportsByUser retrieves reports created by a specific user, paginated.
func (r *MediaCrudRepository) ListReportsByUser(userId uint, page, size int) ([]domain.PostReport, error) {
	var reports []domain.PostReport
	offset := (page - 1) * size
	result := r.conn.Where("user_id = ?", userId).Limit(size).Offset(offset).Find(&reports)
	if result.Error != nil {
		return nil, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return reports, nil
}

// GetPendingReports retrieves unresolved reports, paginated.
func (r *MediaCrudRepository) GetPendingReports(page, size int) ([]domain.PostReport, error) {
	var reports []domain.PostReport
	offset := (page - 1) * size
	result := r.conn.Where("resolved = ?", false).Limit(size).Offset(offset).Find(&reports)
	if result.Error != nil {
		return nil, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return reports, nil
}

// ConfirmReport marks a report as resolved/accepted.
func (r *MediaCrudRepository) ConfirmReport(reportId uint) error {
	result := r.conn.Model(&domain.PostReport{}).Where("id = ?", reportId).Update("resolved", true)
	if result.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	if result.RowsAffected == 0 {
		return errorhandler.NewDomainError(
			errorhandler.ErrReportNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrReportNotFound), reportId),
			nil,
		)
	}
	return nil
}

// RejectReport deletes a report, marking it as rejected.
func (r *MediaCrudRepository) RejectReport(reportId uint) error {
	return r.DeleteReport(reportId)
}
