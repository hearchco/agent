package mojeek

import (
	"fmt"
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

func localeParamString(locale options.Locale) string {
	spl := strings.SplitN(strings.ToLower(locale.String()), "_", 2)
	return fmt.Sprintf("%v=%v&%v=%v", params.Locale, spl[0], params.LocaleSec, spl[1])
}

func safeSearchParamString(safesearch bool) string {
	if safesearch {
		return fmt.Sprintf("%v=%v", params.SafeSearch, "1")
	} else {
		return fmt.Sprintf("%v=%v", params.SafeSearch, "0")
	}
}
