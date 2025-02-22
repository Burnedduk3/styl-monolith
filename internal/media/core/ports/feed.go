package ports

import "styl-monolith/internal/media/core/domain"

type FeedPort interface {
	GetRandomFeed(userId uint, page, size int) ([]domain.Post, error)
}
