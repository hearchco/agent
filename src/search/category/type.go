package category

import (
	"time"

	"github.com/hearchco/agent/src/search/engines"
)

type Category struct {
	Engines                  []engines.Name
	RequiredEngines          []engines.Name
	RequiredByOriginEngines  []engines.Name
	PreferredEngines         []engines.Name
	PreferredByOriginEngines []engines.Name
	Ranking                  Ranking
	Timings                  Timings
}

type Ranking struct {
	// The exponent, multiplier and addition used on the rank itself.
	RankExp float64
	RankMul float64
	RankAdd float64
	// The multiplier and addition used on the rank score (number calculated from dividing 100 with the rank + above variables applied).
	RankScoreMul float64
	RankScoreAdd float64
	// The multiplier and addition used on the number of times the result was returned.
	TimesReturnedMul float64
	TimesReturnedAdd float64
	// The multiplier and addition used on the times returned score (number calculated from doing log(timesReturnedNum + above variables applied)).
	TimesReturnedScoreMul float64
	TimesReturnedScoreAdd float64
	// Multipliers and additions for each engine, applied to the rank score.
	Engines map[engines.Name]EngineRanking
}

type EngineRanking struct {
	Mul float64
	Add float64
}

type Timings struct {
	// Maximum amount of time to wait for the PreferredEngines (or ByOrigin) to respond.
	// If the search is still waiting for the RequiredEngines (or ByOrigin) after this time, the search will continue.
	PreferredTimeout time.Duration
	// Hard timeout after which the search is forcefully stopped (even if the engines didn't respond).
	HardTimeout time.Duration
}
