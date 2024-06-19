package googlescholar

import (
	"github.com/hearchco/agent/src/search/engines"
)

const (
	seName    = engines.GOOGLESCHOLAR
	searchURL = "https://scholar.google.com/scholar"
)

var origins = [...]engines.Name{seName}
