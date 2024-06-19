package bingimages

import (
	"github.com/hearchco/agent/src/search/engines"
)

const (
	seName    = engines.BINGIMAGES
	searchURL = "https://www.bing.com/images/async"
)

var origins = [...]engines.Name{seName}
