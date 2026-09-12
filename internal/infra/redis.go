package infra

import (
	"context"
	"net"

	"github.com/KrukovEgor/exchange-api/internal/config"
	"github.com/KrukovEgor/exchange-api/internal/domain"
	redis "github.com/redis/go-redis/v9"
)

func NewRedisClient(parentCtx context.Context, cfg *config.RedisConfig) (*redis.Client, error) {
	const op = "infra.NewRedisClient"

	if cfg == nil {
		return nil, domain.NewInvalidInputError(op, "config is nil", nil)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         net.JoinHostPort(cfg.Host, cfg.Port),
		Password:     cfg.Password,
		MinIdleConns: cfg.MinIdleConns,
		MaxIdleConns: cfg.MaxIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	pingCtx, cancel := context.WithTimeout(parentCtx, cfg.DialTimeout)
	defer cancel()

	err := rdb.Ping(pingCtx).Err()
	if err != nil {
		_ = rdb.Close()

		return nil, domain.NewInternalError(op, "failed to connect to redis server", err)
	}

	return rdb, nil
}
