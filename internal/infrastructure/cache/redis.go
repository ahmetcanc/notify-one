package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type RedisCache struct {
	client *redis.Client
}

// NewRedisConnection initializes a new Redis wrapper
func NewRedisConnection(ctx context.Context, cfg Config) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("redis connection established successfully")
	return &RedisCache{client: rdb}, nil
}

// PushToQueue adds notification ID to a Redis list
func (r *RedisCache) PushToQueue(ctx context.Context, notificationID string) error {
	return r.client.LPush(ctx, "notification_queue", notificationID).Err()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

// PopFromQueue blocks until a notification ID is available in the queue
func (r *RedisCache) PopFromQueue(ctx context.Context, queueName string) ([]string, error) {
	// 0 means wait indefinitely
	return r.client.BLPop(ctx, 0, queueName).Result()
}
