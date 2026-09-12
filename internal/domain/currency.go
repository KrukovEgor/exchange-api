package domain

import (
	"context"
)

type (
	Currency struct {
		CurrencyTicker string    `json:"currency_ticker"`
		CurrencyName   string    `json:"currency_name"`
		NetworkList    []Network `json:"network_list"`
	}

	Network struct {
		NetworkTicker string `json:"network_ticker"`
		NetworkName   string `json:"network_name"`
	}
)

type CurrencyCache interface {
	Set(ctx context.Context, direction string, value []Currency) error
	Get(ctx context.Context, direction string) ([]Currency, error)
}
