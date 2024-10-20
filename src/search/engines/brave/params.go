package brave

import (
	"fmt"
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

const (
	// Variable params.
	paramQueryK       = "q"
	paramPageK        = "offset"
	cookieLocaleK     = "country"    // Should be last 2 characters of Locale.
	cookieSafeSearchK = "safesearch" // Can be "off" or "strict".

	// Constant params.
	paramSourceK, paramSourceV         = "source", "web"
	paramSpellcheckK, paramSpellcheckV = "spellcheck", "0"
)

func localeCookieString(locale options.Locale) string {
	region := strings.SplitN(strings.ToLower(locale.String()), "_", 2)[1]
	return fmt.Sprintf("%v=%v", cookieLocaleK, region)
}

func safeSearchCookieString(safesearch bool) string {
	if safesearch {
		return fmt.Sprintf("%v=%v", cookieSafeSearchK, "strict")
	} else {
		return fmt.Sprintf("%v=%v", cookieSafeSearchK, "off")
	}
}
