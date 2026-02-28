package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements Cache backed by a Redis server.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis-backed cache. addr is in the form "host:port".
func NewRedisCache(addr string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

// Close releases the underlying Redis connection.
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// Set stores a key-value pair with a TTL.
func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves the value for a key. Returns ErrNotFound when the key does
// not exist.
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	}
	return val, err
}

// Del removes a key from the cache.
func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// ZAdd adds a member with a score to a sorted set.
func (r *RedisCache) ZAdd(ctx context.Context, key string, score float64, member string) error {
	return r.client.ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Err()
}

// ZRemRangeByScore removes members whose scores fall within [min, max].
func (r *RedisCache) ZRemRangeByScore(ctx context.Context, key, min, max string) error {
	return r.client.ZRemRangeByScore(ctx, key, min, max).Err()
}

// ZCard returns the cardinality (number of members) of a sorted set.
func (r *RedisCache) ZCard(ctx context.Context, key string) (int64, error) {
	return r.client.ZCard(ctx, key).Result()
}

// Expire sets a TTL on an existing key.
func (r *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.client.Expire(ctx, key, ttl).Err()
}
