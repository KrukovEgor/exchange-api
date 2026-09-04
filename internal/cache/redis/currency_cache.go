package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/KrukovEgor/exchange-api/internal/domain"
	"github.com/redis/go-redis/v9"
)

const (
	dataKey = "currencies"

	dataTTL = 5 * time.Minute

	sendDirection    = "send"
	receiveDirection = "receive"
)

type CurrencyCache struct {
	rdb *redis.Client
}

func NewCurrencyCache(rdb *redis.Client) (*CurrencyCache, error) {
	const op = "redis.NewCurrencyCache"

	if rdb == nil {
		return nil, fmt.Errorf("%s: redis client is nil")
	}

	return &CurrencyCache{rdb: rdb}, nil
}

func (c *CurrencyCache) Set(ctx context.Context, direction string, value []domain.Currency) error {
	const op = "redis.CurrencyCache.Set"

	if direction != sendDirection && direction != receiveDirection {
		return fmt.Errorf("%s: invalid direction", op)
	}

	dataBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%s: failed to marshal data: %w", op, err)
	}

	err = c.rdb.Set(ctx, dataKey+":"+direction, dataBytes, dataTTL).Err()
	if err != nil {
		return fmt.Errorf("%s: failed to set data: %w", op, err)
	}

	return nil
}

func (c *CurrencyCache) Get(ctx context.Context, direction string) ([]domain.Currency, error) {
	const op = "redis.CurrencyCache.Get"

	if direction != sendDirection && direction != receiveDirection {
		return nil, fmt.Errorf("%s: invalid direction", op)
	}

	var data []domain.Currency

	res, err := c.rdb.Get(ctx, dataKey+":"+direction).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get data: %w", err)
	}

	err = json.Unmarshal([]byte(res), &data)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to extract data: %w", err)
	}

	return data, nil
}
