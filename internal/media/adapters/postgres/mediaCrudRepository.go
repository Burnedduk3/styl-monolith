package rest

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"styl-monolith/internal/media/adapters/postgres/models"
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

// SaveMetadataOfImageOfS3 adds a new image to the database and returns its S3 URL.
func (r *MediaCrudRepository) SaveMetadataOfImageOfS3(image domain.PostImage) (domain.PostImage, error) {
	postgresPostImage := models.NewPostgresPostImageFromDomainPostImage(image)

	result := r.conn.Create(&postgresPostImage)
	if result.Error != nil {
		err := errorhandler.NewDomainError(
			errorhandler.ErrImageDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrImageDatabaseUnableToCompleteOperation),
			result.Error,
		)
		return domain.PostImage{}, err
	}
	return postgresPostImage.ToPostImageDomain(), nil
}

// DeleteImageFromS3 deletes an image record from the database by its S3 URL.
func (r *MediaCrudRepository) DeleteImageFromS3(imageKey string) error {
	result := r.conn.Where("s3_url = ?", imageKey).Delete(&models.PostImage{})
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
	var postgresPostImage models.PostImage
	result := r.conn.Where("s3_url = ?", imageKey).First(&postgresPostImage)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.PostImage{}, errorhandler.NewDomainError(
				errorhandler.ErrImageNotFound,
				fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrImageNotFound), imageKey),
				nil,
			)
		}
		return domain.PostImage{}, errorhandler.NewDomainError(
			errorhandler.ErrImageDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrImageDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	return postgresPostImage.ToPostImageDomain(), nil
}

// ListImagesFromS3 lists images associated with a given post, paginated.
func (r *MediaCrudRepository) ListImagesFromS3(postId uint, page, size int) ([]domain.PostImage, int, error) {
	var postgresPostImages []models.PostImage
	var totalCount int64

	countResult := r.conn.Model(&models.Post{}).Count(&totalCount)

	if countResult.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			countResult.Error,
		)
	}

	offset := (page - 1) * size
	result := r.conn.Where("post_id = ?", postId).Limit(size).Offset(offset).Find(&postgresPostImages)

	if result.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrImageDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrImageDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	var postImages []domain.PostImage
	for _, postgresPostImage := range postgresPostImages {
		postImages = append(postImages, postgresPostImage.ToPostImageDomain())
	}
	// Calculate total number of pages
	totalPages := int((totalCount + int64(size) - 1) / int64(size))
	return postImages, totalPages, nil
}

// SetPostMainImage sets a specific image as the main image for a post.
func (r *MediaCrudRepository) SetPostMainImage(postId, imageId uint) error {
	// Reset the "is_main" flag for all images of the given post
	reset := r.conn.Model(&models.PostImage{}).Where("post_id = ?", postId).Update("is_main", false)
	if reset.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrImageDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrImageDatabaseUnableToCompleteOperation),
			reset.Error,
		)
	}

	// Set the specified image as the main image for the post
	result := r.conn.Model(&models.PostImage{}).Where("id = ? AND post_id = ?", imageId, postId).Update("is_main", true)
	if result.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrImageDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrImageDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	// Check if no rows were affected, indicating the image might not exist
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
	postgresComment := models.NewPostgresPostCommentFromDomainPostComment(comment)
	result := r.conn.Create(&postgresComment)
	if result.Error != nil {
		return domain.PostComment{}, errorhandler.NewDomainError(
			errorhandler.ErrCommentDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrCommentDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}
	return postgresComment.ToPostCommentDomain(), nil
}

// DeleteComment deletes a comment from the database by its ID.
func (r *MediaCrudRepository) DeleteComment(commentId uint) error {
	var postgresComment models.PostComment
	result := r.conn.First(&postgresComment, commentId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errorhandler.NewDomainError(
				errorhandler.ErrCommentNotFound,
				fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrCommentNotFound), commentId),
				nil,
			)
		}
		return errorhandler.NewDomainError(
			errorhandler.ErrCommentDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrCommentDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	deleteResult := r.conn.Delete(&postgresComment)
	if deleteResult.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrCommentDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrCommentDatabaseUnableToCompleteOperation),
			deleteResult.Error,
		)
	}

	return nil
}

