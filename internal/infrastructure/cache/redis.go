package cache

import (
	"context"
	"fmt"
	"log"
	"time"

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

// PushToQueue adds notification ID to a specific priority list
func (r *RedisCache) PushToQueue(ctx context.Context, queueName string, notificationID string) error {
	return r.client.LPush(ctx, queueName, notificationID).Err()
}

// PopFromQueues blocks until an ID is available in the given order of priority
func (r *RedisCache) PopFromQueues(ctx context.Context, queueNames ...string) ([]string, error) {
	return r.client.BLPop(ctx, 0, queueNames...).Result()
}

// IsRateLimited checks if the specific channel has exceeded its limit
func (r *RedisCache) IsRateLimited(ctx context.Context, channel string, limit int) (bool, error) {
	// Create a key for the current second: rate_limit:email:1709661520
	key := fmt.Sprintf("rate_limit:%s:%d", channel, time.Now().Unix())

	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	// Set expiration to 2 seconds to clean up memory
	if count == 1 {
		r.client.Expire(ctx, key, 2*time.Second)
	}

	return count > int64(limit), nil
}

func (r *RedisCache) GetQueueDepth(ctx context.Context, queueName string) (int64, error) {
	return r.client.LLen(ctx, queueName).Result()
}

// AddToDelayedQueue adds a task to the sorted set with a specific execution time (score)
func (r *RedisCache) AddToDelayedQueue(ctx context.Context, key string, value string, score int64) error {
	return r.client.ZAdd(ctx, key, redis.Z{
		Score:  float64(score),
		Member: value,
	}).Err()
}

// GetAndRemoveReadyTasks finds and removes tasks that are ready to be processed
func (r *RedisCache) GetAndRemoveReadyTasks(ctx context.Context, key string, maxScore int64) ([]string, error) {
	var luaScript = redis.NewScript(`
        local val = redis.call('zrangebyscore', KEYS[1], '-inf', ARGV[1])
        if #val > 0 then
            redis.call('zremrangebyscore', KEYS[1], '-inf', ARGV[1])
        end
        return val
    `)

	res, err := luaScript.Run(ctx, r.client, []string{key}, maxScore).StringSlice()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	return res, nil
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
