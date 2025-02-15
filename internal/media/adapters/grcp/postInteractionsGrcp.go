package grpc

import (
	"context"
	"github.com/sirupsen/logrus"
	pb "styl-monolith/generated/proto/media"
	"styl-monolith/internal/media/core/services"
)

// CrudPostGrcp is a gRPC server for handling CRUD operations on posts and related entities.
// It embeds pb.UnimplementedCrudPostServiceServer to ensure forward compatibility.
// It uses MediaService for media-related operations and logrus.Logger for logging.
type CrudPostGrcp struct {
	pb.UnimplementedCrudPostServiceServer
	mediaService services.MediaService
	log          *logrus.Logger
}

// NewCrudPostGrcp initializes and returns a new instance of CrudPostGrcp with the provided MediaService and logger.
func NewCrudPostGrcp(mediaService services.MediaService, logger *logrus.Logger) *CrudPostGrcp {
	return &CrudPostGrcp{mediaService: mediaService, log: logger}
}

// ListPosts fetches a paginated list of posts based on the provided page and size from the request.
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

// ListComments retrieves a paginated list of comments for a specific post, given pagination request details like page and size.
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

// PostCommentsById retrieves comments associated with a specific post ID and returns the response containing the post and comments.
func (s *CrudPostGrcp) PostCommentsById(ctx context.Context, req *pb.PostCommentsByIdRequest) (*pb.PostCommentsByIdResponse, error) {
	s.log.Printf("Received CreatePost request: %v", req.PostId)
	return &pb.PostCommentsByIdResponse{}, nil
}

// GetUserProfile fetches a paginated list of posts related to a specific user based on the provided user information.
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

// GetUserFeed retrieves a paginated feed of posts for a specific user based on the provided request parameters.
func (s *CrudPostGrcp) GetUserFeed(ctx context.Context, req *pb.UserFeed) (*pb.PaginationResponsePost, error) {
	return nil, nil
}
