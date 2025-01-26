package grpc

import (
	"context"
	"log"
	"styl-monolith/internal/media/core/services"

	pb "styl-monolith/generated/proto/media"
)

type CrudPostGrcp struct {
	pb.UnimplementedCrudPostServiceServer
	mediaService services.MediaService
}

func NewCrudPostGrcp(mediaService services.MediaService) *CrudPostGrcp {
	return &CrudPostGrcp{mediaService: mediaService}
}

// ListPosts implements the gRPC ListPosts method
func (s *CrudPostGrcp) ListPosts(ctx context.Context, req *pb.PaginationRequest) (*pb.PaginationResponsePost, error) {
	log.Printf("Received CreatePost request: %v", req.Page)
	log.Printf("Received CreatePost request: %v", req.Size)
	posts, err := s.mediaService.ListPosts(int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}
	var postsResponse []*pb.Post
	for _, post := range posts {
		postsResponse = append(postsResponse, post.ToProtoDomain())
	}

	return &pb.PaginationResponsePost{
		Post: postsResponse,
	}, nil
}

// ListComments implements the gRPC ListComments method
func (s *CrudPostGrcp) ListComments(ctx context.Context, req *pb.PaginationRequest) (*pb.PaginationResponseComments, error) {
	log.Printf("Received CreatePost request: %v", req.Page)
	log.Printf("Received CreatePost request: %v", req.Size)
	// Simulate saving to a database
	return &pb.PaginationResponseComments{}, nil
}

// PostCommentsById implements the gRPC PostCommentsById method
func (s *CrudPostGrcp) PostCommentsById(ctx context.Context, req *pb.PostCommentsByIdRequest) (*pb.PostCommentsByIdResponse, error) {
	log.Printf("Received CreatePost request: %v", req.PostId)
	return &pb.PostCommentsByIdResponse{}, nil
}
