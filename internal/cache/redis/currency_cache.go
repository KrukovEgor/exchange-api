package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/KrukovEgor/exchange-api/internal/config"
	"github.com/KrukovEgor/exchange-api/internal/domain"
	"github.com/redis/go-redis/v9"
)

const (
	basePrefix  = "exchange:currencies:v1"
	staleSuffix = "stale"
)

type CurrencyCache struct {
	rdb          *redis.Client
	ttl          time.Duration
	staleTTL     time.Duration
	refreshDelta time.Duration
}

func NewCurrencyCache(rdb *redis.Client, cfg *config.CurrencyCacheConfig) (*CurrencyCache, error) {
	const op = "redis.NewCurrencyCache"

	if rdb == nil {
		return nil, domain.NewInvalidInputError(op, "redis client is nil", nil)
	}

	return &CurrencyCache{
		rdb:          rdb,
		ttl:          cfg.BaseTTL,
		staleTTL:     cfg.StaleTTL,
		refreshDelta: cfg.RefreshDelta,
	}, nil
}

func (c *CurrencyCache) Set(parentCtx context.Context, dir domain.Direction, value []domain.Currency) error {
	const op = "redis.CurrencyCache.Set"

	if !dir.Valid() {
		return domain.NewInvalidInputError(op, "invalid direction", nil)
	}

	if len(value) == 0 {
		return domain.NewInvalidInputError(op, "empty value", nil)
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return domain.NewInternalError(op, "failed to marshal data", err)
	}

	pipe := c.rdb.TxPipeline()

	pipe.Set(parentCtx, buildBaseKey(dir), bytes, c.ttl)
	pipe.Set(parentCtx, buildStaleKey(dir), bytes, c.ttl)

	_, err = pipe.Exec(parentCtx)
	if err != nil {
		return mapRedisError(op, "failed to execute set pipeline", err)
	}

	return nil
}

func (c *CurrencyCache) Get(parentCtx context.Context, dir domain.Direction) ([]domain.Currency, error) {
	const op = "redis.CurrencyCache.Get"

	if !dir.Valid() {
		return nil, domain.NewInvalidInputError(op, "invalid direction", nil)
	}

	var data []domain.Currency

	err := get(c.rdb, parentCtx, buildBaseKey(dir), &data)
	if err != nil {
		var unmarshalErr *json.UnsupportedTypeError
		if errors.As(err, &unmarshalErr) {
			return nil, domain.NewInternalError(op, "failed to unmarshal data", err)
		}

		return nil, mapRedisError(op, "failed to get data", err)
	}

	return data, nil
}

func buildBaseKey(dir domain.Direction) string {
	return basePrefix + ":" + dir.String()
}

func buildStaleKey(dir domain.Direction) string {
	return buildBaseKey(dir) + ":" + staleSuffix
}
