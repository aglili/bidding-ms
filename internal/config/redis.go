package config

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectToRedis(cfg *Config) (*redis.Client, error) {

	redisURL := cfg.RedisURL

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return client, nil
}
