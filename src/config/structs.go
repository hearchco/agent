package config

import (
	"time"

	"github.com/tminaorg/brzaguza/src/category"
	"github.com/tminaorg/brzaguza/src/engines"
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

// Delegates Timeout, PageTimeout to colly.Collector.SetRequestTimeout(); Note: See https://github.com/gocolly/colly/issues/644
// Delegates Delay, RandomDelay, Parallelism to colly.Collector.Limit()
type Timings struct {
	Timeout     time.Duration `koanf:"timeout"`
	PageTimeout time.Duration `koanf:"pagetimeout"`
	Delay       time.Duration `koanf:"delay"`
	RandomDelay time.Duration `koanf:"randomdelay"`
	Parallelism int           `koanf:"parallelism"`
}

type Settings struct {
	RequestedResultsPerPage int     `koanf:"requestedresults"`
	Shortcut                string  `koanf:"shortcut"`
	Timings                 Timings `koanf:"timings"`
}

type Redis struct {
	Host     string `koanf:"host"`
	Port     uint16 `koanf:"port"`
	Password string `koanf:"password"`
	Database uint8  `koanf:"database"`
}

type Cache struct {
	Type  string `koanf:"type"`
	Redis Redis  `koanf:"redis"`
}

type Server struct {
	Port        int    `koanf:"port"`
	FrontendUrl string `koanf:"frontendurl"`
	Cache       Cache  `koanf:"cache"`
}

type ReaderEngine struct {
	Enabled bool `koanf:"enabled"`
}

type ReaderCategory struct {
	REngines map[string]ReaderEngine `koanf:"engines"`
	Ranking  Ranking                 `koanf:"ranking"`
}

type ReaderConfig struct {
	Server      Server                           `koanf:"server"`
	RCategories map[category.Name]ReaderCategory `koanf:"categories"`
	Settings    map[string]Settings              `koanf:"settings"`
}

type Category struct {
	Engines []engines.Name
	Ranking Ranking
}

type Config struct {
	Server     Server
	Categories map[category.Name]Category
	Settings   map[engines.Name]Settings
}
