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
	REXP    float64                  `koanf:"rexp"`
	A       float64                  `koanf:"a"`
	B       float64                  `koanf:"b"`
	C       float64                  `koanf:"c"`
	D       float64                  `koanf:"d"`
	TRA     float64                  `koanf:"tra"`
	TRB     float64                  `koanf:"trb"`
	TRC     float64                  `koanf:"trc"`
	TRD     float64                  `koanf:"trd"`
	Engines map[string]EngineRanking `koanf:"engines"`
}

type EngineRanking struct {
	Mul   float64 `koanf:"mul"`
	Const float64 `koanf:"const"`
}

type Timings struct {
	// Maximum amount of time to wait for the PreferredEngines (or ByOrigin) to respond.
	// If the search is still waiting for the RequiredEngines (or ByOrigin) after this time, the search will continue.
	PreferredTimeout time.Duration
	// Hard timeout after which the search is forcefully stopped (even if the engines didn't respond).
	HardTimeout time.Duration
}