// ListComments retrieves a paginated list of comments for a specific post.
func (r *MediaCrudRepository) ListComments(postId uint, page, size int) ([]domain.PostComment, int, error) {
	var postgresComments []models.PostComment
	var totalCount int64
	offset := (page - 1) * size

	countResult := r.conn.Model(&models.Post{}).Count(&totalCount)

	if countResult.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			countResult.Error,
		)
	}

	result := r.conn.Where("post_id = ?", postId).Limit(size).Offset(offset).Find(&postgresComments)
	if result.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrCommentDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrCommentDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	var comments []domain.PostComment
	for _, postgresComment := range postgresComments {
		comments = append(comments, postgresComment.ToPostCommentDomain())
	}
	// Calculate total number of pages
	totalPages := int((totalCount + int64(size) - 1) / int64(size))

	return comments, totalPages, nil
}

// GetCommentById retrieves a single comment by its ID.
func (r *MediaCrudRepository) GetCommentById(commentId uint) (domain.PostComment, error) {
	var postgresComment models.PostComment
	result := r.conn.First(&postgresComment, commentId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.PostComment{}, errorhandler.NewDomainError(
				errorhandler.ErrCommentNotFound,
				fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrCommentNotFound), commentId),
				nil,
			)
		}
		return domain.PostComment{}, errorhandler.NewDomainError(
			errorhandler.ErrCommentDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrCommentDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	return postgresComment.ToPostCommentDomain(), nil
}

// UpdateComment updates an existing comment in the database.
func (r *MediaCrudRepository) UpdateComment(commentId uint, updatedComment domain.PostComment) (domain.PostComment, error) {
	var postgresComment models.PostComment
	result := r.conn.First(&postgresComment, commentId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.PostComment{}, errorhandler.NewDomainError(
				errorhandler.ErrCommentNotFound,
				fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrCommentNotFound), commentId),
				nil,
			)
		}
		return domain.PostComment{}, errorhandler.NewDomainError(
			errorhandler.ErrCommentDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrCommentDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	postgresUpdatedComment := models.NewPostgresPostCommentFromDomainPostComment(updatedComment)
	postgresComment.Comment = postgresUpdatedComment.Comment
	postgresComment.IsDeleted = postgresUpdatedComment.IsDeleted
	postgresComment.IsReported = postgresUpdatedComment.IsReported

	updateResult := r.conn.Save(&postgresComment)
	if updateResult.Error != nil {
		return domain.PostComment{}, errorhandler.NewDomainError(
			errorhandler.ErrCommentDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrCommentDatabaseUnableToCompleteOperation),
			updateResult.Error,
		)
	}

	return postgresComment.ToPostCommentDomain(), nil
}

