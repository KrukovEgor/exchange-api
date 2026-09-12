package easybit

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/KrukovEgor/exchange-api/internal/domain"
)

func mapTransportError(op string, err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, context.Canceled):
		return domain.NewInternalError(op, "request canceled by caller", err)
	case errors.Is(err, context.DeadlineExceeded):
		return domain.NewProviderUnavailableError(op, "provider request timed out", err)
	default:
		return domain.NewProviderUnavailableError(op, "provider is unavailable", err)
	}
}

func mapStatusError(op string, status int, body string) error {
	cause := fmt.Errorf("unexpected status %d: %s", status, truncate(body, readerLimitBytes))

	switch {
	case status == http.StatusTooManyRequests:
		return domain.NewProviderUnavailableError(op, "provider rate limit exceeded", cause)
	case status >= http.StatusInternalServerError:
		return domain.NewProviderUnavailableError(op, "provider server error", cause)
	case status == http.StatusUnauthorized, status == http.StatusForbidden:
		return domain.NewProviderRejectedError(op, "provider rejected credentials", cause)
	case status >= http.StatusBadRequest:
		return domain.NewInvalidInputError(op, "provider rejected the request", cause)
	default:
		return domain.NewProviderInvalidDataError(op, "unexpected provider status", cause)
	}
}

func mapAPIError(op string, apiErr apiError) error {
	if apiErr.ErrorCode == nil {
		return domain.NewProviderInvalidDataError(
			op, "provider returned error without an error code",
			fmt.Errorf("message: %q", apiErr.message()),
		)
	}

	code := apiErrorCode(*apiErr.ErrorCode)
	cause := fmt.Errorf("provider error %d: %s", code, apiErr.message())

	switch code {
	case codeRateLimitExceeded, codeOrdersUnavailable, codeCurrencySuspended, codeRateUnavailable:
		return domain.NewProviderUnavailableError(op, "provider is temporarily unavailable", cause)

	case codeBadRequest, codeUnsupportedCurrency, codeUnsupportedPair, codeUnsupportedNetwork,
		codeInvalidAddress, codeTagNotSupported, codeInvalidTag, codeAmountNotAllowed, codeInvalidVPM,
		codeRefundAddrRequired, codeInvalidOrderID, codeInvalidStatus, codeInvalidLimit, codeInvalidSortDirection:
		return domain.NewInvalidInputError(op, "provider rejected request parameters", cause)

	case codeInvalidExtraFee, codePartnersOnly, codeExtraFeeForbidden,
		codeVPMDisabled, codeExtraFeeOverrideDenied:
		return domain.NewProviderRejectedError(op, "provider account restriction", cause)
	}

	switch {
	case code >= transportRangeMin && code < transportRangeMax:
		return domain.NewProviderUnavailableError(op, "provider transport error", cause)
	case code >= currencyRangeMin && code < currencyRangeMax,
		code >= orderRangeMin && code < orderRangeMax:
		return domain.NewInvalidInputError(op, "provider rejected request parameters", cause)
	case code >= accountRangeMin && code < accountRangeMax:
		return domain.NewProviderRejectedError(op, "provider account restriction", cause)
	default:
		return domain.NewProviderInvalidDataError(op, "unknown provider error code", cause)
	}
}

func truncate(s string, limit int) string {
	if len(s) <= limit {
		return s
	}

	return s[:limit] + "..."
}
