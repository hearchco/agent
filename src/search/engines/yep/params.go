package yep

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/hearchco/agent/src/search/engines/options"
// )

// const (
// 	paramKeyPage       = "limit"
// 	paramKeyLocale     = "gl"         // Should be last 2 characters of Locale.
// 	paramKeySafeSearch = "safeSearch" // Can be "off" or "strict".

// 	paramClient     = "client=web"
// 	paramNo_correct = "no_correct=false"
// 	paramType       = "type=web"
// )

// func localeParamString(locale options.Locale) string {
// 	country := strings.Split(locale.String(), "_")[1]
// 	return fmt.Sprintf("%v=%v", paramKeyLocale, country)
// }

// func safeSearchParamString(safesearch bool) string {
// 	if safesearch {
// 		return fmt.Sprintf("%v=%v", paramKeySafeSearch, "strict")
// 	} else {
// 		return fmt.Sprintf("%v=%v", paramKeySafeSearch, "off")
// 	}
// }
