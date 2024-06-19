package brave

import (
	"fmt"
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

const (
	paramKeyPage       = "offset"
	paramKeyLocale     = "country"    // Should be last 2 characters of Locale.
	paramKeySafeSearch = "safesearch" // Can be "off" or "strict".

	paramSource     = "source=web"
	paramSpellcheck = "spellcheck=0"
)

func localeCookieString(locale options.Locale) string {
	region := strings.SplitN(strings.ToLower(locale.String()), "_", 2)[1]
	return fmt.Sprintf("%v=%v", paramKeyLocale, region)
}

func safeSearchCookieString(safesearch bool) string {
	if safesearch {
		return fmt.Sprintf("%v=%v", paramKeySafeSearch, "strict")
	} else {
		return fmt.Sprintf("%v=%v", paramKeySafeSearch, "off")
	}
}
