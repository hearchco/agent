package duckduckgo

import (
	"fmt"
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

const (
	// Variable params.
	paramQueryK   = "q"
	paramPageK    = "dc"
	cookieLocaleK = "kl" // Should be Locale with _ replaced by - and first 2 letters as last and vice versa.
	// paramSafeSearchK = ""   // Always enabled.

	// Suggestions variable params.
	sugParamTypeK, sugParamTypeV = "type", "list"
)

func localeCookieString(locale options.Locale) string {
	spl := strings.SplitN(strings.ToLower(locale.String()), "_", 2)
	return fmt.Sprintf("%v=%v-%v", cookieLocaleK, spl[1], spl[0])
}
