package cache

import (
	"fmt"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/models"
	"github.com/go-redis/redis/v8"
)

func NewCache(cfg models.CacheConfig) models.Cache {
	var cacheInst models.Cache
	if cfg.Type == "redis" {
		options := &redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		}
		client := redis.NewClient(options)
		cacheInst = newRedisCache(client)
	}
	// Add condition for more cache type,. as of only redis
	return cacheInst
}
