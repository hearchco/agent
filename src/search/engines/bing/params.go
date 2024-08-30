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

	// Image search params.
	imgParamKeyLocale       = "m"
	imgParamKeyLocaleAlt    = "mkt"
	imgParamKeyLocaleSec    = "u"
	imgParamKeyLocaleSecAlt = "ui"
	imgParamAsync           = "async=1"
	imgParamCount           = "count=35"
)

func localeParamString(locale options.Locale) string {
	spl := strings.SplitN(strings.ToLower(locale.String()), "_", 2)
	return fmt.Sprintf("%v=%v&%v=%v", paramKeyLocale, spl[0], paramKeyLocaleSec, spl[1])
}

func localeCookieString(locale options.Locale) string {
	spl := strings.SplitN(strings.ToLower(locale.String()), "_", 2)
	return fmt.Sprintf("%v=%v&%v=%v", imgParamKeyLocale, spl[1], imgParamKeyLocaleSec, spl[0])
}

func localeAltCookieString(locale options.Locale) string {
	spl := strings.SplitN(strings.ToLower(locale.String()), "_", 2)
	return fmt.Sprintf("%v=%v&%v=%v", imgParamKeyLocaleAlt, spl[1], imgParamKeyLocaleSecAlt, spl[0])
}
