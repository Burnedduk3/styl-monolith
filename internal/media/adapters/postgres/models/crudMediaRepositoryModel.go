package models

import (
	"gorm.io/gorm"
	"styl-monolith/internal/media/core/domain"
)

// Post represents the main structure for social media posts in the database.
type Post struct {
	gorm.Model
	UserID          uint          `gorm:"index"`
	ViewCount       int           `gorm:"default:0"`
	CommentCount    int           `gorm:"default:0"`
	LikeCount       int           `gorm:"default:0"`
	NumberOfReports int           `gorm:"default:0"`
	SavedCount      int           `gorm:"default:0"`
	Caption         string        `gorm:"type:text"`
	IsDeleted       bool          `gorm:"default:false"`
	IsPublic        bool          `gorm:"default:true"`
	IsApproved      bool          `gorm:"default:false"`
	IsSpam          bool          `gorm:"default:false"`
	IsReported      bool          `gorm:"default:false"`
	PostImages      []PostImage   `gorm:"foreignKey:PostID"`
	PostComments    []PostComment `gorm:"foreignKey:PostID"`
	PostLikes       []PostLike    `gorm:"foreignKey:PostID"`
	PostReports     []PostReport  `gorm:"foreignKey:PostID"`
	PostSaves       []PostSaved   `gorm:"foreignKey:PostID"`
}

// ToPostDomain converts a Post model instance to a domain.Post.
func (p *Post) ToPostDomain() domain.Post {
	post := domain.Post{
		Id:              p.ID,
		UserId:          p.UserID,
		ViewCount:       p.ViewCount,
		CommentCount:    p.CommentCount,
		LikeCount:       p.LikeCount,
		NumberOfReports: p.NumberOfReports,
		SavedCount:      p.SavedCount,
		Caption:         p.Caption,
		IsDeleted:       p.IsDeleted,
		IsPublic:        p.IsPublic,
		IsApproved:      p.IsApproved,
		IsSpam:          p.IsSpam,
		IsReported:      p.IsReported,
		PostImages:      []domain.PostImage{},
		PostComments:    []domain.PostComment{},
		PostLikes:       []domain.PostLike{},
		PostReports:     []domain.PostReport{},
		PostSave:        []domain.PostSaved{},
	}

	// Convert associated PostImages
	for _, img := range p.PostImages {
		post.PostImages = append(post.PostImages, img.ToPostImageDomain())
	}

	// Convert associated PostComments
	for _, comment := range p.PostComments {
		post.PostComments = append(post.PostComments, comment.ToPostCommentDomain())
	}

	// Convert associated PostLikes
	for _, like := range p.PostLikes {
		post.PostLikes = append(post.PostLikes, like.ToPostLikeDomain())
	}

	// Convert associated PostReports
	for _, report := range p.PostReports {
		post.PostReports = append(post.PostReports, report.ToPostReportDomain())
	}

	// Convert associated PostSaves
	for _, save := range p.PostSaves {
		post.PostSave = append(post.PostSave, save.ToPostSavedDomain())
	}

	return post
}

// PostLike represents user likes on a particular post.
type PostLike struct {
	gorm.Model
	UserID uint `gorm:"index"`
	PostID uint `gorm:"index"`
}

// ToPostLikeDomain converts a PostLike model instance to a domain.PostLike.
func (pl *PostLike) ToPostLikeDomain() domain.PostLike {
	return domain.PostLike{
		Id:     pl.ID,
		UserId: pl.UserID,
		PostId: pl.PostID,
	}
}

// PostComment represents user comments for a post.
type PostComment struct {
	gorm.Model
	UserID     uint   `gorm:"index"`
	PostID     uint   `gorm:"index"`
	Comment    string `gorm:"type:text"`
	IsDeleted  bool   `gorm:"default:false"`
	IsReported bool   `gorm:"default:false"`
}

// ToPostCommentDomain converts a PostComment model instance to a domain.PostComment.
func (pc *PostComment) ToPostCommentDomain() domain.PostComment {
	return domain.PostComment{
		Id:         pc.ID,
		UserId:     pc.UserID,
		PostId:     pc.PostID,
		Comment:    pc.Comment,
		IsDeleted:  pc.IsDeleted,
		IsReported: pc.IsReported,
	}
}

// PostReport represents a report made by a user on a particular post.
type PostReport struct {
	gorm.Model
	UserID    uint   `gorm:"index"`
	PostID    uint   `gorm:"index"`
	Reason    string `gorm:"type:text"`
	Confirmed bool   `gorm:"default:false"`
	Resolved  bool   `gorm:"default:false"`
}

// ToPostReportDomain converts a PostReport model instance to a domain.PostReport.
func (pr *PostReport) ToPostReportDomain() domain.PostReport {
	return domain.PostReport{
		Id:        pr.ID,
		UserId:    pr.UserID,
		PostId:    pr.PostID,
		Reason:    pr.Reason,
		Confirmed: pr.Confirmed,
		Resolved:  pr.Resolved,
	}
}

// PostSaved represents when a user saves/bookmarks a post.
type PostSaved struct {
	gorm.Model
	UserID uint `gorm:"index"`
	PostID uint `gorm:"index"`
}

// ToPostSavedDomain converts a PostSaved model instance to a domain.PostSaved.
func (ps *PostSaved) ToPostSavedDomain() domain.PostSaved {
	return domain.PostSaved{
		Id:     ps.ID,
		UserId: ps.UserID,
		PostId: ps.PostID,
	}
}

// PostImage contains images associated with a post.
type PostImage struct {
	gorm.Model
	PostID    uint   `gorm:"index"`
	ImageID   uint   `gorm:"index"`
	S3URL     string `gorm:"type:text"`
	IsDeleted bool   `gorm:"default:false"`
	IsMain    bool   `gorm:"default:false"`
	Order     int
	Width     int
	Height    int
	Size      int
	Type      string
	Format    string
	Exif      string
	Location  string
}

