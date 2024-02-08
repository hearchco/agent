package config

import (
	"time"

	"github.com/hearchco/hearchco/src/category"
	"github.com/hearchco/hearchco/src/engines"
)

type EngineRanking struct {
	Mul   float64 `koanf:"mul"`
	Const float64 `koanf:"const"`
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

type Settings struct {
	RequestedResultsPerPage int    `koanf:"requestedresults"`
	Shortcut                string `koanf:"shortcut"`
}

// ReaderTTL is format in which the config is read from the config file
// in <number><unit> format
// example: 1s, 1m, 1h, 1d, 1w, 1M, 1y
// if unit is not specified, it is assumed to be milliseconds
type ReaderTTL struct {
	// how long to store the results in cache
	// setting this to 0 caches the results forever
	// to disable caching set conf.Cache.Type to "none"
	Time string `koanf:"time"`
	// if the remaining TTL when retrieving from cache is less than this, update the cache entry and reset the TTL
	// setting this to 0 disables this feature
	// setting this to the same value (or higher) as Results will update the cache entry every time
	RefreshTime string `koanf:"refreshtime"`
}
type TTL struct {
	Time        time.Duration
	RefreshTime time.Duration
}

type Badger struct {
	// setting this to false will result in badger not persisting the cache to disk
	// that means that badger will run in memory only
	Persist bool `koanf:"persist"`
}

type Redis struct {
	Host     string `koanf:"host"`
	Port     uint16 `koanf:"port"`
	Password string `koanf:"password"`
	Database uint8  `koanf:"database"`
}

// ReaderCache is format in which the config is read from the config file
type ReaderCache struct {
	// can be "none", "badger" or "redis"
	Type string `koanf:"type"`
	// has no effect if Type is "none"
	TTL ReaderTTL `koanf:"ttl"`
	// badger specific settings
	Badger Badger `koanf:"badger"`
	// redis specific settings
	Redis Redis `koanf:"redis"`
}
type Cache struct {
	Type   string
	TTL    TTL
	Badger Badger
	Redis  Redis
}

// ReaderServer is format in which the config is read from the config file
type ReaderServer struct {
	// port on which the API server listens
	Port int `koanf:"port"`
	// frontend url needed for CORS
	FrontendUrl string `koanf:"frontendurl"`
	// cache settings
	Cache ReaderCache `koanf:"cache"`
}
type Server struct {
	Port        int `koanf:"port"`
	FrontendUrl string
	Cache       Cache
}

// ReaderEngine is format in which the config is read from the config file
type ReaderEngine struct {
	Enabled bool `koanf:"enabled"`
}

// ReaderTimings is format in which the config is read from the config file
// in miliseconds
type ReaderTimings struct {
	// HardTimeout uint `koanf:"hardTimeout"`
	Timeout     uint `koanf:"timeout"`
	PageTimeout uint `koanf:"pagetimeout"`
	Delay       uint `koanf:"delay"`
	RandomDelay uint `koanf:"randomdelay"`
	Parallelism int  `koanf:"parallelism"`
}

// Delegates Timeout, PageTimeout to colly.Collector.SetRequestTimeout(); Note: See https://github.com/gocolly/colly/issues/644
// Delegates Delay, RandomDelay, Parallelism to colly.Collector.Limit()
type Timings struct {
	// HardTimeout time.Duration
	Timeout     time.Duration
	PageTimeout time.Duration
	Delay       time.Duration
	RandomDelay time.Duration
	Parallelism int
}

// ReaderCategory is format in which the config is read from the config file
type ReaderCategory struct {
	REngines map[string]ReaderEngine `koanf:"engines"`
	Ranking  Ranking                 `koanf:"ranking"`
	RTimings ReaderTimings           `koanf:"timings"`
}
type Category struct {
	Engines []engines.Name
	Ranking Ranking
	Timings Timings
}

// ReaderConfig is format in which the config is read from the config file
type ReaderConfig struct {
	Server      ReaderServer                     `koanf:"server"`
	RCategories map[category.Name]ReaderCategory `koanf:"categories"`
	Settings    map[string]Settings              `koanf:"settings"`
}
type Config struct {
	Server     Server
	Categories map[category.Name]Category
	Settings   map[engines.Name]Settings
}
