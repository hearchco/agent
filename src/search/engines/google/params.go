package google

import (
	"fmt"
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

const (
	paramKeyPage       = "start"
	paramKeyLocale     = "hl"   // Should be first 2 characters of Locale.
	paramKeyLocaleSec  = "lr"   // Should be first 2 characters of Locale with prefixed "lang_".
	paramKeySafeSearch = "safe" // Can be "off", "medium or "high".

	paramFilter = "filter=0"

	// Suggestions API params.
	sugParamClient = "client=firefox"
)

func localeParamString(locale options.Locale) string {
	lang := strings.SplitN(strings.ToLower(locale.String()), "_", 2)[0]
	return fmt.Sprintf("%v=%v&%v=lang_%v", paramKeyLocale, lang, paramKeyLocaleSec, lang)
}

func safeSearchParamString(safesearch bool) string {
	if safesearch {
		return fmt.Sprintf("%v=%v", paramKeySafeSearch, "high")
	} else {
		return fmt.Sprintf("%v=%v", paramKeySafeSearch, "off")
	}
}
