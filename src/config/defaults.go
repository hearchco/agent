package config

import (
	"time"

	"github.com/tminaorg/brzaguza/src/category"
	"github.com/tminaorg/brzaguza/src/engines"
)

const DefaultLocale string = "en-US"

func NewRanking() Ranking {
	return Ranking{
		REXP: 0.5,
		A:    1,
		B:    0,
		C:    1,
		D:    0,
		TRA:  1,
		TRB:  0,
		TRC:  1,
		TRD:  0,
		Engines: map[string]EngineRanking{
			engines.BING.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.BRAVE.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.DUCKDUCKGO.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.ETOOLS.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.GOOGLE.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.MOJEEK.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.PRESEARCH.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.QWANT.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.STARTPAGE.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.SWISSCOWS.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.YAHOO.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.YANDEX.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.YEP.ToLower(): {
				Mul:   1,
				Const: 0,
			},
		},
	}
}

func NewSettings() map[engines.Name]Settings {
	return map[engines.Name]Settings{
		engines.BING: {
			Shortcut: "bi",
		},
		engines.BRAVE: {
			Shortcut: "br",
		},
		engines.DUCKDUCKGO: {
			Shortcut: "ddg",
		},
		engines.ETOOLS: {
			Shortcut: "ets",
		},
		engines.GOOGLE: {
			Shortcut: "go",
		},
		engines.MOJEEK: {
			Shortcut: "mjk",
		},
		engines.PRESEARCH: {
			Shortcut: "ps",
		},
		engines.QWANT: {
			Shortcut: "qw",
		},
		engines.STARTPAGE: {
			Shortcut: "sp",
		},
		engines.SWISSCOWS: {
			Shortcut: "sc",
		},
		engines.YAHOO: {
			Shortcut: "yh",
		},
		engines.YEP: {
			Shortcut: "yep",
		},
	}
}

func NewAllEnabled() []engines.Name {
	return []engines.Name{
		engines.BING,
		engines.BRAVE,
		engines.DUCKDUCKGO,
		engines.ETOOLS,
		engines.GOOGLE,
		engines.MOJEEK,
		engines.PRESEARCH,
		engines.QWANT,
		engines.STARTPAGE,
		engines.SWISSCOWS,
		engines.YAHOO,
		engines.YEP,
	}
}

func NewInfo() []engines.Name {
	return []engines.Name{
		engines.BING,
		engines.GOOGLE,
		engines.MOJEEK,
	}
}

func New() *Config {
	return &Config{
		Server: Server{
			Port:         3030,
			FrontendUrls: []string{"http://localhost:8000"},
			Cache: Cache{
				Type: "pebble",
				Redis: Redis{
					Host: "localhost",
					Port: 6379,
				},
			},
		},
		Settings: NewSettings(),
		Categories: map[category.Name]Category{
			category.GENERAL: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     1000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.INFO: {
				Engines: NewInfo(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     1000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.SCIENCE: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     3000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.NEWS: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     1000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.BLOG: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     2500 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.SURF: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     2000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.NEWNEWS: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     1000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
		},
	}
}
