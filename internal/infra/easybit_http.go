package infra

import (
	"fmt"
	"io"
	"net/http"

	"github.com/KrukovEgor/exchange-api/internal/config"
	"github.com/cenkalti/backoff/v7"
	"golang.org/x/time/rate"
)

type retryTransport struct {
	maxRetries        int
	disableKeepAlives bool
	next              http.RoundTripper
}

func (t *retryTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return backoff.Retry(r.Context(), func() (*http.Response, error) {
		reqClone := r.Clone(r.Context())

		if r.Body != nil {
			if r.GetBody == nil {
				return nil, backoff.Permanent(fmt.Errorf("request body getter is not init"))
			}
			body, err := r.GetBody()
			if err != nil {
				return nil, backoff.Permanent(fmt.Errorf("failed to retrieve request body: %w", err))
			}
			reqClone.Body = body
		}

		resp, err := t.next.RoundTrip(reqClone)
		if err != nil {
			return nil, fmt.Errorf("failed to execute next round trip: %w", err)
		}

		if resp.StatusCode == http.StatusTooManyRequests ||
			resp.StatusCode >= http.StatusInternalServerError {
			if !t.disableKeepAlives {
				_, _ = io.Copy(io.Discard, resp.Body)
			}
			resp.Body.Close()
			return nil, fmt.Errorf("response status: %s", resp.Status)
		}

		return resp, nil
	}, backoff.WithMaxTries(uint(t.maxRetries)))
}

type rateLimitedTransport struct {
	rateLimiter *rate.Limiter
	next        http.RoundTripper
}

func (t *rateLimitedTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if err := t.rateLimiter.Wait(r.Context()); err != nil {
		return nil, err
	}

	return t.next.RoundTrip(r)
}

type authTransport struct {
	apiKey string
	next   http.RoundTripper
}

func (t *authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("API-KEY", t.apiKey)
	return t.next.RoundTrip(r)
}

func chainRoundTripper(rt http.RoundTripper, middlewares ...func(http.RoundTripper) http.RoundTripper) http.RoundTripper {
	for _, m := range middlewares {
		rt = m(rt)
	}
	return rt
}

func NewEasyBitHTTPClient(cfg *config.EasyBitConfig) (*http.Client, error) {
	const op = "infra.NewEasyBitHTTPClient"

	if cfg == nil {
		return nil, fmt.Errorf("%s: config is nil", op)
	}

	limiter := rate.NewLimiter(rate.Limit(cfg.LimiterRate), cfg.LimiterBurst)

	baseTransport := http.DefaultTransport.(*http.Transport).Clone()
	baseTransport.DisableKeepAlives = cfg.DisableKeepAlives

	transportChain := chainRoundTripper(
		baseTransport,
		func(rt http.RoundTripper) http.RoundTripper {
			return &authTransport{
				apiKey: cfg.APIKey,
				next:   rt,
			}
		},
		func(rt http.RoundTripper) http.RoundTripper {
			return &rateLimitedTransport{
				rateLimiter: limiter,
				next:        rt,
			}
		},
		func(rt http.RoundTripper) http.RoundTripper {
			return &retryTransport{
				maxRetries: cfg.MaxRetries,
				next:       rt,
			}
		},
	)

	return &http.Client{
		Transport: transportChain,
		Timeout:   cfg.RequestTimeout,
	}, nil
}
