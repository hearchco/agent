package config

import (
	"time"

	"github.com/hearchco/agent/src/search/engines"
)

var suggestionsEngines = []engines.Name{
	engines.DUCKDUCKGO,
	engines.GOOGLE,
}

var suggestionsRequiredEngines = []engines.Name{
	engines.DUCKDUCKGO,
	engines.GOOGLE,
}

var suggestionsRequiredByOriginEngines = []engines.Name{}

var suggestionsPreferredEngines = []engines.Name{}

var suggestionsPreferredByOriginEngines = []engines.Name{}

func suggestionsRanking() CategoryRanking {
	return EmptyRanking(suggestionsEngines)
}

var suggestionsTimings = CategoryTimings{
	PreferredTimeout: 300 * time.Millisecond,
	HardTimeout:      500 * time.Millisecond,
}
