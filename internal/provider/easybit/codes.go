package easybit

type apiErrorCode int

const (
	transportRangeMin apiErrorCode = 400
	transportRangeMax apiErrorCode = 1000

	currencyRangeMin apiErrorCode = 1000
	currencyRangeMax apiErrorCode = 2000

	orderRangeMin apiErrorCode = 2000
	orderRangeMax apiErrorCode = 3000

	accountRangeMin apiErrorCode = 3000
	accountRangeMax apiErrorCode = 4000
)

const (
	codeBadRequest        apiErrorCode = 400
	codeRateLimitExceeded apiErrorCode = 429
)

const (
	codeOrdersUnavailable   apiErrorCode = 1000
	codeUnsupportedCurrency apiErrorCode = 1001
	codeUnsupportedPair     apiErrorCode = 1002
	codeUnsupportedNetwork  apiErrorCode = 1003
	codeCurrencySuspended   apiErrorCode = 1004
	codeInvalidAddress      apiErrorCode = 1005
	codeTagNotSupported     apiErrorCode = 1006
	codeInvalidTag          apiErrorCode = 1007
	codeAmountNotAllowed    apiErrorCode = 1008
	codeInvalidVPM          apiErrorCode = 1009
	codeRefundAddrRequired  apiErrorCode = 1010
	codeRateUnavailable     apiErrorCode = 1011
)

const (
	codeInvalidOrderID       apiErrorCode = 2001
	codeInvalidStatus        apiErrorCode = 2002
	codeInvalidLimit         apiErrorCode = 2003
	codeInvalidSortDirection apiErrorCode = 2004
)

const (
	codeInvalidExtraFee        apiErrorCode = 3001
	codePartnersOnly           apiErrorCode = 3002
	codeExtraFeeForbidden      apiErrorCode = 3003
	codeVPMDisabled            apiErrorCode = 3004
	codeExtraFeeOverrideDenied apiErrorCode = 3005
)
