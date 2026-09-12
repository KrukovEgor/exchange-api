package infra

import (
	"net/http"

	"github.com/KrukovEgor/exchange-api/internal/config"
	"github.com/KrukovEgor/exchange-api/internal/domain"
	"golang.org/x/time/rate"
)

type authTransport struct {
	apiKey string
	next   http.RoundTripper
}

func (t *authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	reqClone := r.Clone(r.Context())
	reqClone.Header.Set("API-KEY", t.apiKey)

	return t.next.RoundTrip(reqClone)
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

func NewEasyBitHTTPClient(cfg *config.EasyBitConfig) (*http.Client, error) {
	const op = "infra.NewEasyBitHTTPClient"

	if cfg == nil {
		return nil, domain.NewInvalidInputError(op, "config is nil", nil)
	}

	if cfg.APIKey == "" {
		return nil, domain.NewInvalidInputError(op, "api key is empty", nil)
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = cfg.MaxIdleConnsPerHost

	var rt http.RoundTripper = transport
	rt = &authTransport{apiKey: cfg.APIKey, next: rt}
	rt = &rateLimitedTransport{
		rateLimiter: rate.NewLimiter(rate.Limit(cfg.LimiterRate), cfg.LimiterBurst),
		next:        rt,
	}

	return &http.Client{
		Transport: rt,
		Timeout:   cfg.RequestTimeout,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}, nil
}
