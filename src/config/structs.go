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

// Config struct for Koanf
type Config struct {
	Server  Server            `koanf:"server"`
	Engines map[string]Engine `koanf:"engines"`
}
