package presearch

import (
	"fmt"
)

const (
	// Variable params.
	paramQueryK       = "q"
	paramPageK        = "page"
	cookieSafeSearchK = "use_safe_search" // Can be "true" or "false".
)

func safeSearchCookieString(safesearch bool) string {
	if safesearch {
		return fmt.Sprintf("%v=%v", cookieSafeSearchK, "true")
	} else {
		return fmt.Sprintf("%v=%v", cookieSafeSearchK, "false")
	}
}
