package googleimages

import (
	"github.com/hearchco/agent/src/search/engines"
)

const (
	seName    = engines.GOOGLEIMAGES
	searchURL = "https://www.google.com/search"
)

var origins = [...]engines.Name{seName}
