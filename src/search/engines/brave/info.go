package brave

import (
	"github.com/hearchco/agent/src/search/engines"
)

const (
	seName    = engines.BRAVE
	searchURL = "https://search.brave.com/search"
)

var origins = [...]engines.Name{seName, engines.GOOGLE}
