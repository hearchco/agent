package config

import (
	"time"

	"github.com/hearchco/agent/src/search/engines"
)

var imagesEngines = []engines.Name{
	engines.BING,
	engines.GOOGLE,
}

var imagesRequiredEngines = []engines.Name{}

var imagesRequiredByOriginEngines = []engines.Name{
	engines.BING,
	engines.GOOGLE,
}

var imagesPreferredEngines = []engines.Name{}

var imagesPreferredByOriginEngines = []engines.Name{}

func imagesRanking() CategoryRanking {
	return EmptyRanking(imagesEngines)
}

var imagesTimings = CategoryTimings{
	PreferredTimeout: 500 * time.Millisecond,
	HardTimeout:      1500 * time.Millisecond,
}
