package config

import (
	"time"

	"github.com/hearchco/agent/src/search/engines"
)

var scienceEngines = []engines.Name{
	engines.GOOGLESCHOLAR,
}

var scienceRequiredEngines = []engines.Name{
	engines.GOOGLESCHOLAR,
}

var scienceRequiredByOriginEngines = []engines.Name{}

var sciencePreferredEngines = []engines.Name{}

var sciencePreferredByOriginEngines = []engines.Name{}

func scienceRanking() CategoryRanking {
	return ReqPrefOthRanking(scienceRequiredEngines, sciencePreferredEngines, scienceEngines)
}

var scienceTimings = CategoryTimings{
	PreferredTimeout: 700 * time.Millisecond,
	HardTimeout:      3 * time.Second,
}
