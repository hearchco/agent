package config

import "github.com/tminaorg/brzaguza/src/engines"

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
		Engines: map[string]Engine{
			engines.Bing.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "bi",
				},
			},
			engines.Brave.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "br",
				},
			},
			engines.DuckDuckGo.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "ddg",
				},
			},
			engines.Etools.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "ets",
				},
			},
			engines.Google.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "go",
				},
			},
			engines.Mojeek.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "mjk",
				},
			},
			engines.Qwant.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "qw",
				},
			},
			engines.Startpage.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "sp",
				},
			},
			engines.Swisscows.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "sc",
				},
			},
			engines.Yep.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "yep",
				},
			},
		},
		Ranking: NewRanking(),
	}
}