// CreatePost adds a new post to the database and returns the created post.
func (r *MediaCrudRepository) CreatePost(post domain.Post) (domain.Post, error) {
	postgresPost := models.NewPostgresPostFromDomainPost(post)

	result := r.conn.Create(&postgresPost)
	if result.Error != nil {
		return domain.Post{}, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	return postgresPost.ToPostDomain(), nil
}

// DeletePost deletes a post from the database by its ID.
func (r *MediaCrudRepository) DeletePost(postId uint) error {
	result := r.conn.Delete(&models.Post{}, postId)
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

// HardDeletePost deletes a post and all its related data from the database by its ID.
func (r *MediaCrudRepository) HardDeletePost(postId uint) error {
	// Begin a transaction
	tx := r.conn.Begin()
	if tx.Error != nil {
		r.log.Errorf("Failed to begin transaction: %v", tx.Error)
		return errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			tx.Error,
		)
	}

	// Delete related PostImages
	if err := tx.Where("post_id = ?", postId).Delete(&models.PostImage{}).Error; err != nil {
		r.log.Errorf("Failed to delete related images for post ID %d: %v", postId, err)
		tx.Rollback()
		return errorhandler.NewDomainError(
			errorhandler.ErrImageDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrImageDatabaseUnableToCompleteOperation),
			err,
		)
	}

	// Delete related PostComments
	if err := tx.Where("post_id = ?", postId).Delete(&models.PostComment{}).Error; err != nil {
		r.log.Errorf("Failed to delete related comments for post ID %d: %v", postId, err)
		tx.Rollback()
		return errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			"Failed to delete related comments for the post.",
			err,
		)
	}

	// Delete related PostLikes
	if err := tx.Where("post_id = ?", postId).Delete(&models.PostLike{}).Error; err != nil {
		r.log.Errorf("Failed to delete related likes for post ID %d: %v", postId, err)
		tx.Rollback()
		return errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			"Failed to delete related likes for the post.",
			err,
		)
	}

	// Delete related PostReports
	if err := tx.Where("post_id = ?", postId).Delete(&models.PostReport{}).Error; err != nil {
		r.log.Errorf("Failed to delete related reports for post ID %d: %v", postId, err)
		tx.Rollback()
		return errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			"Failed to delete related reports for the post.",
			err,
		)
	}

	// Delete related PostSaves
	if err := tx.Where("post_id = ?", postId).Delete(&models.PostSaved{}).Error; err != nil {
		r.log.Errorf("Failed to delete related saves for post ID %d: %v", postId, err)
		tx.Rollback()
		return errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			"Failed to delete related saves for the post.",
			err,
		)
	}

	// Delete the post itself
	if err := tx.Delete(&models.Post{}, postId).Error; err != nil {
		r.log.Errorf("Failed to delete post ID %d: %v", postId, err)
		tx.Rollback()
		return errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			err,
		)
	}

	// Check if the post actually existed
	if tx.RowsAffected == 0 {
		r.log.Warnf("Post ID %d not found.", postId)
		tx.Rollback()
		return errorhandler.NewDomainError(
			errorhandler.ErrPostNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrPostNotFound), postId),
			nil,
		)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		r.log.Errorf("Failed to commit transaction for deleting post ID %d: %v", postId, err)
		return errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			err,
		)
	}

	r.log.Infof("Successfully deleted post ID %d and all related data.", postId)
	return nil
}

// ListPosts retrieves a paginated list of posts along with the total number of pages,
// including PostImages relation and counts for likes and comments for each post.
func (r *MediaCrudRepository) ListPosts(page, size int) ([]domain.Post, int, error) {
	var postgresPosts []models.Post
	var totalCount int64
	offset := (page - 1) * size

	// Count total records
	countResult := r.conn.Model(&models.Post{}).Count(&totalCount)
	if countResult.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			countResult.Error,
		)
	}

	// Retrieve current page's records with PostImages relation
	queryResult := r.conn.Preload("PostImages").Limit(size).Offset(offset).Find(&postgresPosts)
	if queryResult.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			queryResult.Error,
		)
	}

	// Convert Postgres models to domain models and fetch counts
	var posts []domain.Post
	for _, postgresPost := range postgresPosts {
		// Count comments for the current post
		var commentCount int64
		r.conn.Model(&models.PostComment{}).Where("post_id = ?", postgresPost.ID).Count(&commentCount)

		// Count likes for the current post
		var likeCount int64
		r.conn.Model(&models.PostLike{}).Where("post_id = ?", postgresPost.ID).Count(&likeCount)

		// Convert the Post model to the domain model and set counts
		domainPost := postgresPost.ToPostDomain()
		domainPost.CommentCount = int(commentCount)
		domainPost.LikeCount = int(likeCount)

		posts = append(posts, domainPost)
	}

	// Calculate total number of pages
	totalPages := int((totalCount + int64(size) - 1) / int64(size))

	return posts, totalPages, nil
}

