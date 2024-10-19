package google

import (
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

const (
	// Variable params.
	paramQueryK      = "q"
	paramPageK       = "start"
	paramLocaleK     = "hl"   // Should be first 2 characters of Locale.
	paramLocaleSecK  = "lr"   // Should be first 2 characters of Locale with prefixed "lang_".
	paramSafeSearchK = "safe" // Can be "off", "medium or "high".

	// Constant params.
	paramFilterK, paramFilterV = "filter", "0"

	// Image search variable params.
	imgParamPageK, imgParamPageVPrefix = "async", "_fmt:json,p:1,ijn:"

	// Image search constant params.
	imgParamTbmK, imgParamTbmV         = "tbm", "isch"
	imgParamAsearchK, imgParamAsearchV = "asearch", "isch"

	// Suggestions constant params.
	sugParamClientK, sugParamClientV = "client", "firefox"
)

func localeParamValues(locale options.Locale) (string, string) {
	lang := strings.SplitN(strings.ToLower(locale.String()), "_", 2)[0]
	return lang, "lang_" + lang
}

func safeSearchParamValue(safesearch bool) string {
	if safesearch {
		return "high"
	} else {
		return "off"
	}
}
