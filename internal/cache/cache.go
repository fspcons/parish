package cache

//go:generate moq -out mocks.go . Cache

import (
	"context"
	"errors"
	"time"
)

// ErrNotFound is returned when a key does not exist in the cache.
var ErrNotFound = errors.New("cache: key not found")

// Cache defines a general-purpose cache interface used for token storage
// and rate limiting.
type Cache interface {
	// Set stores a key-value pair with a time-to-live.
	Set(ctx context.Context, key string, value string, ttl time.Duration) error

	// Get retrieves the value for a key. Returns ErrNotFound when the key
	// does not exist.
	Get(ctx context.Context, key string) (string, error)

	// Del removes a key from the cache.
	Del(ctx context.Context, key string) error

	// ZAdd adds a member with a score to a sorted set.
	ZAdd(ctx context.Context, key string, score float64, member string) error

	// ZRemRangeByScore removes members from a sorted set whose scores fall
	// within the given range (inclusive, string representation).
	ZRemRangeByScore(ctx context.Context, key, min, max string) error

	// ZCard returns the number of members in a sorted set.
	ZCard(ctx context.Context, key string) (int64, error)

	// Expire sets a time-to-live on an existing key.
	Expire(ctx context.Context, key string, ttl time.Duration) error
}
