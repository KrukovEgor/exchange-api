package infra

import (
	"fmt"
	"net/http"

	"github.com/KrukovEgor/exchange-api/internal/config"
	"golang.org/x/time/rate"
)

type rateLimitedTransport struct {
	rateLimiter *rate.Limiter
	next        http.RoundTripper
}

func (t *rateLimitedTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if err := t.rateLimiter.Wait(r.Context()); err != nil {
		return nil, err
	}

	reqClone := r.Clone(r.Context())

	if r.Body != nil {
		if r.GetBody == nil {
			return nil, fmt.Errorf("request body getter is not init")
		}
		body, err := r.GetBody()
		if err != nil {
			return nil, fmt.Errorf("request body getter is not init")
		}
		reqClone.Body = body
	}

	return t.next.RoundTrip(reqClone)
}

type authTransport struct {
	apiKey string
	next   http.RoundTripper
}

func (t *authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("API-KEY", t.apiKey)
	return t.next.RoundTrip(r)
}

func chainRoundTripper(
	rt http.RoundTripper,
	middlewares ...func(http.RoundTripper) http.RoundTripper,
) http.RoundTripper {
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
	baseTransport.MaxIdleConnsPerHost = cfg.MaxIdleConnsPerHost

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
	)

	return &http.Client{
		Transport: transportChain,
		Timeout:   cfg.RequestTimeout,
	}, nil
}
