package grpc

import (
	"context"
	"log"

	pb "styl-monolith/generated/proto/media"
)

type CrudPostService struct {
	pb.UnimplementedCrudPostServiceServer // Embedding the unimplemented server to follow the interface
}

// CreatePostLike implements the gRPC CreatePostLike method
func (s *CrudPostService) CreatePostLike(ctx context.Context, req *pb.CreateLikeRequest) (*pb.CreateLikeResponse, error) {
	log.Printf("Received CreatePost request: %v", req.Content)
	// Simulate saving to a database
	newPostID := "12345"
	return &pb.CreateLikeResponse{
		Id: newPostID,
	}, nil
}
