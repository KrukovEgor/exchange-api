package redis

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"

	"github.com/KrukovEgor/exchange-api/internal/domain"
)

func mapRedisError(op, msg string, err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, redis.Nil):
		return domain.NewCacheMissError(op, "cache miss", err)
	case errors.Is(err, context.Canceled):
		return domain.NewInternalError(op, "cache request canceled by caller", err)
	default:
		return domain.NewCacheUnavailableError(op, msg, err)
	}
}
