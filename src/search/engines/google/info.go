package google

import (
	"github.com/hearchco/agent/src/search/engines"
)

const (
	seName         = engines.GOOGLE
	searchURL      = "https://www.google.com/search"
	imageSearchURL = "https://www.google.com/search"
	suggestURL     = "https://suggestqueries.google.com/complete/search"
)

var origins = [...]engines.Name{seName}
