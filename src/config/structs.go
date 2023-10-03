package config

import "time"

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

// Server config
type Server struct {
	Port         int      `koanf:"port"`
	FrontendUrls []string `koanf:"frontendUrls"`
	RedisUrl     string   `koanf:"redisUrl"`
}

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

// Config struct for Koanf
type Config struct {
	Server  Server            `koanf:"server"`
	Engines map[string]Engine `koanf:"engines"`
	Ranking Ranking           `koanf:"ranking"`
}