// GetPost retrieves a single post by its ID, includes its PostImages relation,
// and gets the total count of comments and likes.
func (r *MediaCrudRepository) GetPost(postId uint) (domain.Post, error) {
	var postgresPost models.Post
	var commentCount int64
	var likeCount int64

	result := r.conn.Preload("PostImages").First(&postgresPost, postId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.Post{}, errorhandler.NewDomainError(
				errorhandler.ErrPostNotFound,
				fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrPostNotFound), postId),
				nil,
			)
		}
		return domain.Post{}, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	r.conn.Model(&models.PostComment{}).Where("post_id = ?", postId).Count(&commentCount)

	r.conn.Model(&models.PostLike{}).Where("post_id = ?", postId).Count(&likeCount)

	domainPost := postgresPost.ToPostDomain()
	domainPost.CommentCount = int(commentCount)
	domainPost.LikeCount = int(likeCount)

	return domainPost, nil
}

// UpdatePost updates an existing post in the database.
func (r *MediaCrudRepository) UpdatePost(postId uint, updatedPost domain.Post) (domain.Post, error) {
	// Retrieve the existing Post model
	var postgresPost models.Post
	result := r.conn.First(&postgresPost, postId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.Post{}, errorhandler.NewDomainError(
				errorhandler.ErrPostNotFound,
				fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrPostNotFound), postId),
				nil,
			)
		}
		return domain.Post{}, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	postgresPost.Caption = updatedPost.Caption       // Update caption
	postgresPost.IsPublic = updatedPost.IsPublic     // Update visibility (public/private)
	postgresPost.IsApproved = updatedPost.IsApproved // Update approval status
	postgresPost.IsSpam = updatedPost.IsSpam         // Update spam marker
	postgresPost.ViewCount = updatedPost.ViewCount   // Update view count, if manual adjustment is allowed

	// Save the updated Post model back to the database
	updateResult := r.conn.Save(&postgresPost)

	if updateResult.Error != nil {
		return domain.Post{}, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			updateResult.Error,
		)
	}
	newPost := postgresPost.ToPostDomain()
	newPost.PostImages = updatedPost.PostImages

	return newPost, nil
}

// IncreasePostViewCount increments the view count for a post.
func (r *MediaCrudRepository) IncreasePostViewCount(postId uint) error {
	result := r.conn.Model(&models.Post{}).Where("id = ?", postId).Update("view_count", gorm.Expr("view_count + ?", 1))
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
func (r *MediaCrudRepository) ListPostsByUser(userId uint, page, size int) ([]domain.Post, int, error) {
	var postgresPosts []models.Post
	var totalCount int64

	// Count the total number of posts by userId
	countResult := r.conn.Model(&models.Post{}).Where("user_id = ?", userId).Count(&totalCount)
	if countResult.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			countResult.Error,
		)
	}

	// Paginate the query
	offset := (page - 1) * size
	result := r.conn.Where("user_id = ?", userId).Limit(size).Offset(offset).Find(&postgresPosts)
	if result.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	// Convert Postgres models to domain models
	var posts []domain.Post
	for _, postgresPost := range postgresPosts {
		posts = append(posts, postgresPost.ToPostDomain())
	}

	// Calculate total number of pages
	totalPages := int((totalCount + int64(size) - 1) / int64(size))

	return posts, totalPages, nil
}

