package presearch

import (
	"github.com/hearchco/agent/src/search/engines"
)

const (
	seName    = engines.PRESEARCH
	searchURL = "https://presearch.com/search"
)

var origins = [...]engines.Name{seName, engines.GOOGLE}
