package domain

type (
	Currency struct {
		CurrencyTicker string
		CurrencyName   string
		NetworkList    []Network
	}

	Network struct {
		NetworkTicker string
		NetworkName   string
	}
)
