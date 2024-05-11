package config

import (
	"time"

	"github.com/hearchco/hearchco/src/search/engines"
)

// ReaderCategory is format in which the config is read from the config file
type ReaderCategory struct {
	REngines map[string]ReaderCategoryEngine `koanf:"engines"`
	Ranking  CategoryRanking                 `koanf:"ranking"`
	RTimings ReaderCategoryTimings           `koanf:"timings"`
}
type Category struct {
	Engines []engines.Name
	Ranking CategoryRanking
	Timings CategoryTimings
}

// ReaderEngine is format in which the config is read from the config file
type ReaderCategoryEngine struct {
	Enabled bool `koanf:"enabled"`
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

// ReaderTimings is format in which the config is read from the config file
// In <number><unit> format
// Example: 1s, 1m, 1h, 1d, 1w, 1M, 1y
// If unit is not specified, it is assumed to be milliseconds
// Delegates Timeout, PageTimeout to colly.Collector.SetRequestTimeout(); Note: See https://github.com/gocolly/colly/issues/644
// Delegates Delay, RandomDelay, Parallelism to colly.Collector.Limit()
type ReaderCategoryTimings struct {
	// Minimum amount of time to wait before starting to check the number of results
	// Search will wait for at least this amount of time (unless all engines respond)
	PreferredTimeoutMin string `koanf:"preferredtimeoutmin"`
	// Maximum amount of time to wait until the number of results is satisfactory
	// Search will wait for at most this amount of time (unless all engines respond or the preferred number of results is found)
	PreferredTimeoutMax string `koanf:"preferredtimeoutmax"`
	// Preferred number of results to find
	PreferredResultsNumber int `koanf:"preferredresultsnumber"`
	// Time of the steps for checking if the number of results is satisfactory
	StepTime string `koanf:"steptime"`
	// Minimum number of results required after the maximum preferred time
	// If this number isn't met, the search will continue after the maximum preferred time
	MinimumResultsNumber int `koanf:"minimumresultsnumber"`
	// Hard timeout after which the search is forcefully stopped (even if the engines didn't respond)
	HardTimeout string `koanf:"hardtimeout"`
	// Colly delay
	Delay string `koanf:"delay"`
	// Colly random delay
	RandomDelay string `koanf:"randomdelay"`
	// Colly parallelism
	Parallelism int `koanf:"parallelism"`
}

// Delegates Timeout, PageTimeout to colly.Collector.SetRequestTimeout(); Note: See https://github.com/gocolly/colly/issues/644
// Delegates Delay, RandomDelay, Parallelism to colly.Collector.Limit()
type CategoryTimings struct {
	// Minimum amount of time to wait before starting to check the number of results
	// Search will wait for at least this amount of time (unless all engines respond)
	PreferredTimeoutMin time.Duration
	// Maximum amount of time to wait until the number of results is satisfactory
	// Search will wait for at most this amount of time (unless all engines respond or the preferred number of results is found)
	PreferredTimeoutMax time.Duration
	// Preferred number of results to find
	PreferredResultsNumber int
	// Time of the steps for checking if the number of results is satisfactory
	StepTime time.Duration
	// Minimum number of results required after the maximum preferred time
	// If this number isn't met, the search will continue after the maximum preferred time
	MinimumResultsNumber int
	// Hard timeout after which the search is forcefully stopped (even if the engines didn't respond)
	HardTimeout time.Duration
	// Colly delay
	Delay time.Duration
	// Colly random delay
	RandomDelay time.Duration
	// Colly parallelism
	Parallelism int
}