// GetPostsByTag retrieves all posts associated with a specific tag, paginated.
func (r *MediaCrudRepository) GetPostsByTag(tagId uint, page, size int) ([]domain.Post, error) {
	var postgresPosts []models.Post
	offset := (page - 1) * size
	result := r.conn.
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Where("post_tags.tag_id = ?", tagId).
		Limit(size).Offset(offset).Find(&postgresPosts)

	if result.Error != nil {
		return nil, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	// Convert Postgres models to domain models
	var posts []domain.Post
	for _, postgresPost := range postgresPosts {
		posts = append(posts, postgresPost.ToPostDomain())
	}

	return posts, nil
}

// ReportPost adds a report to a post.
func (r *MediaCrudRepository) ReportPost(postId uint, reason string, userId uint) error {
	postgresReport := models.NewPostgresPostReportFromDomainPostReport(domain.PostReport{
		PostId:   postId,
		Reason:   reason,
		UserId:   userId,
		Resolved: false,
	})

	result := r.conn.Create(&postgresReport)
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
	result := r.conn.Model(&models.Post{}).Where("id = ?", postId).Update("is_spam", true)
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
	result := r.conn.Model(&models.Post{}).Where("id = ?", postId).Update("is_approved", true)
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
func (r *MediaCrudRepository) ListFlaggedPosts(page, size int) ([]domain.Post, int, error) {
	var postgresPosts []models.Post
	var totalCount int64

	// Count the total number of flagged or spam posts
	countResult := r.conn.Model(&models.Post{}).Where("is_spam = ? OR is_flagged = ?", true, true).Count(&totalCount)
	if countResult.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			countResult.Error,
		)
	}

	// Paginate the query
	offset := (page - 1) * size
	result := r.conn.Where("is_spam = ? OR is_flagged = ?", true, true).Limit(size).Offset(offset).Find(&postgresPosts)
	if result.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrPostDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrPostDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	// Convert Postgres models to domain models
	var posts []domain.Post
	for _, postgresPost := range postgresPosts {
		posts = append(posts, postgresPost.ToPostDomain())
	}

	// Calculate total number of pages
	totalPages := int((totalCount + int64(size) - 1) / int64(size))

	return posts, totalPages, nil
}

// TogglePostPrivacy updates the privacy of a post.
func (r *MediaCrudRepository) TogglePostPrivacy(postId uint, isPublic bool) error {
	result := r.conn.Model(&models.Post{}).Where("id = ?", postId).Update("is_public", isPublic)
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
	postgresLike := models.NewPostgresPostLikeFromDomainPostLike(domain.PostLike{
		UserId: userId,
		PostId: postId,
	})

	result := r.conn.Create(&postgresLike)
	if result.Error != nil {
		return domain.PostLike{}, errorhandler.NewDomainError(
			errorhandler.ErrLikeDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrLikeDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	// Convert back to a domain object
	return postgresLike.ToPostLikeDomain(), nil
}

// DeleteLike removes a like from the database by its ID.
func (r *MediaCrudRepository) DeleteLike(likeId uint) error {
	result := r.conn.Delete(&models.PostLike{}, likeId)
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
func (r *MediaCrudRepository) ListLikes(page, size int) ([]domain.PostLike, int, error) {
	var postgresLikes []models.PostLike
	var totalCount int64

	// Count the total number of likes
	countResult := r.conn.Model(&models.PostLike{}).Count(&totalCount)
	if countResult.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrLikeDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrLikeDatabaseUnableToCompleteOperation),
			countResult.Error,
		)
	}

	// Paginate the query
	offset := (page - 1) * size
	result := r.conn.Limit(size).Offset(offset).Find(&postgresLikes)
	if result.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrLikeDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrLikeDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	// Convert Postgres models to domain models
	var likes []domain.PostLike
	for _, postgresLike := range postgresLikes {
		likes = append(likes, postgresLike.ToPostLikeDomain())
	}

	// Calculate total number of pages
	totalPages := int((totalCount + int64(size) - 1) / int64(size))

	return likes, totalPages, nil
}

// GetLike retrieves a single like by its ID.
func (r *MediaCrudRepository) GetLike(likeId uint) (domain.PostLike, error) {
	var postgresLike models.PostLike
	result := r.conn.First(&postgresLike, likeId)
	if result.Error != nil {
		return domain.PostLike{}, errorhandler.NewDomainError(
			errorhandler.ErrLikeNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrLikeNotFound), likeId),
			result.Error,
		)
	}

	// Convert to a domain object
	return postgresLike.ToPostLikeDomain(), nil
}

