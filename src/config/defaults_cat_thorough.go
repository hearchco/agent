package config

import (
	"time"

	"github.com/hearchco/agent/src/search/engines"
)

var thoroughEngines = []engines.Name{
	engines.BING,
	engines.BRAVE,
	engines.DUCKDUCKGO,
	engines.ETOOLS,
	engines.GOOGLE,
	engines.MOJEEK,
	engines.PRESEARCH,
	engines.QWANT,
	engines.STARTPAGE,
	engines.SWISSCOWS,
	engines.YAHOO,
	// engines.YEP,
}

var thoroughRequiredEngines = []engines.Name{
	engines.BING,
	engines.BRAVE,
	engines.DUCKDUCKGO,
	engines.ETOOLS,
	engines.GOOGLE,
	engines.MOJEEK,
	engines.PRESEARCH,
	engines.QWANT,
	engines.STARTPAGE,
	engines.SWISSCOWS,
	engines.YAHOO,
	// engines.YEP,
}

var thoroughRequiredByOriginEngines = []engines.Name{}

var thoroughPreferredEngines = []engines.Name{}

var thoroughPreferredByOriginEngines = []engines.Name{}

func thoroughRanking() CategoryRanking {
	return EmptyRanking(thoroughEngines)
}

var thoroughTimings = CategoryTimings{
	PreferredTimeout: 3 * time.Second,
	HardTimeout:      5 * time.Second,
}
