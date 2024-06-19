package googleimages

import (
	"fmt"
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

const (
	paramKeyPage       = "async=_fmt:json,p:1,ijn"
	paramKeyLocale     = "hl"   // Should be first 2 characters of Locale.
	paramKeyLocaleSec  = "lr"   // Should be first 2 characters of Locale with prefixed "lang_".
	paramKeySafeSearch = "safe" // Can be "off", "medium or "high".

	paramTbm     = "tbm=isch"
	paramAsearch = "asearch=isch"
	paramFilter  = "filter=0"
)

func localeParamString(locale options.Locale) string {
	lang := strings.SplitN(strings.ToLower(locale.String()), "_", 2)[0]
	return fmt.Sprintf("%v=%v&%v=lang_%v", paramKeyLocale, lang, paramKeyLocaleSec, lang)
}

func safeSearchParamString(safesearch bool) string {
	if safesearch {
		return fmt.Sprintf("%v=%v", paramKeySafeSearch, "high")
	} else {
		return fmt.Sprintf("%v=%v", paramKeySafeSearch, "off")
	}
}
