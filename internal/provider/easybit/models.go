package easybit

type (
	apiResponse[T any] struct {
		Success int `json:"success"`
		Data    *T  `json:"data,omitempty"`
		apiError
	}

	apiError struct {
		ErrorMessage *string `json:"errorMessage,omitempty"`
		ErrorCode    *int    `json:"errorCode,omitempty"`
	}

	network struct {
		Network              string  `json:"network"`
		Name                 string  `json:"name"`
		IsDefault            bool    `json:"isDefault"`
		SendStatus           bool    `json:"sendStatus"`
		ReceiveStatus        bool    `json:"receiveStatus"`
		ReceiveDecimals      int     `json:"receiveDecimals"`
		ConfirmationsMinimum int     `json:"confirmationsMinimum"`
		ConfirmationsMaximum int     `json:"confirmationsMaximum"`
		Explorer             string  `json:"explorer"`
		ExplorerHash         string  `json:"explorerHash"`
		ExplorerAddress      string  `json:"explorerAddress"`
		HasTag               bool    `json:"hasTag"`
		TagName              *string `json:"tagName,omitempty"`
		ContractAddress      *string `json:"contractAddress,omitempty"`
		ExplorerContract     *string `json:"explorerContract,omitempty"`
	}

	currency struct {
		Currency         string    `json:"currency"`
		Name             string    `json:"name"`
		SendStatusAll    bool      `json:"sendStatusAll"`
		ReceiveStatusAll bool      `json:"receiveStatusAll"`
		NetworkList      []network `json:"networkList"`
	}
)
