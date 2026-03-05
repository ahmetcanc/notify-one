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

// NewRedisConnection initializes a new Redis client with the given configuration
func NewRedisConnection(ctx context.Context, cfg Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test the connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("failed to connect to redis: %v", err)
		return nil, err
	}

	log.Println("redis connection established successfully")
	return rdb, nil
}
