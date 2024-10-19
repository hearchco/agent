package bing

import (
	"fmt"
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

const (
	// Variables params.
	paramQueryK     = "q"
	paramPageK      = "first"
	paramLocaleK    = "setlang" // Should be first 2 characters of Locale.
	paramLocaleSecK = "cc"      // Should be last 2 characters of Locale.
	// paramSafeSearchK = ""        // Always enabled.

	// Image variable params.
	imgCookieLocaleK       = "m"
	imgCookieLocaleSecK    = "u"
	imgCookieLocaleAltK    = "mkt"
	imgCookieLocaleAltSecK = "ui"

	// Image constant params.
	imgParamAsyncK, imgParamAsyncV = "async", "1"
	imgParamCountK, imgParamCountV = "count", "35"
)

func localeParamValues(locale options.Locale) (string, string) {
	spl := strings.SplitN(strings.ToLower(locale.String()), "_", 2)
	return spl[0], spl[1]
}

func localeCookieString(locale options.Locale) string {
	spl := strings.SplitN(strings.ToLower(locale.String()), "_", 2)
	return fmt.Sprintf("%v=%v&%v=%v", imgCookieLocaleK, spl[1], imgCookieLocaleSecK, spl[0])
}

func localeAltCookieString(locale options.Locale) string {
	spl := strings.SplitN(strings.ToLower(locale.String()), "_", 2)
	return fmt.Sprintf("%v=%v&%v=%v", imgCookieLocaleAltK, spl[1], imgCookieLocaleAltSecK, spl[0])
}
