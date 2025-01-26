package domain

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "styl-monolith/generated/proto/media"
	"time"
)

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
	Resolved  bool
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

func (p *Post) ToProtoDomain() *pb.Post {
	return &pb.Post{
		Id:              uint64(p.ID),
		UserId:          uint64(p.UserId),
		ViewCount:       uint32(p.ViewCount),
		CommentCount:    uint32(p.CommentCount),
		LikeCount:       uint32(p.LikeCount),
		NumberOfReports: uint32(p.NumberOfReports),
		SavedCount:      uint32(p.SavedCount),
		Caption:         p.Caption,
		Created:         timestamppb.New(p.Created),
		Updated:         timestamppb.New(p.Updated),
		Deleted:         timestamppb.New(p.Deleted),
		IsDeleted:       p.IsDeleted,
		IsPublic:        p.IsPublic,
		IsApproved:      p.IsApproved,
		IsSpam:          p.IsSpam,
		IsReported:      p.IsReported,
		TagIds:          convertUintSliceToUint64Slice(p.TagIds),
		PostImages:      convertPostImagesToProto(p.PostImages),
		PostComments:    convertPostCommentsToProto(p.PostComments),
		PostLikes:       convertPostLikesToProto(p.PostLikes),
		PostReports:     convertPostReportsToProto(p.PostReports),
	}
}

// Helper to convert domain.PostImage slices to proto equivalents
func convertPostImagesToProto(images []PostImage) []*pb.PostImage {
	var protoImages []*pb.PostImage
	for _, image := range images {
		protoImages = append(protoImages, &pb.PostImage{
			Id:        uint64(image.Id),
			PostId:    uint64(image.PostId),
			ImageId:   uint64(image.ImageId),
			S3Url:     image.S3Url,
			Created:   timestamppb.New(image.Created),
			Deleted:   timestamppb.New(image.Deleted),
			IsDeleted: image.IsDeleted,
			IsMain:    image.IsMain,
			Order:     int32(image.Order),
			Width:     int32(image.Width),
			Height:    int32(image.Height),
			Size:      int32(image.Size),
			Type:      image.Type,
			Format:    image.Format,
			Exif:      image.Exif,
			Location:  image.Location,
		})
	}
	return protoImages
}

// Helper to convert domain.PostComment slices to proto equivalents
func convertPostCommentsToProto(comments []PostComment) []*pb.PostComment {
	var protoComments []*pb.PostComment
	for _, comment := range comments {
		protoComments = append(protoComments, &pb.PostComment{
			Id:         uint64(comment.Id),
			UserId:     uint64(comment.UserId),
			PostId:     uint64(comment.PostId),
			Comment:    comment.Comment,
			Created:    timestamppb.New(comment.Created),
			IsDeleted:  comment.IsDeleted,
			IsReported: comment.IsReported,
		})
	}
	return protoComments
}

// Helper to convert domain.PostLike slices to proto equivalents
func convertPostLikesToProto(likes []PostLike) []*pb.PostLike {
	var protoLikes []*pb.PostLike
	for _, like := range likes {
		protoLikes = append(protoLikes, &pb.PostLike{
			Id:      uint64(like.Id),
			UserId:  uint64(like.UserId),
			PostId:  uint64(like.PostId),
			Created: timestamppb.New(like.Created),
		})
	}
	return protoLikes
}

// Helper to convert domain.PostReport slices to proto equivalents
func convertPostReportsToProto(reports []PostReport) []*pb.PostReport {
	var protoReports []*pb.PostReport
	for _, report := range reports {
		protoReports = append(protoReports, &pb.PostReport{
			Id:        uint64(report.Id),
			UserId:    uint64(report.UserId),
			PostId:    uint64(report.PostId),
			Reason:    report.Reason,
			Confirmed: report.Confirmed,
		})
	}
	return protoReports
}

// Helper to convert []uint to []uint64
func convertUintSliceToUint64Slice(input []uint) []uint64 {
	var result []uint64
	for _, value := range input {
		result = append(result, uint64(value))
	}
	return result
}
