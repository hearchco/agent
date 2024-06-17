package config

import (
	"time"

	"github.com/hearchco/agent/src/search/engines"
)

var imagesEngines = []engines.Name{
	engines.BINGIMAGES,
	engines.GOOGLEIMAGES,
}

var imagesRequiredEngines = []engines.Name{}

var imagesRequiredByOriginEngines = []engines.Name{
	engines.BINGIMAGES,
	engines.GOOGLEIMAGES,
}

var imagesPreferredEngines = []engines.Name{}

var imagesPreferredByOriginEngines = []engines.Name{}

func imagesRanking() CategoryRanking {
	return EmptyRanking(imagesEngines)
}

var imagesTimings = CategoryTimings{
	PreferredTimeout: 700 * time.Millisecond,
	HardTimeout:      3 * time.Second,
}
