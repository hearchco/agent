package duckduckgo

import (
	"fmt"
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

const (
	paramKeyPage       = "dc"
	paramKeyLocale     = "kl" // Should be Locale with _ replaced by - and first 2 letters as last and vice versa.
	paramKeySafeSearch = ""   // Always enabled.
)

func localeCookieString(locale options.Locale) string {
	spl := strings.SplitN(strings.ToLower(locale.String()), "_", 2)
	return fmt.Sprintf("%v=%v-%v", paramKeyLocale, spl[1], spl[0])
}
