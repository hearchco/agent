package config

import "github.com/tminaorg/brzaguza/src/engines"

const DefaultLocale string = "en-US"

var DefaultConfig Config = Config{
	Server: Server{
		Port:        3030,
		FrontendUrl: "http://localhost:8000",
		RedisUrl:    "http://localhost:6379",
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
}
