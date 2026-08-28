package easybit

import (
	"fmt"
	"net/http"
)

const defaultEasyBitURL = "https://api.easybit.com"

type EasyBitClient struct {
	HTTPClient *http.Client
	BaseURL    string
}

func New(baseURL string, client *http.Client) (*EasyBitClient, error) {
	const op = "easybit.New"

	if client == nil {
		return nil, fmt.Errorf("%s: HTTP client is nil", op)
	}

	if baseURL == "" {
		baseURL = defaultEasyBitURL
	}

	return &EasyBitClient{
		HTTPClient: client,
		BaseURL:    baseURL,
	}, nil
}
