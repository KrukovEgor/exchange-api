package easybit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v7"
)

const defaultEasyBitURL = "https://api.easybit.com"
const maxRetriesNumber = 3

type EasyBitClient struct {
	HTTPClient *http.Client
	BaseURL    string
}

func New(httpClient *http.Client, baseURL string) (*EasyBitClient, error) {
	const op = "easybit.New"

	if httpClient == nil {
		return nil, fmt.Errorf("%s: HTTP client is nil", op)
	}

	if baseURL == "" {
		baseURL = defaultEasyBitURL
	}

	return &EasyBitClient{
		HTTPClient: httpClient,
		BaseURL:    baseURL,
	}, nil
}

func (c *EasyBitClient) doRequest(
	baseCtx context.Context,
	method string,
	endpoint string,
	attemptTimeout time.Duration,
	out any,
) error {
	const op = "easybit.EasyBitClient.doRequest"

	_, err := backoff.Retry(baseCtx, func() (struct{}, error) {
		attemptCtx, cancel := context.WithTimeout(baseCtx, attemptTimeout)
		defer cancel()

		req, err := http.NewRequestWithContext(attemptCtx, method, c.BaseURL+endpoint, nil)
		if err != nil {
			return struct{}{}, backoff.Permanent(fmt.Errorf("failed to create request: %w", err))
		}

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return struct{}{}, backoff.Permanent(fmt.Errorf("failed to send request: %w", err))
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests ||
			resp.StatusCode >= http.StatusInternalServerError {
			return struct{}{}, fmt.Errorf("response status: %s", resp.Status)
		}

		err = json.NewDecoder(resp.Body).Decode(out)
		if err != nil {
			return struct{}{}, fmt.Errorf("%s: failed to decode response: %w", op, err)
		}

		return struct{}{}, nil
	}, backoff.WithMaxTries(maxRetriesNumber))

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
