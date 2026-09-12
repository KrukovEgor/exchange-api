package easybit

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/KrukovEgor/exchange-api/internal/config"
	"github.com/KrukovEgor/exchange-api/internal/domain"
	"github.com/cenkalti/backoff/v7"
)

const (
	defaultEasyBitURL = "https://api.easybit.com"

	readerLimitBytes = 2 << 20
	attemptBuffer    = 50 * time.Millisecond
)

type Client struct {
	httpClient     *http.Client
	baseURL        *url.URL
	attemptTimeout time.Duration
	maxRetries     uint
	newBackOff     func() backoff.BackOff
}

func New(httpClient *http.Client, cfg *config.EasyBitConfig) (*Client, error) {
	const op = "easybit.New"

	if httpClient == nil {
		return nil, domain.NewInternalError(op, "HTTP client is nil", nil)
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultEasyBitURL
	}

	parsedURL, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, domain.NewInternalError(op, "invalid base URL", nil)
	}

	return &Client{
		httpClient:     httpClient,
		baseURL:        parsedURL,
		attemptTimeout: cfg.AttemptTimeout,
		maxRetries:     cfg.MaxRetries,
		newBackOff:     func() backoff.BackOff { return backoff.NewExponentialBackOff() },
	}, nil
}

func idempotent(method string) bool {
	switch method {
	case http.MethodGet:
		return true
	default:
		return false
	}
}

func doRequest[T any](parentCtx context.Context, client *Client, method, endpoint string, query url.Values, out T) error {
	const op = "easybit.doRequest"

	target := client.baseURL.JoinPath(endpoint)
	if len(query) > 0 {
		target.RawQuery = query.Encode()
	}

	maxRetries := client.maxRetries
	if !idempotent(method) {
		maxRetries = 1
	}

	bodyBytes, err := backoff.Retry(
		parentCtx,
		func() ([]byte, error) {
			return attempt(parentCtx, client, op, method, target.String())
		},
		backoff.WithBackOff(client.newBackOff()),
		backoff.WithMaxElapsedTime(0),
		backoff.WithMaxTries(maxRetries),
	)
	if err != nil {
		if retryErr := backoff.AsRetryError(err); retryErr != nil && retryErr.LastErr != nil {
			return retryErr.LastErr
		}

		return mapTransportError(op, err)
	}

	err = json.Unmarshal(bodyBytes, out)
	if err != nil {
		return domain.NewProviderInvalidDataError(
			op, "failder to unmarshal provider response", err,
		)
	}

	return nil
}

func attempt(parentCtx context.Context, client *Client, op, method, target string) ([]byte, error) {
	attemptCtx, cancel := context.WithTimeout(parentCtx, client.attemptTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(attemptCtx, method, target, nil)
	if err != nil {
		return nil, backoff.Permanent(domain.NewInternalError(op, "failed to build request", err))
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) && parentCtx.Err() != nil {
			return nil, backoff.Permanent(mapTransportError(op, err))
		}

		return nil, mapTransportError(op, err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, readerLimitBytes))
	if err != nil {
		return nil, mapTransportError(op, err)
	}

	switch {
	case resp.StatusCode == http.StatusTooManyRequests:
		statusErr := mapStatusError(op, resp.StatusCode, string(bodyBytes))
		if duration, ok := checkRetryAfterHeader(resp.Header); ok {
			if !fitsInCtxDeadline(parentCtx, duration) {
				return nil, backoff.Permanent(mapTransportError(op, err))
			}

			return nil, backoff.RetryAfter(duration, statusErr)
		}

		return nil, statusErr
	case resp.StatusCode >= http.StatusInternalServerError:
		return nil, mapStatusError(op, resp.StatusCode, string(bodyBytes))
	case resp.StatusCode >= http.StatusBadRequest:
		return nil, backoff.Permanent(mapStatusError(op, resp.StatusCode, string(bodyBytes)))
	case resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices:
		return nil, backoff.Permanent(mapStatusError(op, resp.StatusCode, string(bodyBytes)))
	}

	return bodyBytes, nil
}

func checkRetryAfterHeader(h http.Header) (time.Duration, bool) {
	data := h.Get("Retry-After")
	if data == "" {
		return 0, false
	}

	seconds, err := strconv.Atoi(data)
	if err == nil {
		return time.Duration(seconds) * time.Second, true
	}

	date, err := http.ParseTime(data)
	if err == nil {
		return time.Until(date), true
	}

	return 0, false
}

func fitsInCtxDeadline(ctx context.Context, d time.Duration) bool {
	ctxDeadline, ok := ctx.Deadline()
	if !ok {
		return true
	}

	return time.Until(ctxDeadline) > d+attemptBuffer
}
