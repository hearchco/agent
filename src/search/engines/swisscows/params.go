package swisscows

import (
	"fmt"
	"strings"

	"github.com/hearchco/agent/src/search/engines/options"
)

const (
	paramKeyPage   = "offset"
	paramKeyLocale = "region" // Should be the same as Locale, only with "_" replaced by "-".

	paramFreshness = "freshness=All"
	paramItems     = "itemsCount=10"
)

func localeParamString(locale options.Locale) string {
	region := strings.Replace(locale.String(), "_", "-", 1)
	return fmt.Sprintf("%v=%v", paramKeyLocale, region)
}
