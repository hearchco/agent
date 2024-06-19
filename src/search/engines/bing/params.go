package bing

import (
	"fmt"
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

const (
	paramKeyPage       = "first"
	paramKeyLocale     = "setlang" // Should be first 2 characters of Locale.
	paramKeyLocaleSec  = "cc"      // Should be last 2 characters of Locale.
	paramKeySafeSearch = ""        // Always enabled.
)

func localeParamString(locale options.Locale) string {
	spl := strings.SplitN(strings.ToLower(locale.String()), "_", 2)
	return fmt.Sprintf("%v=%v&%v=%v", paramKeyLocale, spl[0], paramKeyLocaleSec, spl[1])
}
