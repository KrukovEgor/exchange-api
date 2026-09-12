package domain

import (
	"context"
)

type CryptoProvider interface {
	GetCurrencies(ctx context.Context, dir Direction) ([]Currency, error)
}
