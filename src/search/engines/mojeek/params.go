package mojeek

import (
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

const (
	// Variable params.
	paramQueryK      = "q"
	paramPageK       = "s"
	paramLocaleK     = "lb"   // Should be first 2 characters of Locale.
	paramLocaleSecK  = "arc"  // Should be last 2 characters of Locale.
	paramSafeSearchK = "safe" // Can be "0" or "1".
)

func localeParamValues(locale options.Locale) (string, string) {
	spl := strings.SplitN(strings.ToLower(locale.String()), "_", 2)
	return spl[0], spl[1]
}

func safeSearchParamValue(safesearch bool) string {
	if safesearch {
		return "1"
	} else {
		return "0"
	}
}
