package presearch

import (
	"fmt"
)

func safeSearchCookieString(safesearch bool) string {
	if safesearch {
		return fmt.Sprintf("%v=%v", params.SafeSearch, "true")
	} else {
		return fmt.Sprintf("%v=%v", params.SafeSearch, "false")
	}
}
