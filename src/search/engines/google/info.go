package google

import (
	"github.com/hearchco/agent/src/search/engines"
)

const (
	seName    = engines.GOOGLE
	searchURL = "https://www.google.com/search"
)

var origins = [...]engines.Name{seName}
