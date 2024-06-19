package startpage

import (
	"fmt"
)

const (
	paramKeyPage       = "page"
	paramKeySafeSearch = "qadf" // Can be "none" or empty param (empty means it's enabled).
)

func safeSearchParamString(safesearch bool) string {
	if safesearch {
		return ""
	} else {
		return fmt.Sprintf("%v=%v", paramKeySafeSearch, "none")
	}
}
