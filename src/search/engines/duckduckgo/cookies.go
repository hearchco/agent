package duckduckgo

import (
	"fmt"
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

func localeCookieString(locale options.Locale) string {
	spl := strings.SplitN(strings.ToLower(locale.String()), "_", 2)
	return fmt.Sprintf("%v=%v-%v", params.Locale, spl[1], spl[0])
}
