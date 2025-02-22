package services

import (
	"github.com/sirupsen/logrus"
	"styl-monolith/internal/media/core/domain"
	"styl-monolith/internal/media/core/ports"
)

// FeedService defines an interface for handling media operations such as uploads, reports, posts, comments, and likes.
type FeedService interface {
	GetRandomFeed(userId uint, page, size int) ([]domain.Post, error)
}

// FeedServiceStruct defines the structure for managing media services and integrations with various ports and resources.
type FeedServiceStruct struct {
	logger   *logrus.Logger
	feedPort ports.FeedPort
}

// NewFeedService initializes and returns a FeedService instance with the provided dependencies and configurations.
func NewFeedService(
	log *logrus.Logger,
	feedPort ports.FeedPort,
) FeedService {
	return &FeedServiceStruct{
		logger:   log,
		feedPort: feedPort,
	}
}

func (fs *FeedServiceStruct) GetRandomFeed(userId uint, page, size int) ([]domain.Post, error) {
	return fs.feedPort.GetRandomFeed(userId)
}
