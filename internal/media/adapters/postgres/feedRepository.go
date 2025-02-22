package postgres

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math/rand"
	"styl-monolith/internal/media/core/domain"
	"styl-monolith/internal/media/core/ports"
	"time"
)

type FeedRepositoryStruct struct {
	log  *logrus.Logger
	conn *gorm.DB
}

// NewFeedRepository creates a new instance of FeedRepository.
func NewFeedRepository(log *logrus.Logger, conn *gorm.DB) ports.FeedRepository {
	return &FeedRepositoryStruct{log: log, conn: conn}
}

func (f *FeedRepositoryStruct) GetRandomFeed(userId uint, page, size int) ([]domain.Post, error) {
	const limit = 5
	var posts []domain.Post

	// Seed the random generator (optional, for randomness)
	rand.Seed(time.Now().UnixNano())

	// Query to retrieve posts excluding posts made by the given userId
	// Order by creation time to give priority to newer posts
	// Then apply random ordering and limit the result to 5
	result := f.conn.
		Where("user_id != ?", userId). // Exclude posts by this user
		Order("created_at DESC"). // Prioritize newer posts
		Limit(50). // Consider the newest 50 posts for randomness
		Find(&posts)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // No posts found
		}
		return nil, result.Error
	}

	if len(posts) <= limit {
		return posts, nil // If 5 or fewer posts exist, return them directly
	}

	// Randomly pick 5 posts from the query result
	randomPosts := randomSelect(posts, limit)
	return randomPosts, nil
}

// Helper function to randomly select `n` elements from a slice
func randomSelect(posts []domain.Post, n int) []domain.Post {
	postCount := len(posts)
	randIndexes := rand.Perm(postCount)[:n]
	randomPosts := make([]domain.Post, 0, n)

	for _, idx := range randIndexes {
		randomPosts = append(randomPosts, posts[idx])
	}
	return randomPosts
}