// GetPostLikes retrieves a list of likes for a specific post, paginated.
func (r *MediaCrudRepository) GetPostLikes(postId uint, page, size int) ([]domain.PostLike, int, error) {
	var postgresLikes []models.PostLike
	var totalCount int64

	// Count the total number of likes for the specific post
	countResult := r.conn.Model(&models.PostLike{}).Where("post_id = ?", postId).Count(&totalCount)
	if countResult.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrLikeDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrLikeDatabaseUnableToCompleteOperation),
			countResult.Error,
		)
	}

	// Paginate the query
	offset := (page - 1) * size
	result := r.conn.Where("post_id = ?", postId).Limit(size).Offset(offset).Find(&postgresLikes)
	if result.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrLikeDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrLikeDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	// Convert Postgres likes to domain likes
	var likes []domain.PostLike
	for _, postgresLike := range postgresLikes {
		likes = append(likes, postgresLike.ToPostLikeDomain())
	}

	// Calculate total number of pages
	totalPages := int((totalCount + int64(size) - 1) / int64(size))

	return likes, totalPages, nil
}

// ResolveReport marks a specific report as resolved.
func (r *MediaCrudRepository) ResolveReport(reportId uint) error {
	result := r.conn.Model(&models.PostReport{}).Where("id = ?", reportId).Update("resolved", true)
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
func (r *MediaCrudRepository) ListReportsForPost(postId uint, page, size int) ([]domain.PostReport, int, error) {
	var postgresReports []models.PostReport
	var totalCount int64

	// Count the total number of reports for the given postId
	countResult := r.conn.Model(&models.PostReport{}).Where("post_id = ?", postId).Count(&totalCount)
	if countResult.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			countResult.Error,
		)
	}

	// Paginate the query
	offset := (page - 1) * size
	result := r.conn.Where("post_id = ?", postId).Limit(size).Offset(offset).Find(&postgresReports)
	if result.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	// Convert Postgres reports to domain reports
	var reports []domain.PostReport
	for _, postgresReport := range postgresReports {
		reports = append(reports, postgresReport.ToPostReportDomain())
	}

	// Calculate total number of pages
	totalPages := int((totalCount + int64(size) - 1) / int64(size))

	return reports, totalPages, nil
}

// CreateReport creates a new report for a specific post.
func (r *MediaCrudRepository) CreateReport(report domain.PostReport) (domain.PostReport, error) {
	postgresReport := models.NewPostgresPostReportFromDomainPostReport(report)

	result := r.conn.Create(&postgresReport)
	if result.Error != nil {
		return domain.PostReport{}, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	// Convert back to a domain object
	return postgresReport.ToPostReportDomain(), nil
}

// DeleteReport deletes a report by its ID.
func (r *MediaCrudRepository) DeleteReport(reportId uint) error {
	result := r.conn.Delete(&models.PostReport{}, reportId)
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
func (r *MediaCrudRepository) ListReports(page, size int) ([]domain.PostReport, int, error) {
	var postgresReports []models.PostReport
	var totalCount int64

	// Count the total number of reports
	countResult := r.conn.Model(&models.PostReport{}).Count(&totalCount)
	if countResult.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			countResult.Error,
		)
	}

	// Paginate the query
	offset := (page - 1) * size
	result := r.conn.Limit(size).Offset(offset).Find(&postgresReports)
	if result.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	// Convert Postgres reports to domain reports
	var reports []domain.PostReport
	for _, postgresReport := range postgresReports {
		reports = append(reports, postgresReport.ToPostReportDomain())
	}

	// Calculate total number of pages
	totalPages := int((totalCount + int64(size) - 1) / int64(size))

	return reports, totalPages, nil
}

