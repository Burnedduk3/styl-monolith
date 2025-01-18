package grpc

import (
	"context"
	"log"

	pb "styl-monolith/generated/proto/media"
)

type CrudPostService struct {
	pb.UnimplementedCrudPostServiceServer // Embedding the unimplemented server to follow the interface
}

// ListPosts implements the gRPC ListPosts method
func (s *CrudPostService) ListPosts(ctx context.Context, req *pb.PaginationRequest) (*pb.PaginationResponsePost, error) {
	log.Printf("Received CreatePost request: %v", req.Limit)
	log.Printf("Received CreatePost request: %v", req.Page)
	log.Printf("Received CreatePost request: %v", req.Size)
	return &pb.PaginationResponsePost{}, nil
}

// ListComments implements the gRPC ListComments method
func (s *CrudPostService) ListComments(ctx context.Context, req *pb.PaginationRequest) (*pb.PaginationResponseComments, error) {
	log.Printf("Received CreatePost request: %v", req.Limit)
	log.Printf("Received CreatePost request: %v", req.Page)
	log.Printf("Received CreatePost request: %v", req.Size)
	// Simulate saving to a database
	return &pb.PaginationResponseComments{}, nil
}

// PostCommentsById implements the gRPC PostCommentsById method
func (s *CrudPostService) PostCommentsById(ctx context.Context, req *pb.PostCommentsByIdRequest) (*pb.PostCommentsByIdResponse, error) {
	log.Printf("Received CreatePost request: %v", req.PostId)
	return &pb.PostCommentsByIdResponse{}, nil
}
