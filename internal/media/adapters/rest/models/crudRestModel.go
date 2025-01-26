package models

import (
	"errors"
	"styl-monolith/internal/media/core/domain"
	"time"
)

// --------------------------------------------------
// POST MODELS
// --------------------------------------------------

// PostPayload represents the payload for creating/updating a post
type PostPayload struct {
	UserId    uint     `json:"userId"`
	Caption   string   `json:"caption"`
	IsPublic  bool     `json:"isPublic"`
	TagIds    []uint   `json:"tagIds"`
	ImageUrls []string `json:"imageUrls"`
}

func (p *PostPayload) Validate() error {
	if p.UserId == 0 {
		return errors.New("userId is required")
	}
	if p.Caption == "" {
		return errors.New("caption is required")
	}
	return nil
}

func (p *PostPayload) ToDomainPost() domain.Post {
	var images []domain.PostImage
	for _, url := range p.ImageUrls {
		images = append(images, domain.PostImage{
			S3Url: url,
		})
	}

	return domain.Post{
		UserId:     p.UserId,
		Caption:    p.Caption,
		IsPublic:   p.IsPublic,
		TagIds:     p.TagIds,
		PostImages: images,
	}
}

// PostResponse defines the JSON structure returned for a post
type PostResponse struct {
	Id           uint            `json:"id"`
	UserId       uint            `json:"userId"`
	Caption      string          `json:"caption"`
	ViewCount    int             `json:"viewCount"`
	CommentCount int             `json:"commentCount"`
	LikeCount    int             `json:"likeCount"`
	Created      time.Time       `json:"created"`
	Updated      time.Time       `json:"updated"`
	Images       []ImageResponse `json:"images"`
	Tags         []uint          `json:"tags"`
}

func NewPostResponseFromDomain(post domain.Post) PostResponse {
	var imageResponses []ImageResponse
	for _, image := range post.PostImages {
		imageResponses = append(imageResponses, NewImageResponseFromDomain(image))
	}

	return PostResponse{
		Id:           post.ID,
		UserId:       post.UserId,
		Caption:      post.Caption,
		ViewCount:    post.ViewCount,
		CommentCount: post.CommentCount,
		LikeCount:    post.LikeCount,
		Created:      post.Created,
		Updated:      post.Updated,
		Images:       imageResponses,
		Tags:         post.TagIds,
	}
}

// --------------------------------------------------
// COMMENT MODELS
// --------------------------------------------------

// CommentPayload represents the payload for creating/updating a comment
type CommentPayload struct {
	UserId  uint   `json:"userId"`
	PostId  uint   `json:"postId"`
	Comment string `json:"comment"`
}

func (c *CommentPayload) Validate() error {
	if c.UserId == 0 {
		return errors.New("userId is required")
	}
	if c.PostId == 0 {
		return errors.New("postId is required")
	}
	if c.Comment == "" {
		return errors.New("comment is required")
	}
	return nil
}

func (c *CommentPayload) ToDomainComment() domain.PostComment {
	return domain.PostComment{
		UserId:  c.UserId,
		PostId:  c.PostId,
		Comment: c.Comment,
	}
}

// CommentResponse defines the JSON structure returned for a comment
type CommentResponse struct {
	Id         uint      `json:"id"`
	UserId     uint      `json:"userId"`
	PostId     uint      `json:"postId"`
	Comment    string    `json:"comment"`
	Created    time.Time `json:"created"`
	IsDeleted  bool      `json:"isDeleted"`
	IsReported bool      `json:"isReported"`
}

func NewCommentResponseFromDomain(comment domain.PostComment) CommentResponse {
	return CommentResponse{
		Id:         comment.Id,
		UserId:     comment.UserId,
		PostId:     comment.PostId,
		Comment:    comment.Comment,
		Created:    comment.Created,
		IsDeleted:  comment.IsDeleted,
		IsReported: comment.IsReported,
	}
}

// --------------------------------------------------
// REPORT MODELS
// --------------------------------------------------

// ReportPayload represents the payload for creating a report
type ReportPayload struct {
	UserId uint   `json:"userId"`
	PostId uint   `json:"postId"`
	Reason string `json:"reason"`
}

func (r *ReportPayload) Validate() error {
	if r.UserId == 0 {
		return errors.New("userId is required")
	}
	if r.PostId == 0 {
		return errors.New("postId is required")
	}
	if r.Reason == "" {
		return errors.New("reason is required")
	}
	return nil
}

func (r *ReportPayload) ToDomainReport() domain.PostReport {
	return domain.PostReport{
		UserId: r.UserId,
		PostId: r.PostId,
		Reason: r.Reason,
	}
}

// ReportResponse defines the JSON structure returned for a report
type ReportResponse struct {
	Id        uint   `json:"id"`
	UserId    uint   `json:"userId"`
	PostId    uint   `json:"postId"`
	Reason    string `json:"reason"`
	Confirmed bool   `json:"confirmed"`
	Resolved  bool   `json:"resolved"`
}

