package etools

import (
	"fmt"
)

func safeSearchParamString(safesearch bool) string {
	if safesearch {
		return fmt.Sprintf("%v=%v", params.SafeSearch, "true")
	} else {
		return fmt.Sprintf("%v=%v", params.SafeSearch, "false")
	}
}