// ToPostImageDomain converts a PostImage model instance to a domain.PostImage.
func (pi *PostImage) ToPostImageDomain() domain.PostImage {
	return domain.PostImage{
		Id:        pi.ID,
		PostId:    pi.PostID,
		ImageId:   pi.ImageID,
		S3Url:     pi.S3URL,
		IsDeleted: pi.IsDeleted,
		IsMain:    pi.IsMain,
		Order:     pi.Order,
		Width:     pi.Width,
		Height:    pi.Height,
		Size:      pi.Size,
		Type:      pi.Type,
		Format:    pi.Format,
		Exif:      pi.Exif,
		Location:  pi.Location,
	}
}

// NewPostgresPostFromDomainPost converts a domain.Post into a Postgres-compatible Post model.
func NewPostgresPostFromDomainPost(post domain.Post) Post {
	// Convert the main Post structure
	mp := Post{
		Model: gorm.Model{
			ID: post.Id,
		},
		UserID:          post.UserId,
		ViewCount:       post.ViewCount,
		CommentCount:    post.CommentCount,
		LikeCount:       post.LikeCount,
		NumberOfReports: post.NumberOfReports,
		SavedCount:      post.SavedCount,
		Caption:         post.Caption,
		IsDeleted:       post.IsDeleted,
		IsPublic:        post.IsPublic,
		IsApproved:      post.IsApproved,
		IsSpam:          post.IsSpam,
		IsReported:      post.IsReported,
		PostImages:      make([]PostImage, 0),
		PostComments:    make([]PostComment, 0),
		PostLikes:       make([]PostLike, 0),
		PostReports:     make([]PostReport, 0),
		PostSaves:       make([]PostSaved, 0),
	}

	// Convert PostImages
	for _, image := range post.PostImages {
		mp.PostImages = append(mp.PostImages, PostImage{
			Model: gorm.Model{
				ID: image.Id,
			},
			PostID:    image.PostId,
			ImageID:   image.ImageId,
			S3URL:     image.S3Url,
			IsDeleted: image.IsDeleted,
			IsMain:    image.IsMain,
			Order:     image.Order,
			Width:     image.Width,
			Height:    image.Height,
			Size:      image.Size,
			Type:      image.Type,
			Format:    image.Format,
			Exif:      image.Exif,
			Location:  image.Location,
		})
	}

	// Convert PostComments
	for _, comment := range post.PostComments {
		mp.PostComments = append(mp.PostComments, NewPostgresPostCommentFromDomainPostComment(comment))
	}

	// Convert PostLikes
	for _, like := range post.PostLikes {
		mp.PostLikes = append(mp.PostLikes, NewPostgresPostLikeFromDomainPostLike(like))
	}

	// Convert PostReports
	for _, report := range post.PostReports {
		mp.PostReports = append(mp.PostReports, NewPostgresPostReportFromDomainPostReport(report))
	}

	// Convert PostSaves
	for _, save := range post.PostSave {
		mp.PostSaves = append(mp.PostSaves, NewPostgresPostSavedFromDomainPostSaved(save))
	}

	return mp
}

// NewPostgresPostLikeFromDomainPostLike converts a domain.PostLike into a Postgres-compatible PostLike model.
func NewPostgresPostLikeFromDomainPostLike(like domain.PostLike) PostLike {
	return PostLike{
		Model: gorm.Model{
			ID: like.Id,
		},
		UserID: like.UserId,
		PostID: like.PostId,
	}
}

// NewPostgresPostCommentFromDomainPostComment converts a domain.PostComment into a Postgres-compatible PostComment model.
func NewPostgresPostCommentFromDomainPostComment(comment domain.PostComment) PostComment {
	return PostComment{
		Model: gorm.Model{
			ID: comment.Id,
		},
		UserID:     comment.UserId,
		PostID:     comment.PostId,
		Comment:    comment.Comment,
		IsDeleted:  comment.IsDeleted,
		IsReported: comment.IsReported,
	}
}

// NewPostgresPostReportFromDomainPostReport converts a domain.PostReport into a Postgres-compatible PostReport model.
func NewPostgresPostReportFromDomainPostReport(report domain.PostReport) PostReport {
	return PostReport{
		Model: gorm.Model{
			ID: report.Id,
		},
		UserID:    report.UserId,
		PostID:    report.PostId,
		Reason:    report.Reason,
		Confirmed: report.Confirmed,
		Resolved:  report.Resolved,
	}
}

// NewPostgresPostSavedFromDomainPostSaved converts a domain.PostSaved into a Postgres-compatible PostSaved model.
func NewPostgresPostSavedFromDomainPostSaved(saved domain.PostSaved) PostSaved {
	return PostSaved{
		Model: gorm.Model{
			ID: saved.Id,
		},
		UserID: saved.UserId,
		PostID: saved.PostId,
	}
}

// NewPostgresPostImageFromDomainPostImage converts a domain.PostImage into a Postgres-compatible PostImage model.
func NewPostgresPostImageFromDomainPostImage(image domain.PostImage) PostImage {
	return PostImage{
		Model: gorm.Model{
			ID: image.Id,
		},
		PostID:    image.PostId,
		ImageID:   image.ImageId,
		S3URL:     image.S3Url,
		IsDeleted: image.IsDeleted,
		IsMain:    image.IsMain,
		Order:     image.Order,
		Width:     image.Width,
		Height:    image.Height,
		Size:      image.Size,
		Type:      image.Type,
		Format:    image.Format,
		Exif:      image.Exif,
		Location:  image.Location,
	}
}
