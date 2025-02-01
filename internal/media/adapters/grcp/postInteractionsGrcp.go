package grpc

import (
	"context"
	"github.com/sirupsen/logrus"
	"styl-monolith/internal/media/core/services"

	pb "styl-monolith/generated/proto/media"
)

type CrudPostGrcp struct {
	pb.UnimplementedCrudPostServiceServer
	mediaService services.MediaService
	log          *logrus.Logger
}

func NewCrudPostGrcp(mediaService services.MediaService, logger *logrus.Logger) *CrudPostGrcp {
	return &CrudPostGrcp{mediaService: mediaService, log: logger}
}

// ListPosts implements the gRPC ListPosts method
func (s *CrudPostGrcp) ListPosts(ctx context.Context, req *pb.PaginationRequest) (*pb.PaginationResponsePost, error) {
	s.log.Debug("Received ListPosts GRPC request")
	posts, TotalPages, err := s.mediaService.ListPosts(int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}
	var postsResponse []*pb.Post
	for _, post := range posts {
		postsResponse = append(postsResponse, post.ToProtoDomain())
	}

	return &pb.PaginationResponsePost{
		CurrentPage: uint64(req.Page),
		PageSize:    uint64(req.Size),
		TotalPages:  uint64(TotalPages),
		Post:        postsResponse,
	}, nil
}

// ListComments implements the gRPC ListComments method
func (s *CrudPostGrcp) ListComments(ctx context.Context, req *pb.PaginationRequest) (*pb.PaginationResponseComments, error) {
	dcomments, TotalPages, err := s.mediaService.ListComments(uint(req.PostId), int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}
	var postsCommentsResponse []*pb.PostComment
	for _, comment := range dcomments {
		postsCommentsResponse = append(postsCommentsResponse, comment.ToProtoDomain())
	}
	// Simulate saving to a database
	return &pb.PaginationResponseComments{
		CurrentPage: uint64(req.Page),
		PageSize:    uint64(req.Size),
		TotalPages:  uint64(TotalPages),
		Comments:    postsCommentsResponse,
	}, nil
}

// PostCommentsById implements the gRPC PostCommentsById method
func (s *CrudPostGrcp) PostCommentsById(ctx context.Context, req *pb.PostCommentsByIdRequest) (*pb.PostCommentsByIdResponse, error) {
	s.log.Printf("Received CreatePost request: %v", req.PostId)
	return &pb.PostCommentsByIdResponse{}, nil
}
