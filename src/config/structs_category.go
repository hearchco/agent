package config

import (
	"time"

	"github.com/hearchco/agent/src/search/engines"
)

// ReaderCategory is format in which the config is read from the config file and environment variables.
type ReaderCategory struct {
	REngines map[string]ReaderCategoryEngine `koanf:"engines"`
	Ranking  CategoryRanking                 `koanf:"ranking"`
	RTimings ReaderCategoryTimings           `koanf:"timings"`
}
type Category struct {
	Engines                  []engines.Name
	RequiredEngines          []engines.Name
	RequiredByOriginEngines  []engines.Name
	PreferredEngines         []engines.Name
	PreferredByOriginEngines []engines.Name
	Ranking                  CategoryRanking
	Timings                  CategoryTimings
}

// ReaderEngine is format in which the config is read from the config file and environment variables.
type ReaderCategoryEngine struct {
	// If false, the engine will not be used and other options will be ignored.
	Enabled bool `koanf:"enabled"`
	// If true, the engine will be awaited unless the hard timeout is reached.
	Required bool `koanf:"required"`
	// If true, the fastest engine that has this engine in "Origins" will be awaited unless the hard timeout is reached.
	// This means that we want to get results from this engine or any engine that has this engine in "Origins", whichever responds the fastest.
	RequiredByOrigin bool `koanf:"requiredbyorigin"`
	// If true, the engine will be awaited unless the preferred timeout is reached.
	Preferred bool `koanf:"preferred"`
	// If true, the fastest engine that has this engine in "Origins" will be awaited unless the preferred timeout is reached.
	// This means that we want to get results from this engine or any engine that has this engine in "Origins", whichever responds the fastest.
	PreferredByOrigin bool `koanf:"preferredbyorigin"`
}

type CategoryRanking struct {
	REXP    float64                          `koanf:"rexp"`
	A       float64                          `koanf:"a"`
	B       float64                          `koanf:"b"`
	C       float64                          `koanf:"c"`
	D       float64                          `koanf:"d"`
	TRA     float64                          `koanf:"tra"`
	TRB     float64                          `koanf:"trb"`
	TRC     float64                          `koanf:"trc"`
	TRD     float64                          `koanf:"trd"`
	Engines map[string]CategoryEngineRanking `koanf:"engines"`
}

type CategoryEngineRanking struct {
	Mul   float64 `koanf:"mul"`
	Const float64 `koanf:"const"`
}

// ReaderTimings is format in which the config is read from the config file and environment variables.
// In <number><unit> format.
// Example: 1s, 1m, 1h, 1d, 1w, 1M, 1y.
// If unit is not specified, it is assumed to be milliseconds.
// Delegates Delay, RandomDelay, Parallelism to colly.Collector.Limit().
type ReaderCategoryTimings struct {
	// Maximum amount of time to wait for the PreferredEngines (or ByOrigin) to respond.
	// If the search is still waiting for the RequiredEngines (or ByOrigin) after this time, the search will continue.
	PreferredTimeout string `koanf:"preferredtimeout"`
	// Hard timeout after which the search is forcefully stopped (even if the engines didn't respond).
	HardTimeout string `koanf:"hardtimeout"`
	// Colly delay.
	Delay string `koanf:"delay"`
	// Colly random delay.
	RandomDelay string `koanf:"randomdelay"`
	// Colly parallelism.
	Parallelism int `koanf:"parallelism"`
}

// Delegates Delay, RandomDelay, Parallelism to colly.Collector.Limit().
type CategoryTimings struct {
	// Maximum amount of time to wait for the PreferredEngines (or ByOrigin) to respond.
	// If the search is still waiting for the RequiredEngines (or ByOrigin) after this time, the search will continue.
	PreferredTimeout time.Duration
	// Hard timeout after which the search is forcefully stopped (even if the engines didn't respond).
	HardTimeout time.Duration
	// Colly delay.
	Delay time.Duration
	// Colly random delay.
	RandomDelay time.Duration
	// Colly parallelism.
	Parallelism int
}
