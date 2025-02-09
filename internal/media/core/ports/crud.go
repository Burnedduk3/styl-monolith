package ports

import (
	"mime/multipart"
	"styl-monolith/internal/media/core/domain"
)

type MediaPort interface {
	SaveImageToS3(imageMetadata domain.PostImage, imageFile *multipart.FileHeader, bucketName, objectKey string, postId uint) (domain.PostImage, error)
	DeleteImageFromS3(bucketName, objectKey string) error
}

type CommentPort interface {
	CreateComment(comment domain.PostComment) (domain.PostComment, error)
	DeleteComment(commentId uint) error
	ListComments(postId uint, page, size int) ([]domain.PostComment, int, error)
	GetCommentById(commentId uint) (domain.PostComment, error)
	UpdateComment(commentId uint, comment domain.PostComment) (domain.PostComment, error)
}

type PostPort interface {
	CreatePost(post domain.Post) (domain.Post, error)
	SaveMetadataOfImageOfS3(image domain.PostImage) (domain.PostImage, error)
	SetPostMainImage(postId, imageId uint) error
	DeletePost(postId uint) error
	HardDeletePost(postId uint) error
	ListPosts(page, size int) ([]domain.Post, int, error)
	GetPost(postId uint) (domain.Post, error)
	UpdatePost(postId uint, updatedData domain.Post) (domain.Post, error)
	IncreasePostViewCount(postId uint) error
	ListPostsByUser(userId uint, page, size int) ([]domain.Post, int, error)
	GetPostsByTag(tagId uint, page, size int) ([]domain.Post, error)
	ReportPost(postId uint, reason string, userId uint) error
	MarkPostAsSpam(postId uint) error
	ApprovePost(postId uint) error
	ListFlaggedPosts(page, size int) ([]domain.Post, int, error)
	TogglePostPrivacy(postId uint, isPublic bool) error
}

type LikePort interface {
	CreateLike(userId, postId uint) (domain.PostLike, error)
	DeleteLike(likeId uint) error
	ListLikes(page, size int) ([]domain.PostLike, int, error)
	GetLike(likeId uint) (domain.PostLike, error)
	GetPostLikes(postId uint, page, size int) ([]domain.PostLike, int, error) // Likes for a specific post
}

type ReportPort interface {
	CreateReport(report domain.PostReport) (domain.PostReport, error)
	DeleteReport(reportId uint) error
	ListReports(page, size int) ([]domain.PostReport, int, error)
	GetReport(reportId uint) (domain.PostReport, error)
	UpdateReport(reportId uint, updatedReport domain.PostReport) (domain.PostReport, error)
	ListReportsByPost(postId uint, page, size int) ([]domain.PostReport, int, error)
	ListReportsByUser(userId uint, page, size int) ([]domain.PostReport, int, error)
	GetPendingReports(page, size int) ([]domain.PostReport, error)
	ConfirmReport(reportId uint) error
	RejectReport(reportId uint) error
}
