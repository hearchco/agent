package yahoo

import (
	"fmt"
)

const (
	paramKeyPage       = "b"
	paramKeySafeSearch = "vm" // Can be "p" (disabled) or "r" (enabled).

	paramSafeSearchPrefix = "sB=v=1&pn=10&rw=new&userset=0"
)

func safeSearchCookieString(safesearch bool) string {
	if safesearch {
		return fmt.Sprintf("%v&%v=%v", paramSafeSearchPrefix, paramKeySafeSearch, "r")
	} else {
		return fmt.Sprintf("%v&%v=%v", paramSafeSearchPrefix, paramKeySafeSearch, "p")
	}
}
