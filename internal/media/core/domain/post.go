package domain

import "time"

type Post struct {
	ID              uint
	UserId          uint
	ViewCount       int
	CommentCount    int
	LikeCount       int
	NumberOfReports int
	SavedCount      int
	Caption         string
	Created         time.Time
	Updated         time.Time
	Deleted         time.Time
	IsDeleted       bool
	IsPublic        bool
	IsApproved      bool
	IsSpam          bool
	IsReported      bool
	TagIds          []uint // from tag domain
	PostImages      []PostImage
	PostComments    []PostComment
	PostLikes       []PostLike
	PostReports     []PostReport
}

type PostLike struct {
	Id      uint
	UserId  uint
	PostId  uint
	Created time.Time
}

type PostComment struct {
	Id         uint
	UserId     uint
	PostId     uint
	Comment    string
	Created    time.Time
	Deleted    time.Time
	IsDeleted  bool
	IsReported bool
}

type PostReport struct {
	Id        uint
	UserId    uint
	PostId    uint
	Reason    string
	Confirmed bool
}

type PostSaved struct {
	Id      uint
	UserId  uint
	PostId  uint
	Created time.Time
}

type PostImage struct {
	Id        uint
	PostId    uint
	ImageId   uint
	S3Url     string
	Created   time.Time
	Deleted   time.Time
	IsDeleted bool
	IsMain    bool
	Order     int
	Width     int
	Height    int
	Size      int
	Type      string
	Format    string
	Exif      string
	Location  string
}
