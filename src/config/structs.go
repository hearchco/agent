package config

import "time"

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
	PageTimeout time.Duration `koanf:"pageTimeout"`
	Delay       time.Duration `koanf:"delay"`
	RandomDelay time.Duration `koanf:"randomDelay"`
	Parallelism int           `koanf:"parallelism"`
}

type Settings struct {
	RequestedResultsPerPage int     `koanf:"requestedResults"`
	Shortcut                string  `koanf:"shortcut"`
	Timings                 Timings `koanf:"timings"`
}

type Engine struct {
	Enabled  bool     `koanf:"enabled"`
	Settings Settings `koanf:"settings"`
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
	Port         int      `koanf:"port"`
	FrontendUrls []string `koanf:"frontendUrls"`
	Cache        Cache    `koanf:"cache"`
}

type Category struct {
	Engines map[string]Engine `koanf:"engines"`
}

type Config struct {
	Server     Server              `koanf:"server"`
	Categories map[string]Category `koanf:"categories"`
	Ranking    Ranking             `koanf:"ranking"`
}
