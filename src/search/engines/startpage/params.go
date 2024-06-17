package startpage

import (
	"fmt"
)

func safeSearchParamString(safesearch bool) string {
	if safesearch {
		return ""
	} else {
		return fmt.Sprintf("%v=%v", params.SafeSearch, "none")
	}
}
