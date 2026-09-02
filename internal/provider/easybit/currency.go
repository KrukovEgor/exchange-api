package easybit

import (
	"context"
	"fmt"
	"time"

	"github.com/KrukovEgor/exchange-api/internal/domain"
)

const attemptTimeout = 2 * time.Second

func (c *EasyBitClient) GetCurrencies(ctx context.Context, direction string) ([]domain.Currency, error) {
	var rawData apiResponse[[]currency]

	err := c.doRequest(ctx, "GET", "/currencyList", attemptTimeout, &rawData)
	if err != nil {
		return nil, fmt.Errorf("failed to get available currencies: %w", err)
	}

	return mapCurrencies(*rawData.Data, direction), nil
}
