package presearch

import (
	"fmt"
)

const (
	paramKeyPage       = "page"
	paramKeySafeSearch = "use_safe_search" // // Can be "true" or "false".
)

func safeSearchCookieString(safesearch bool) string {
	if safesearch {
		return fmt.Sprintf("%v=%v", paramKeySafeSearch, "true")
	} else {
		return fmt.Sprintf("%v=%v", paramKeySafeSearch, "false")
	}
}
