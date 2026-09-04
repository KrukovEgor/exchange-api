package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/KrukovEgor/exchange-api/internal/config"
	redis "github.com/redis/go-redis/v9"
)

const (
	defaultHost = "0.0.0.0"

	pingTimeout = 3 * time.Second
)

func NewRedisClient(ctx context.Context, cfg *config.RedisConfig) (*redis.Client, error) {
	const op = "infra.NewRedisClient"

	if cfg == nil {
		return nil, fmt.Errorf("%s: config is nil", op)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         defaultHost + ":" + cfg.Port,
		Password:     cfg.Password,
		MinIdleConns: cfg.MinIdleConns,
		MaxIdleConns: cfg.MaxIdleConns,
	})

	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	err := rdb.Ping(pingCtx).Err()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to redis server: %w", op, err)
	}

	return rdb, nil
}
