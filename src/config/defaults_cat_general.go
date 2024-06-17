package config

import (
	"time"

	"github.com/hearchco/agent/src/search/engines"
)

var generalEngines = []engines.Name{
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

var generalRequiredEngines = []engines.Name{}

var generalRequiredByOriginEngines = []engines.Name{
	engines.BING,
	engines.GOOGLE,
}

var generalPreferredEngines = []engines.Name{
	engines.ETOOLS,
}

var generalPreferredByOriginEngines = []engines.Name{
	engines.BRAVE,
	engines.MOJEEK,
	// engines.YEP,
}

func generalRanking() CategoryRanking {
	return EmptyRanking(generalEngines)
}

var generalTimings = CategoryTimings{
	PreferredTimeout: 500 * time.Millisecond,
	HardTimeout:      1500 * time.Millisecond,
}
