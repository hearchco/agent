package mojeek

import (
	"fmt"
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

const (
	paramKeyPage       = "s"
	paramKeyLocale     = "lb"   // Should be first 2 characters of Locale.
	paramKeyLocaleSec  = "arc"  // Should be last 2 characters of Locale.
	paramKeySafeSearch = "safe" // Can be "0" or "1".
)

func localeParamString(locale options.Locale) string {
	spl := strings.SplitN(strings.ToLower(locale.String()), "_", 2)
	return fmt.Sprintf("%v=%v&%v=%v", paramKeyLocale, spl[0], paramKeyLocaleSec, spl[1])
}

func safeSearchParamString(safesearch bool) string {
	if safesearch {
		return fmt.Sprintf("%v=%v", paramKeySafeSearch, "1")
	} else {
		return fmt.Sprintf("%v=%v", paramKeySafeSearch, "0")
	}
}
