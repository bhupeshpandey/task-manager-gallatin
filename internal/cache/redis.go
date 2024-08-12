package cache

import (
	"context"
	. "github.com/bhupeshpandey/task-manager-gallatin/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
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

func (r *redisCache) GetTask(id uuid.UUID) ([]byte, error) {
	data, err := r.client.Get(r.ctx, id.String()).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	return []byte(data), err
}

func (r *redisCache) SetTask(id uuid.UUID, data []byte) error {
	return r.client.Set(r.ctx, id.String(), data, 10*time.Minute).Err() // Cache for 10 minutes
}

func (r *redisCache) DeleteTask(id uuid.UUID) error {
	return r.client.Del(r.ctx, id.String()).Err()
}
