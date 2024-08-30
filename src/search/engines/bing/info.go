package bing

import (
	"github.com/hearchco/agent/src/search/engines"
)

const (
	seName         = engines.BING
	searchURL      = "https://www.bing.com/search"
	imageSearchURL = "https://www.bing.com/images/async"
)

var origins = [...]engines.Name{seName}
