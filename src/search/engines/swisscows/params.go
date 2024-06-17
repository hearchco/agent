package swisscows

import (
	"fmt"
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

func localeParamString(locale options.Locale) string {
	region := strings.Replace(locale.String(), "_", "-", 1)
	return fmt.Sprintf("%v=%v", params.Locale, region)
}
