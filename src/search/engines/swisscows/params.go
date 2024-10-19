package swisscows

import (
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

const (
	// Variable params.
	paramQueryK  = "query"
	paramPageK   = "offset"
	paramLocaleK = "region" // Should be the same as Locale, only with "_" replaced by "-".

	// Constant params.
	paramFreshnessK, paramFreshnessV = "freshness", "All"
	paramItemsK, paramItemsV         = "itemsCount", "10"
)

func localeParamValue(locale options.Locale) string {
	return strings.Replace(locale.String(), "_", "-", 1)
}
