package cache

import (
	"context"
	"errors"
	. "github.com/bhupeshpandey/task-manager-gallatin/internal/models"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisCache struct {
	client *redis.Client
	ctx    context.Context
}

func newRedisCache(client *redis.Client) Cache {
	return &redisCache{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *redisCache) GetValue(id string) ([]byte, error) {
	data, err := r.client.Get(r.ctx, id).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // Cache miss
	}
	return []byte(data), err
}

func (r *redisCache) SetValue(id string, data []byte) error {
	return r.client.Set(r.ctx, id, data, 10*time.Minute).Err() // Cache for 10 minutes
}

func (r *redisCache) DeleteEntry(id string) error {
	return r.client.Del(r.ctx, id).Err()
}
