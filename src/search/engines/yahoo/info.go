package yahoo

import (
	"github.com/hearchco/agent/src/search/engines"
)

const (
	seName    = engines.YAHOO
	searchURL = "https://search.yahoo.com/search"
)

var origins = [...]engines.Name{engines.YAHOO, engines.BING}
