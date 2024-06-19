package qwant

import (
	"github.com/hearchco/agent/src/search/engines"
)

const (
	seName    = engines.QWANT
	searchURL = "https://api.qwant.com/v3/search/web"
)

var origins = [...]engines.Name{engines.QWANT, engines.BING}
