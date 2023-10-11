package config

import (
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
			engines.Bing.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.Brave.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.DuckDuckGo.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.Etools.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.Google.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.Mojeek.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.Presearch.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.Qwant.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.Startpage.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.Swisscows.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.Yahoo.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.Yandex.ToLower(): {
				Mul:   1,
				Const: 0,
			},
			engines.Yep.ToLower(): {
				Mul:   1,
				Const: 0,
			},
		},
	}
}

func NewSettings() map[string]Settings {
	return map[string]Settings{
		engines.Bing.ToLower(): {
			Shortcut: "bi",
		},
		engines.Brave.ToLower(): {
			Shortcut: "br",
		},
		engines.DuckDuckGo.ToLower(): {
			Shortcut: "ddg",
		},
		engines.Etools.ToLower(): {
			Shortcut: "ets",
		},
		engines.Google.ToLower(): {
			Shortcut: "go",
		},
		engines.Mojeek.ToLower(): {
			Shortcut: "mjk",
		},
		engines.Presearch.ToLower(): {
			Shortcut: "ps",
		},
		engines.Qwant.ToLower(): {
			Shortcut: "qw",
		},
		engines.Startpage.ToLower(): {
			Shortcut: "sp",
		},
		engines.Swisscows.ToLower(): {
			Shortcut: "sc",
		},
		engines.Yahoo.ToLower(): {
			Shortcut: "yh",
		},
		engines.Yep.ToLower(): {
			Shortcut: "yep",
		},
	}
}

func NewAllEnabled() map[string]Engine {
	return map[string]Engine{
		engines.Bing.ToLower(): {
			Enabled: true,
		},
		engines.Brave.ToLower(): {
			Enabled: true,
		},
		engines.DuckDuckGo.ToLower(): {
			Enabled: true,
		},
		engines.Etools.ToLower(): {
			Enabled: true,
		},
		engines.Google.ToLower(): {
			Enabled: true,
		},
		engines.Mojeek.ToLower(): {
			Enabled: true,
		},
		engines.Presearch.ToLower(): {
			Enabled: true,
		},
		engines.Qwant.ToLower(): {
			Enabled: true,
		},
		engines.Startpage.ToLower(): {
			Enabled: true,
		},
		engines.Swisscows.ToLower(): {
			Enabled: true,
		},
		engines.Yahoo.ToLower(): {
			Enabled: true,
		},
		engines.Yep.ToLower(): {
			Enabled: true,
		},
	}
}

func NewInfo() map[string]Engine {
	return map[string]Engine{
		engines.Bing.ToLower(): {
			Enabled: true,
		},
		engines.Brave.ToLower(): {
			Enabled: false,
		},
		engines.DuckDuckGo.ToLower(): {
			Enabled: false,
		},
		engines.Etools.ToLower(): {
			Enabled: false,
		},
		engines.Google.ToLower(): {
			Enabled: true,
		},
		engines.Mojeek.ToLower(): {
			Enabled: true,
		},
		engines.Presearch.ToLower(): {
			Enabled: false,
		},
		engines.Qwant.ToLower(): {
			Enabled: false,
		},
		engines.Startpage.ToLower(): {
			Enabled: false,
		},
		engines.Swisscows.ToLower(): {
			Enabled: false,
		},
		engines.Yahoo.ToLower(): {
			Enabled: false,
		},
		engines.Yep.ToLower(): {
			Enabled: false,
		},
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
			},
			category.INFO: {
				Engines: NewInfo(),
				Ranking: NewRanking(),
			},
			category.SCIENCE: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
			},
			category.NEWS: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
			},
			category.BLOG: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
			},
			category.SURF: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
			},
			category.NEWNEWS: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
			},
		},
	}
}
