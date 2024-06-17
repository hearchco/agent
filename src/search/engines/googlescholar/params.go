package googlescholar

import (
	"fmt"
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

func localeParamString(locale options.Locale) string {
	lang := strings.SplitN(strings.ToLower(locale.String()), "_", 2)[0]
	return fmt.Sprintf("%v=%v&%v=lang_%v", params.Locale, lang, params.LocaleSec, lang)
}

func safeSearchParamString(safesearch bool) string {
	if safesearch {
		return fmt.Sprintf("%v=%v", params.SafeSearch, "high")
	} else {
		return fmt.Sprintf("%v=%v", params.SafeSearch, "off")
	}
}