// GetReport retrieves a single report by its ID.
func (r *MediaCrudRepository) GetReport(reportId uint) (domain.PostReport, error) {
	var postgresReport models.PostReport
	result := r.conn.First(&postgresReport, reportId)
	if result.Error != nil {
		return domain.PostReport{}, errorhandler.NewDomainError(
			errorhandler.ErrReportNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrReportNotFound), reportId),
			result.Error,
		)
	}

	// Convert to domain object
	return postgresReport.ToPostReportDomain(), nil
}

// UpdateReport updates an existing report by its ID.
func (r *MediaCrudRepository) UpdateReport(reportId uint, updatedReport domain.PostReport) (domain.PostReport, error) {
	postgresReport := models.NewPostgresPostReportFromDomainPostReport(updatedReport)

	result := r.conn.Model(&models.PostReport{}).Where("id = ?", reportId).Updates(postgresReport)
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

	return postgresReport.ToPostReportDomain(), nil
}

// ListReportsByPost retrieves reports associated with a specific post, paginated.
func (r *MediaCrudRepository) ListReportsByPost(postId uint, page, size int) ([]domain.PostReport, int, error) {
	var postgresReports []models.PostReport
	var totalCount int64

	// Count the total number of reports for the specific post
	countResult := r.conn.Model(&models.PostReport{}).Where("post_id = ?", postId).Count(&totalCount)
	if countResult.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			countResult.Error,
		)
	}

	// Paginate the query
	offset := (page - 1) * size
	result := r.conn.Where("post_id = ?", postId).Limit(size).Offset(offset).Find(&postgresReports)
	if result.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	// Convert Postgres reports to domain reports
	var reports []domain.PostReport
	for _, postgresReport := range postgresReports {
		reports = append(reports, postgresReport.ToPostReportDomain())
	}

	// Calculate total number of pages
	totalPages := int((totalCount + int64(size) - 1) / int64(size))

	return reports, totalPages, nil
}

// ListReportsByUser retrieves reports created by a specific user, paginated.
func (r *MediaCrudRepository) ListReportsByUser(userId uint, page, size int) ([]domain.PostReport, int, error) {
	var postgresReports []models.PostReport
	var totalCount int64

	// Count the total number of reports created by the specific user
	countResult := r.conn.Model(&models.PostReport{}).Where("user_id = ?", userId).Count(&totalCount)
	if countResult.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			countResult.Error,
		)
	}

	// Paginate the query
	offset := (page - 1) * size
	result := r.conn.Where("user_id = ?", userId).Limit(size).Offset(offset).Find(&postgresReports)
	if result.Error != nil {
		return nil, 0, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	// Convert Postgres reports to domain reports
	var reports []domain.PostReport
	for _, postgresReport := range postgresReports {
		reports = append(reports, postgresReport.ToPostReportDomain())
	}

	// Calculate total number of pages
	totalPages := int((totalCount + int64(size) - 1) / int64(size))

	return reports, totalPages, nil
}

// GetPendingReports retrieves unresolved reports, paginated.
func (r *MediaCrudRepository) GetPendingReports(page, size int) ([]domain.PostReport, error) {
	var postgresReports []models.PostReport
	offset := (page - 1) * size
	result := r.conn.Where("resolved = ?", false).Limit(size).Offset(offset).Find(&postgresReports)
	if result.Error != nil {
		return nil, errorhandler.NewDomainError(
			errorhandler.ErrReportDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrReportDatabaseUnableToCompleteOperation),
			result.Error,
		)
	}

	// Convert Postgres reports to domain reports
	var reports []domain.PostReport
	for _, postgresReport := range postgresReports {
		reports = append(reports, postgresReport.ToPostReportDomain())
	}

	return reports, nil
}

// ConfirmReport marks a report as resolved/accepted.
func (r *MediaCrudRepository) ConfirmReport(reportId uint) error {
	result := r.conn.Model(&models.PostReport{}).Where("id = ?", reportId).Update("resolved", true)
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
