package etools

import (
	"fmt"
)

const (
	paramKeyPage       = "page"
	paramKeySafeSearch = "safeSearch" // Can be "true" or "false".

	paramCountry  = "country=web"
	paramLanguage = "language=all"
)

func safeSearchParamString(safesearch bool) string {
	if safesearch {
		return fmt.Sprintf("%v=%v", paramKeySafeSearch, "true")
	} else {
		return fmt.Sprintf("%v=%v", paramKeySafeSearch, "false")
	}
}
