package easybit

import (
	"github.com/KrukovEgor/exchange-api/internal/domain"
)

func mapCurrency(dto currency, dir domain.Direction) domain.Currency {
	return domain.Currency{
		CurrencyTicker: dto.Currency,
		CurrencyName:   dto.Name,
		NetworkList:    mapNetworks(dto.NetworkList, dir),
	}
}

func mapCurrencies(dtos []currency, dir domain.Direction) []domain.Currency {
	domainCurrencies := make([]domain.Currency, 0, len(dtos))

	for _, dto := range dtos {
		if dto.available(dir) {
			mappedCurrency := mapCurrency(dto, dir)
			if len(mappedCurrency.NetworkList) > 0 {
				domainCurrencies = append(domainCurrencies, mappedCurrency)
			}
		}
	}

	return domainCurrencies
}

func mapNetwork(dto network) domain.Network {
	return domain.Network{
		NetworkTicker: dto.Network,
		NetworkName:   dto.Name,
	}
}

func mapNetworks(dtos []network, dir domain.Direction) []domain.Network {
	domainNetworks := make([]domain.Network, 0, len(dtos))

	for _, dto := range dtos {
		if dto.available(dir) {
			domainNetworks = append(domainNetworks, mapNetwork(dto))
		}
	}

	return domainNetworks
}

func (c currency) available(dir domain.Direction) bool {
	switch dir {
	case domain.SendDirection:
		return c.SendStatusAll
	case domain.ReceiveDirection:
		return c.ReceiveStatusAll
	default:
		return false
	}
}

func (n network) available(dir domain.Direction) bool {
	switch dir {
	case domain.SendDirection:
		return n.SendStatus
	case domain.ReceiveDirection:
		return n.ReceiveStatus
	default:
		return false
	}
}
