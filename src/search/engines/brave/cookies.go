package brave

import (
	"fmt"
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

func localeCookieString(locale options.Locale) string {
	region := strings.SplitN(strings.ToLower(locale.String()), "_", 2)[1]
	return fmt.Sprintf("%v=%v", params.Locale, region)
}

func safeSearchCookieString(safesearch bool) string {
	if safesearch {
		return fmt.Sprintf("%v=%v", params.SafeSearch, "strict")
	} else {
		return fmt.Sprintf("%v=%v", params.SafeSearch, "off")
	}
}
