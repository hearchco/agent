package swisscows

import (
	"github.com/hearchco/agent/src/search/engines"
)

const (
	seName    = engines.SWISSCOWS
	searchURL = "https://api.swisscows.com/web/search"
)

var origins = [...]engines.Name{engines.SWISSCOWS, engines.BING}
