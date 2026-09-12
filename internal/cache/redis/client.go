package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func set[T any](c *redis.Client, parentCtx context.Context, key string, value T, ttl time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data")
	}

	err = c.Set(parentCtx, key, bytes, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func get[T any](c *redis.Client, parentCtx context.Context, key string, out T) error {
	bytes, err := c.Get(parentCtx, key).Bytes()
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, out)
	if err != nil {
		return err
	}

	return nil
}
