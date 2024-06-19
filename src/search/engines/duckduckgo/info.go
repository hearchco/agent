package duckduckgo

import (
	"github.com/hearchco/agent/src/search/engines"
)

const (
	seName    = engines.DUCKDUCKGO
	searchURL = "https://lite.duckduckgo.com/lite/"
)

var origins = [...]engines.Name{seName, engines.BING}
