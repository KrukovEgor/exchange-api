package easybit

import (
	"github.com/KrukovEgor/exchange-api/internal/domain"
)

const (
	sendDirection    = "send"
	receiveDirection = "receive"
)

func mapCurrency(dto currency, direction string) domain.Currency {
	return domain.Currency{
		CurrencyTicker: dto.Currency,
		CurrencyName:   dto.Name,
		NetworkList:    mapNetworks(dto.NetworkList, direction),
	}
}

func mapCurrencies(dtos []currency, direction string) []domain.Currency {
	domainCurrencies := make([]domain.Currency, 0, len(dtos))

	for _, dto := range dtos {
		if direction == sendDirection && dto.SendStatusAll ||
			direction == receiveDirection && dto.ReceiveStatusAll {
			mappedCurrency := mapCurrency(dto, direction)
			if len(mappedCurrency.NetworkList) > 0 {
				domainCurrencies = append(domainCurrencies, mapCurrency(dto, direction))
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

func mapNetworks(dtos []network, direction string) []domain.Network {
	domainNetworks := make([]domain.Network, 0, len(dtos))

	for _, dto := range dtos {
		if direction == sendDirection && dto.SendStatus ||
			direction == receiveDirection && dto.ReceiveStatus {
			domainNetworks = append(domainNetworks, mapNetwork(dto))
		}
	}

	return domainNetworks
}