func NewReportResponseFromDomain(report domain.PostReport) ReportResponse {
	return ReportResponse{
		Id:        report.Id,
		UserId:    report.UserId,
		PostId:    report.PostId,
		Reason:    report.Reason,
		Confirmed: report.Confirmed,
		Resolved:  report.Resolved,
	}
}

func NewPaginatedReportsResponse(reports []domain.PostReport, page, size int) PaginatedReportsResponse {
	var responses []ReportResponse
	for _, report := range reports {
		responses = append(responses, NewReportResponseFromDomain(report))
	}

	return PaginatedReportsResponse{
		CurrentPage: page,
		PageSize:    size,
		TotalPages:  len(responses) / size,
		Data:        responses,
	}
}

type PaginatedReportsResponse struct {
	CurrentPage int              `json:"currentPage"`
	PageSize    int              `json:"pageSize"`
	TotalPages  int              `json:"totalPages"`
	Data        []ReportResponse `json:"data"`
}

// NewReportResponseFromDomainReport converts a domain `PostReport` into a `ReportResponse`.
func NewReportResponseFromDomainReport(report domain.PostReport) ReportResponse {
	return ReportResponse{
		Id:        report.Id,
		UserId:    report.UserId,
		PostId:    report.PostId,
		Reason:    report.Reason,
		Confirmed: report.Confirmed,
		Resolved:  report.Resolved,
	}
}

// --------------------------------------------------
// IMAGE MODELS
// --------------------------------------------------

// ImagePayload represents the payload for uploading an image
type ImagePayload struct {
	PostId   uint   `json:"postId"`
	ImageId  uint   `json:"imageId"`
	S3Url    string `json:"s3Url"`
	IsMain   bool   `json:"isMain"`
	Order    int    `json:"order"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int    `json:"size"`
	Type     string `json:"type"`
	Format   string `json:"format"`
	Exif     string `json:"exif"`
	Location string `json:"location"`
}

func (i *ImagePayload) Validate() error {
	if i.PostId == 0 {
		return errors.New("postId is required")
	}
	if i.ImageId == 0 {
		return errors.New("imageId is required")
	}
	if i.S3Url == "" {
		return errors.New("s3Url is required")
	}
	if i.Width <= 0 {
		return errors.New("width must be greater than zero")
	}
	if i.Height <= 0 {
		return errors.New("height must be greater than zero")
	}
	if i.Size <= 0 {
		return errors.New("size must be greater than zero")
	}
	return nil
}

func (i *ImagePayload) ToDomainImage() domain.PostImage {
	return domain.PostImage{
		PostId:   i.PostId,
		ImageId:  i.ImageId,
		S3Url:    i.S3Url,
		IsMain:   i.IsMain,
		Order:    i.Order,
		Width:    i.Width,
		Height:   i.Height,
		Size:     i.Size,
		Type:     i.Type,
		Format:   i.Format,
		Exif:     i.Exif,
		Location: i.Location,
		Created:  time.Now(),
	}
}

// ImageResponse defines the JSON structure returned for an image
type ImageResponse struct {
	Id     uint   `json:"id"`
	PostId uint   `json:"postId"`
	S3Url  string `json:"s3Url"`
	IsMain bool   `json:"isMain"`
	Order  int    `json:"order"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Size   int    `json:"size"`
}

func NewImageResponseFromDomain(image domain.PostImage) ImageResponse {
	return ImageResponse{
		Id:     image.Id,
		PostId: image.PostId,
		S3Url:  image.S3Url,
		IsMain: image.IsMain,
		Order:  image.Order,
		Width:  image.Width,
		Height: image.Height,
		Size:   image.Size,
	}
}

func NewPaginatedImagesResponse(postImages []domain.PostImage, page, size int) PaginatedImagesResponse {
	var imageResponses []ImageResponse
	for _, image := range postImages {
		imageResponses = append(imageResponses, NewImageResponseFromDomain(image))
	}

	totalPages := (len(imageResponses) + size - 1) / size // Calculate total pages (rounding up)

	return PaginatedImagesResponse{
		CurrentPage: page,
		PageSize:    size,
		TotalPages:  totalPages,
		Data:        imageResponses,
	}
}

type PaginatedImagesResponse struct {
	CurrentPage int             `json:"currentPage"`
	PageSize    int             `json:"pageSize"`
	TotalPages  int             `json:"totalPages"`
	Data        []ImageResponse `json:"data"`
}

// --------------------------------------------------
// GENERIC MODELS
// --------------------------------------------------

// ErrorMessage is a generic struct for JSON error responses
type ErrorMessage struct {
	Message string `json:"message"`
}

// ImageDeletePayload represents the payload for deleting an image
type ImageDeletePayload struct {
	ImageKey string `json:"imageKey"`
}

func (p *ImageDeletePayload) Validate() error {
	if p.ImageKey == "" {
		return errors.New("imageKey is required")
	}
	return nil
}
