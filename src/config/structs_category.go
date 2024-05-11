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
// Delegates Timeout, PageTimeout to colly.Collector.SetRequestTimeout(); Note: See https://github.com/gocolly/colly/issues/644
// Delegates Delay, RandomDelay, Parallelism to colly.Collector.Limit()
type ReaderCategoryTimings struct {
	// Preferred timeout if enough results are found
	PreferredTimeout string `koanf:"preferredtimeout"`
	// Number of results which if not met will trigger the additional timeout
	PreferredTimeoutResults int `koanf:"preferredtimeoutresults"`
	// Additional timeout if not enough results are found (delay after which the number of results is checked)
	AdditionalTimeout string `koanf:"additionaltimeout"`
	// Hard timeout after which the search is forcefully stopped
	HardTimeout string `koanf:"hardtimeout"`
	// Colly collector timeout (should be less than or equal to HardTimeout)
	Timeout string `koanf:"timeout"`
	// Colly collector page timeout (should be less than or equal to HardTimeout)
	PageTimeout string `koanf:"pagetimeout"`
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
	// Preferred timeout if enough results are found
	PreferredTimeout time.Duration
	// Number of results which if not met will trigger the additional timeout
	PreferredTimeoutResults int
	// Additional timeout if not enough results are found (delay after which the number of results is checked)
	AdditionalTimeout time.Duration
	// Hard timeout after which the search is forcefully stopped
	HardTimeout time.Duration
	// Colly collector timeout (should be less than or equal to HardTimeout)
	Timeout time.Duration
	// Colly collector page timeout (should be less than or equal to HardTimeout)
	PageTimeout time.Duration
	// Colly delay
	Delay time.Duration
	// Colly random delay
	RandomDelay time.Duration
	// Colly parallelism
	Parallelism int
}
