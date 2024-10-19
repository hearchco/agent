package yahoo

import (
	"fmt"
)

const (
	// Variable params.
	paramQueryK       = "p"
	paramPageK        = "b"
	cookieSafeSearchK = "vm" // Can be "p" (disabled) or "r" (enabled).

	// Constant params.
	cookieSafeSearchPrefix = "sB=v=1&pn=10&rw=new&userset=0"
	// paramSbK, paramSbV           = "sB", "v=1"
	// paramPnK, paramPnV           = "pn", "10"
	// paramRwK, paramRwV           = "rw", "new"
	// paramUsersetK, paramUsersetV = "userset", "0"
)

func safeSearchCookieString(safesearch bool) string {
	if safesearch {
		return fmt.Sprintf("%v&%v=%v", cookieSafeSearchPrefix, cookieSafeSearchK, "r")
	} else {
		return fmt.Sprintf("%v&%v=%v", cookieSafeSearchPrefix, cookieSafeSearchK, "p")
	}
}
