package easybit

import (
	"context"
	"net/http"

	"github.com/KrukovEgor/exchange-api/internal/domain"
)

const (
	currencyEndpoint = "/currencyList"
)

func (c *Client) GetCurrencies(parentCtx context.Context, dir domain.Direction) ([]domain.Currency, error) {
	const op = "easybit.Client.GetCurrencies"

	if !dir.Valid() {
		return nil, domain.NewInvalidInputError(op, "invalid direction", nil)
	}

	var rawData apiResponse[[]currency]

	err := doRequest(parentCtx, c, http.MethodGet, currencyEndpoint, map[string][]string{}, &rawData)

	if err != nil {
		return nil, err
	}

	if rawData.failed() {
		return nil, mapAPIError(op, rawData.apiError)
	}

	if rawData.Data == nil {
		return nil, domain.NewProviderInvalidDataError(op, "provider returned successful response without data", nil)
	}

	return mapCurrencies(*rawData.Data, dir), nil
}
