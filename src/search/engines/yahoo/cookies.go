package yahoo

import (
	"fmt"
)

func safeSearchCookieString(safesearch bool) string {
	if safesearch {
		return fmt.Sprintf("%v=%v", params.SafeSearch, "r")
	} else {
		return fmt.Sprintf("%v=%v", params.SafeSearch, "p")
	}
}
