package grpc

import (
	"context"
	"github.com/sirupsen/logrus"
	pb "styl-monolith/generated/proto/media"
	"styl-monolith/internal/media/core/services"
)

// CrudPostGrcp is a gRPC service implementation for managing CRUD operations for posts and associated actions.
// It combines a MediaService interface for business logic and logging for request tracking and debugging.
type CrudPostGrcp struct {
	pb.UnimplementedCrudPostServiceServer
	mediaService services.MediaService
	feedService  services.FeedService
	log          *logrus.Logger
}

// NewCrudPostGrcp creates a new instance of CrudPostGrcp with the provided MediaService and logger dependencies.
func NewCrudPostGrcp(mediaService services.MediaService, feedService services.FeedService, logger *logrus.Logger) *CrudPostGrcp {
	return &CrudPostGrcp{mediaService: mediaService, log: logger, feedService: feedService}
}

// ListPosts retrieves a paginated list of posts based on the provided pagination request.
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

// ListComments retrieves a paginated list of comments for a specific post based on the provided post ID, page, and size.
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

// PostCommentsById retrieves comments for a specific post by its PostId.
// It returns a PostCommentsByIdResponse containing the post ID and associated comments.
func (s *CrudPostGrcp) PostCommentsById(ctx context.Context, req *pb.PostCommentsByIdRequest) (*pb.PostCommentsByIdResponse, error) {
	s.log.Printf("Received CreatePost request: %v", req.PostId)
	return &pb.PostCommentsByIdResponse{}, nil
}

// GetUserProfile retrieves a paginated list of posts by a specific user and returns it in the response.
// It accepts the user ID, current page, and page size via the request and communicates with the media service for data retrieval.
// Returns a PaginationResponsePost containing the list of posts, current page, page size, and total pages, or an error.
func (s *CrudPostGrcp) GetUserProfile(ctx context.Context, req *pb.UserProfile) (*pb.PaginationResponsePost, error) {
	reqUserId := req.User.UserId
	page := req.PagRequest.Page
	size := req.PagRequest.Size
	posts, totalPages, err := s.mediaService.ListPostsByUser(uint(reqUserId), int(page), int(size))
	if err != nil {
		return nil, err
	}
	var postsResponse []*pb.Post
	for _, post := range posts {
		postsResponse = append(postsResponse, post.ToProtoDomain())
	}

	return &pb.PaginationResponsePost{
		CurrentPage: uint64(req.PagRequest.Page),
		PageSize:    uint64(req.PagRequest.Size),
		TotalPages:  uint64(totalPages),
		Post:        postsResponse,
	}, nil
}

func (s *CrudPostGrcp) GetUserFeed(ctx context.Context, req *pb.UserFeed) (*pb.PaginationResponsePost, error) {
	reqUserId := req.User.UserId
	page := req.PagRequest.Page
	size := req.PagRequest.Size

	posts, err := s.feedService.GetRandomFeed(uint(reqUserId), int(page), int(size))
	if err != nil {
		return nil, err
	}

	var postsResponse []*pb.Post
	for _, post := range posts {
		postsResponse = append(postsResponse, post.ToProtoDomain())
	}

	return &pb.PaginationResponsePost{
		CurrentPage: uint64(page),
		PageSize:    uint64(size),
		TotalPages:  0,
		Post:        postsResponse,
	}, nil
}
