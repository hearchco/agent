package bingimages

import (
	"fmt"
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

func localeParamString(locale options.Locale) string {
	spl := strings.SplitN(strings.ToLower(locale.String()), "_", 2)
	return fmt.Sprintf("%v=%v&%v=%v", params.Locale, spl[0], params.LocaleSec, spl[1])
}
